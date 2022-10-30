package voice

import (
	"bothoi/config"
	"bothoi/models/discord_models"
	"bothoi/models/types"
	"bothoi/references/voice_opcode"
	"bothoi/util"
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"log"
	"math"
	"math/big"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

type client struct {
	sync.RWMutex
	ctx                context.Context
	ctxCancel          context.CancelFunc
	guildID            types.Snowflake
	sessionID          *string
	voiceServer        *discord_models.VoiceServer
	c                  *websocket.Conn
	uc                 *net.UDPConn
	sessionDescription *discord_models.SessionDescription
	udpInfo            *discord_models.UDPInfo
	running            bool
	playerRunning      bool
	udpReady           bool
	udpReadyWait       *sync.Cond
	speaking           bool
	pausing            bool
	pauseWait          *sync.Cond
	playing            bool
	destroyed          bool
	skip               bool
	isWaitForExit      bool
	stopWaitForExit    chan struct{}
	clm                *ClientManager
	vcCtx              context.Context
	vcCtxCancel        context.CancelFunc
	resume             bool
}

func (client *client) connWriteJSON(v any) (err error) {
	log.Println(client.guildID, "outgoing voice: ", v)
	err = client.c.WriteJSON(v)
	return
}

func (client *client) connCloseNormal() {
	err := client.c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println(client.guildID, err)
	}
}

func (client *client) connectionRestart(resume bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(client.guildID, "voiceRestart panic occurred:", err)
		}
	}()
	client.Lock()
	client.udpReadyWait.L.Lock()
	client.udpReady = false
	client.udpReadyWait.Signal()
	client.udpReadyWait.L.Unlock()
	client.resume = resume
	client.speaking = false
	client.Unlock()
	client.vcCtxCancel()
	err := client.c.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(util.Ternary(resume, 4006, websocket.CloseTryAgainLater), ""),
	)
	if err != nil {
		log.Println(client.guildID, err)
	}
}

func (client *client) connect() {
	for {
		if client.destroyed {
			return
		}
		ctx, cancel := context.WithCancel(context.Background())
		client.vcCtx = ctx
		client.vcCtxCancel = cancel
		client.connection()
	}
}

func (client *client) connection() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(client.guildID, "voice connection panic occurred:", err)
		}
	}()
	// only connect once
	client.Lock()
	if client.running {
		client.Unlock()
		return
	}
	client.running = true
	client.Unlock()
	defer func() {
		client.Lock()
		client.running = false
		client.Unlock()
	}()

	c, _, err := websocket.DefaultDialer.Dial("wss://"+client.voiceServer.Endpoint+"?v="+config.VoiceGatewayVersion, nil)
	if err != nil {
		log.Panicln(client.guildID, err)
	}
	client.c = c
	defer func() {
		err := c.Close()
		if err != nil {
			log.Println(client.guildID, err)
		}
	}()

	if !client.resume {
		err = client.connWriteJSON(discord_models.NewVoiceIdentify(client.guildID, config.BotID, *client.sessionID, client.voiceServer.Token))
		if err != nil {
			log.Println(client.guildID, err)
		}
	} else {
		err = client.connWriteJSON(discord_models.NewVoiceResume(client.guildID, *client.sessionID, client.voiceServer.Token))
		if err != nil {
			log.Println(client.guildID, err)
		}
	}

	heartbeatInterval := make(chan int)
	heartbeatAcked := make(chan int64)

	// receive the gateway response
	go func() {
		for {
			var payload discord_models.VoiceGatewayPayload
			err := c.ReadJSON(&payload)
			if err != nil {
				client.RLock()
				if client.destroyed {
					client.RUnlock()
					return
				}
				client.RUnlock()
				log.Println(client.guildID, "read err", err)
				if websocket.IsUnexpectedCloseError(err, 4015) {
					client.connectionRestart(false)
				} else {
					client.connectionRestart(true)
				}
				return
			}
			if config.Development {
				jsonDat, _ := json.Marshal(payload)
				log.Println(client.guildID, "incoming voice: ", payload, string(jsonDat))
			} else {
				log.Println(client.guildID, "incoming voice: ", payload)
			}
			switch payload.Op {
			case voice_opcode.Hello:
				heartbeatInterval <- int(payload.D.(map[string]any)["heartbeat_interval"].(float64))
			case voice_opcode.HeartbeatAck:
				n, _ := strconv.ParseInt(payload.D.(string), 10, 64)
				heartbeatAcked <- n
			case voice_opcode.Ready:
				var data discord_models.UDPInfo
				err := mapstructure.WeakDecode(payload.D, &data)
				if err != nil {
					log.Println(client.guildID, err)
					continue
				}
				if !util.ContainsStr(data.Modes, config.PreferredMode) {
					log.Panicln(client.guildID, "Preferred mode not available")
				}
				client.udpInfo = &data
				// start the UDP connection
				go client.udpConnect()
			case voice_opcode.SessionDescription:
				var data discord_models.SessionDescription
				err := mapstructure.WeakDecode(payload.D, &data)
				if err != nil {
					log.Println(client.guildID, err)
					continue
				}
				client.sessionDescription = &data
			case voice_opcode.Resumed:
				log.Println(client.guildID, "Resumed")
				client.udpConnect()
			}
		}
	}()

	// keeping heartbeats and prevent application from closing.
	interval := <-heartbeatInterval
	sentNonce := int64(0)

	heartbeatIntervalTicker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	defer heartbeatIntervalTicker.Stop()

	for {
		select {
		case <-heartbeatIntervalTicker.C:
		case <-client.ctx.Done():
			return
		case <-client.vcCtx.Done():
			return
		}

		randNum, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			log.Println(client.guildID, err)
		}
		sentNonce = randNum.Int64()

		err = client.connWriteJSON(discord_models.NewVoiceHeartbeat(sentNonce))
		if err != nil {
			log.Println(client.guildID, err)
		}

		select {
		case hbNonce := <-heartbeatAcked:
			if sentNonce != hbNonce {
				// handle nonce error
				log.Println(client.guildID, "Nonce invalid")
				client.connectionRestart(false)
				return
			}
		case <-heartbeatIntervalTicker.C:
			log.Println(client.guildID, "Timeout")
			client.connectionRestart(true)
			return
		case <-client.ctx.Done():
			return
		case <-client.vcCtx.Done():
			return
		}
	}
}

