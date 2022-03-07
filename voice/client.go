package voice

import (
	"bothoi/config"
	"bothoi/models"
	"bothoi/references/voice_opcode"
	"bothoi/states"
	"bothoi/util"
	"bothoi/util/ws_util"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jonas747/dca/v2"
	"github.com/kkdai/youtube/v2"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/nacl/secretbox"
)

type VoiceClient struct {
	SongQueue          *models.SongQueue
	c                  *websocket.Conn
	SessionDescription *models.SessionDescription
	UDPInfo            *models.UDPInfo
	waitDesc           sync.WaitGroup
}

func (client *VoiceClient) Connect() {
	defer states.RemoveGuildFromSongQueue(client.SongQueue.GuildID)

	c, _, err := websocket.DefaultDialer.Dial("wss://"+client.SongQueue.VoiceServer.Endpoint+"?v="+config.VOICE_GATEWAY_VERSION, nil)
	if err != nil {
		log.Fatalln(err)
	}
	client.c = c
	defer c.Close()

	ws_util.WriteJSONLog(client.c, models.NewVoiceIdentify(client.SongQueue.GuildID, config.BOT_ID, *client.SongQueue.SessionID, client.SongQueue.VoiceServer.Token))

	heatbeatInterval := make(chan int)
	heatbeatAcked := make(chan int64)

	// receive the gateway response
	go func() {
		for {
			var payload models.GatewayPayload
			err := c.ReadJSON(&payload)
			if err != nil {
				log.Println(err)
				continue
			}
			if config.DEVELOPMENT {
				jsonDat, _ := json.Marshal(payload)
				log.Println("incoming: ", payload, string(jsonDat))
			} else {
				log.Println("incoming: ", payload)
			}
			switch payload.Op {
			case voice_opcode.Hello:
				heatbeatInterval <- int(payload.D.(map[string]interface{})["heartbeat_interval"].(float64))
			case voice_opcode.HeartbeatAck:
				n, _ := strconv.ParseInt(payload.D.(string), 10, 64)
				heatbeatAcked <- n
			case voice_opcode.Ready:
				var data models.UDPInfo
				mapstructure.Decode(payload.D, &data)
				if !util.ContainsStr(data.Modes, config.PREFERRED_MODE) {
					log.Panicln("Preferred mode not available")
				}
				client.UDPInfo = &data
				// start the UDP connection
				go client.StartVoice()
			case voice_opcode.SessionDescription:
				var data models.SessionDescription
				mapstructure.Decode(payload.D, &data)
				client.SessionDescription = &data
				client.waitDesc.Done()
			}
		}
	}()

	// keeping heartbeats and prevent application from closing.
	interval := <-heatbeatInterval
	prevNonce := int64(0)
	// here to skip first ack in loop
	heatbeatAcked <- 0
	for {
		// wait for heartbeat ack
		select {
		case hbNonce := <-heatbeatAcked:
			if prevNonce != hbNonce {
				// handle nonce error
				c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(4006, ""))
				return
			}
		case <-time.After(time.Duration(interval) * time.Millisecond):
			// uh oh timeout, reconnect
			log.Println("timeout, attempting to reconnect")
			c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(4009, ""))
			return
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)

		randNum, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			log.Println(err)
		}
		prevNonce = randNum.Int64()

		ws_util.WriteJSONLog(c, models.NewVoiceHeartbeat(prevNonce))
	}
}

