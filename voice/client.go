package voice

import (
	"bothoi/config"
	"bothoi/gateway"
	"bothoi/models"
	"bothoi/references/voice_opcode"
	"bothoi/util"
	"bothoi/util/ws_util"
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"log"
	"math"
	"math/big"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/nacl/secretbox"
)

type VoiceClient struct {
	sync.RWMutex
	ctx                context.Context
	ctxCancel          context.CancelFunc
	guildID            string
	voiceChannelID     string
	sessionID          *string
	songQueue          []*models.SongItemWData
	voiceServer        *models.VoiceServer
	c                  *websocket.Conn
	uc                 *net.UDPConn
	sessionDescription *models.SessionDescription
	udpInfo            *models.UDPInfo
	running            bool
	playerRunning      bool
	udpReady           bool
	udpReadyWait       *sync.Cond
	speaking           bool
	frameData          chan []byte
	pausing            bool
	pauseWait          *sync.Cond
	playing            bool
	destroyed          bool
	skip               bool
	isWaitForExit      bool
	stopWaitForExit    chan any
}

// start the client if not started already
func StartClient(guildID, channelID string) error {
	clientList.RLock()
	var client = clientList.c[guildID]
	if client != nil {
		clientList.RUnlock()
		client.RLock()
		defer client.RUnlock()
		if (client.voiceChannelID != "") && (client.voiceChannelID != channelID) {
			return errors.New("Already in a different voice channel")
		} else if client.running {
			return nil
		}
	}
	clientList.RUnlock()

	client = addGuildToClient(guildID, channelID)
	sessionIdChan := make(chan string)
	voiceServerChan := make(chan *models.VoiceServer)
	err := gateway.JoinVoiceChannel(guildID, channelID, sessionIdChan, voiceServerChan)
	if err != nil {
		return err
	}

	// wait for session id and voice server
	go func() {
		sessionId, voiceServer := <-sessionIdChan, <-voiceServerChan
		client.Lock()
		defer client.Unlock()
		client.sessionID = &sessionId
		client.voiceServer = voiceServer
		go client.connect()
	}()

	return nil
}

// remove client from list and properly leave
func StopClient(guildID string) error {
	err := removeClient(guildID)
	if err != nil {
		return err
	}
	gateway.LeaveVoiceChannel(guildID)
	return nil
}

func (client *VoiceClient) connect() {
	// only connect once
	client.Lock()
	if client.running {
		client.Unlock()
		return
	}
	client.running = true
	client.Unlock()

	c, _, err := websocket.DefaultDialer.Dial("wss://"+client.voiceServer.Endpoint+"?v="+config.VOICE_GATEWAY_VERSION, nil)
	if err != nil {
		log.Fatalln(err)
	}
	client.c = c
	defer c.Close()

	ws_util.WriteJSONLog(client.c, models.NewVoiceIdentify(client.guildID, config.BOT_ID, *client.sessionID, client.voiceServer.Token), true)

	heatbeatInterval := make(chan int)
	heatbeatAcked := make(chan int64)

	// receive the gateway response
	go func() {
		for {
			var payload models.VoiceGatewayPayload
			err := c.ReadJSON(&payload)
			if err != nil {
				client.RLock()
				if !client.destroyed {
					log.Println(err)
				}
				client.RUnlock()
				return
			}
			if config.DEVELOPMENT {
				jsonDat, _ := json.Marshal(payload)
				log.Println("incoming voice: ", payload, string(jsonDat))
			} else {
				log.Println("incoming voice: ", payload)
			}
			switch payload.Op {
			case voice_opcode.Hello:
				heatbeatInterval <- int(payload.D.(map[string]any)["heartbeat_interval"].(float64))
			case voice_opcode.HeartbeatAck:
				n, _ := strconv.ParseInt(payload.D.(string), 10, 64)
				heatbeatAcked <- n
			case voice_opcode.Ready:
				var data models.UDPInfo
				mapstructure.Decode(payload.D, &data)
				if !util.ContainsStr(data.Modes, config.PREFERRED_MODE) {
					log.Panicln("Preferred mode not available")
				}
				client.udpInfo = &data
				// start the UDP connection
				go client.connectUdp()
				go client.sendSongFromFrameData()
			case voice_opcode.SessionDescription:
				var data models.SessionDescription
				mapstructure.Decode(payload.D, &data)
				client.sessionDescription = &data
				client.udpReadyWait.L.Lock()
				client.udpReady = true
				client.udpReadyWait.Signal()
				client.udpReadyWait.L.Unlock()
			}
		}
	}()

	// keeping heartbeats and prevent application from closing.
	interval := <-heatbeatInterval
	prevNonce := int64(0)
	for {
		time.Sleep(time.Duration(interval) * time.Millisecond)

		randNum, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			log.Println(err)
		}
		prevNonce = randNum.Int64()

		ws_util.WriteJSONLog(c, models.NewVoiceHeartbeat(prevNonce), true)

		select {
		case hbNonce := <-heatbeatAcked:
			if prevNonce != hbNonce {
				// handle nonce error
				log.Println(client.guildID, "Nonce invalid")
				StopClient(client.guildID)
				c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(4006, ""))
				return
			}
		case <-time.After(time.Duration(interval) * time.Millisecond):
			// uh oh timeout, goodbye
			log.Println(client.guildID, "Timeout")
			StopClient(client.guildID)
			c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(4009, ""))
			return
		case <-client.ctx.Done():
			log.Println(client.guildID, "Done")
			c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, ""))
			return
		}
	}
}