func (client *client) udpConnect() {
	raddr, err := net.ResolveUDPAddr("udp", client.udpInfo.IP+":"+strconv.Itoa(int(client.udpInfo.Port)))
	if err != nil {
		log.Println(client.guildID, err)
		return
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Println(client.guildID, err)
		return
	}
	client.uc = conn
	ip, port := client.udpIpDiscovery()
	err = client.connWriteJSON(discord_models.NewVoiceSelectProtocol(ip, port, config.PreferredMode))
	if err != nil {
		log.Println(client.guildID, err)
	}
	client.udpReadyWait.L.Lock()
	client.udpReady = true
	client.udpReadyWait.Signal()
	client.udpReadyWait.L.Unlock()

	go client.udpKeepAlive(time.Second * 5)
}

func (client *client) udpIpDiscovery() (ip string, port uint16) {
	send := make([]byte, 70)

	binary.BigEndian.PutUint32(send, client.udpInfo.SSRC)
	_, err := client.uc.Write(send)
	if err != nil {
		log.Panicln(client.guildID, err)
	}

	receive := make([]byte, 70)
	rlen, _, err := client.uc.ReadFromUDP(receive)
	if err != nil {
		log.Panicln(client.guildID, err)
	}

	if rlen < 70 {
		log.Panicln(client.guildID, "UDP packet too short")
	}

	var ipb strings.Builder
	for i := 4; i < 20; i++ {
		if receive[i] == 0 {
			break
		}
		ipb.WriteString(string(receive[i]))
	}
	ip = ipb.String()

	// Grab port from position 68 and 69
	port = binary.BigEndian.Uint16(receive[68:70])
	return
}

func (client *client) udpKeepAlive(i time.Duration) {
	defer func(uc *net.UDPConn) {
		err := uc.Close()
		if err != nil {
			log.Println(client.guildID, err)
		}
	}(client.uc)

	var sequence uint64

	packet := make([]byte, 8)

	ticker := time.NewTicker(i)
	defer ticker.Stop()
	for {
		binary.LittleEndian.PutUint64(packet, sequence)
		sequence++

		_, err := client.uc.Write(packet)
		if err != nil {
			log.Println(client.guildID, "Udp err", err)
			client.connectionRestart(true)
			return
		}

		select {
		case <-client.vcCtx.Done():
			return
		case <-client.ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (client *client) waitForExit() {
	client.Lock()
	if client.isWaitForExit == true {
		client.stopWaitForExit <- struct{}{}
	}
	client.isWaitForExit = true
	client.Unlock()

	idleTimeoutCountdown := time.NewTimer(config.IdleTimeout)
	defer func() {
		// clean timer just in case
		if !idleTimeoutCountdown.Stop() {
			select {
			case <-idleTimeoutCountdown.C:
			default:
			}
		}
	}()

	select {
	case <-client.stopWaitForExit:
		return
	case <-client.ctx.Done():
		return
	case <-idleTimeoutCountdown.C:
		client.RLock()
		if client.running {
			client.RUnlock()
			log.Println(client.guildID, "Idle Timeout")
			err := client.clm.ClientStop(client.guildID)
			if err != nil {
				log.Println(client.guildID, "Cannot stop", err)
			}
		} else {
			client.RUnlock()
			return
		}
	}
}