func (client *VoiceClient) StartVoice() {
	client.waitDesc.Add(1)

	raddr, err := net.ResolveUDPAddr("udp", client.UDPInfo.IP+":"+strconv.Itoa(int(client.UDPInfo.Port)))
	conn, err := net.DialUDP("udp", nil, raddr)
	ip, port := client.performIpDiscovery(conn)
	ws_util.WriteJSONLog(client.c, models.NewVoiceSelectProtocol(ip, port, config.PREFERRED_MODE))

	client.waitDesc.Wait()
	go keepAlive(conn, time.Second*5)

	videoID := "nppkHBfM34g"
	ytClient := youtube.Client{}

	video, err := ytClient.GetVideo(videoID)
	if err != nil {
		panic(err)
	}

	formats := video.Formats.WithAudioChannels() // only get videos with audio
	for _, format := range formats {
		fmt.Println(format.AudioSampleRate, format.AudioQuality, format.Bitrate, format.QualityLabel, format.AudioChannels, format.Width, format.Height)
	}
	sUrl, err := ytClient.GetStreamURL(video, &formats[1])

	options := dca.StdEncodeOptions
	options.FrameRate = 48000
	options.FrameDuration = 20
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"

	encodeSession, err := dca.EncodeFile(sUrl, options)
	if err != nil {
		log.Panicln(err)
	}
	// Make sure everything is cleaned up, that for example the encoding process if any issues happened isnt lingering around
	defer encodeSession.Cleanup()
	log.Println("encoded")
	fr := encodeSession.Options().FrameRate
	fd := encodeSession.Options().FrameDuration
	timeStampWidth := uint32(fr * fd / 1000)
	ticker := time.NewTicker(time.Millisecond * time.Duration(fd))
	ws_util.WriteJSONLog(client.c, models.NewVoiceSpeaking(client.UDPInfo.Ssrc))

	randNum, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		log.Println(err)
	}
	var sequenceNumber = uint16(randNum.Uint64())
	randNum, err = rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		log.Println(err)
	}
	var timeStamp = uint32(randNum.Uint64())

	header := make([]byte, 12)
	header[0] = 0x80
	header[1] = 0x78
	binary.BigEndian.PutUint32(header[8:], client.UDPInfo.Ssrc)

	for {
		binary.BigEndian.PutUint16(header[2:], uint16(sequenceNumber))
		binary.BigEndian.PutUint16(header[4:], uint16(timeStamp))

		var nonce [24]byte
		copy(nonce[:], header[0:12])

		frame, err := encodeSession.OpusFrame()
		if err != nil {
			if err != io.EOF {
				log.Panicln(err)
			}
		}
		packet := secretbox.Seal(header, frame, &nonce, &client.SessionDescription.SecretKey)

		<-ticker.C

		_, err = conn.Write(packet)
		if err != nil {
			log.Panicln(err)
		}

		if (sequenceNumber) == 0xFFFF {
			sequenceNumber = 0
		} else {
			sequenceNumber++
		}

		if (timeStamp + timeStampWidth) == 0xFFFFFFFF {
			timeStamp = 0
		} else {
			timeStamp += timeStampWidth
		}
	}
}

func (client *VoiceClient) performIpDiscovery(conn *net.UDPConn) (string, uint16) {
	send := make([]byte, 70)

	binary.BigEndian.PutUint32(send, client.UDPInfo.Ssrc)
	_, err := conn.Write(send)
	if err != nil {
		log.Panicln(err)
	}

	receive := make([]byte, 70)
	rlen, _, err := conn.ReadFromUDP(receive)
	if err != nil {
		log.Panicln(err)
	}

	if rlen < 70 {
		log.Panicln("UDP packet too short")
	}

	var ip string
	for i := 4; i < 20; i++ {
		if receive[i] == 0 {
			break
		}
		ip += string(receive[i])
	}

	// Grab port from position 68 and 69
	port := binary.BigEndian.Uint16(receive[68:70])
	return ip, port

}

func keepAlive(conn *net.UDPConn, i time.Duration) {

	var sequence uint64

	packet := make([]byte, 8)

	ticker := time.NewTicker(i)
	defer ticker.Stop()
	for {

		binary.LittleEndian.PutUint64(packet, sequence)
		sequence++

		_, err := conn.Write(packet)
		if err != nil {
			log.Println(err)
			return
		}

		<-ticker.C
	}
}