func (client *VoiceClient) connectUdp() {
	raddr, err := net.ResolveUDPAddr("udp", client.udpInfo.IP+":"+strconv.Itoa(int(client.udpInfo.Port)))
	if err != nil {
		log.Println(err)
		return
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Println(err)
		return
	}
	client.uc = conn
	ip, port := client.performIpDiscovery()
	ws_util.WriteJSONLog(client.c, models.NewVoiceSelectProtocol(ip, port, config.PREFERRED_MODE), true)

	go client.keepAlive(time.Second * 5)
}

func (client *VoiceClient) performIpDiscovery() (ip string, port uint16) {
	send := make([]byte, 70)

	binary.BigEndian.PutUint32(send, client.udpInfo.SSRC)
	_, err := client.uc.Write(send)
	if err != nil {
		log.Panicln(err)
	}

	receive := make([]byte, 70)
	rlen, _, err := client.uc.ReadFromUDP(receive)
	if err != nil {
		log.Panicln(err)
	}

	if rlen < 70 {
		log.Panicln("UDP packet too short")
	}

	for i := 4; i < 20; i++ {
		if receive[i] == 0 {
			break
		}
		ip += string(receive[i])
	}

	// Grab port from position 68 and 69
	port = binary.BigEndian.Uint16(receive[68:70])
	return
}

func (client *VoiceClient) keepAlive(i time.Duration) {
	defer client.uc.Close()

	var sequence uint64

	packet := make([]byte, 8)

	ticker := time.NewTicker(i)
	defer ticker.Stop()
	for {

		binary.LittleEndian.PutUint64(packet, sequence)
		sequence++

		_, err := client.uc.Write(packet)
		if err != nil {
			log.Println(err)
			client.ctxCancel()
			return
		}

		select {
		case <-client.ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (client *VoiceClient) sendSongFromFrameData() {
	client.udpReadyWait.L.Lock()
	for !client.udpReady {
		client.udpReadyWait.Wait()
	}
	client.udpReadyWait.L.Unlock()

	frameTime := uint32(config.DCA_FRAMERATE * config.DCA_FRAMEDURATION / 1000)
	ticker := time.NewTicker(time.Millisecond * time.Duration(config.DCA_FRAMEDURATION))

	client.RLock()
	if !client.speaking {
		ws_util.WriteJSONLog(client.c, models.NewVoiceSpeaking(client.udpInfo.SSRC), true)
	}
	client.RUnlock()

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
	var nonce [24]byte

	header[0] = 0x80
	header[1] = 0x78
	binary.BigEndian.PutUint32(header[8:], client.udpInfo.SSRC)

	for {
		binary.BigEndian.PutUint16(header[2:], uint16(sequenceNumber))
		binary.BigEndian.PutUint32(header[4:], timeStamp)

		copy(nonce[:], header)

		var frame []byte
		var ok bool
		select {
		case <-client.ctx.Done():
			return
		case frame, ok = <-client.frameData:
			if !ok {
				return
			}
		}

		packet := secretbox.Seal(header, frame, &nonce, &client.sessionDescription.SecretKey)

		select {
		case <-client.ctx.Done():
			return
		case <-ticker.C:
		}

		_, err = client.uc.Write(packet)

		if err != nil {
			log.Println(err)
		}

		if (sequenceNumber) == 0xFFFF {
			sequenceNumber = 0
		} else {
			sequenceNumber++
		}

		if (timeStamp + frameTime) >= 0xFFFFFFFF {
			timeStamp = 0
		} else {
			timeStamp += frameTime
		}
	}
}

func (client *VoiceClient) waitForExit() {
	client.RLock()
	client.isWaitForExit = true
	client.RUnlock()

	select {
	case <-client.stopWaitForExit:
		return
	case <-client.ctx.Done():
		return
	case <-time.After(config.IDLE_TIMEOUT):
		client.RLock()
		if client.running {
			client.RUnlock()
			log.Println(client.guildID, "Idle Timeout")
			StopClient(client.guildID)
		} else {
			client.RUnlock()
			return
		}
	}
}
