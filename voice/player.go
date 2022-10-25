package voice

import (
	"bothoi/config"
	"bothoi/models/discord_models"
	"bothoi/repo"
	"bothoi/util/yt_util.go"
	"crypto/rand"
	"encoding/binary"
	"golang.org/x/crypto/nacl/secretbox"
	"io"
	"log"
	"math"
	"math/big"
	"time"

	"github.com/jonas747/dca/v2"
)

// play song if not start already
func (client *client) play() {
	// ensure only one player is running at a time
	client.Lock()
	if client.playerRunning {
		client.Unlock()
		return
	}
	client.playerRunning = true
	client.Unlock()

	// encode settings
	options := dca.StdEncodeOptions
	options.Volume = 128
	options.FrameRate = config.DcaFramerate
	options.FrameDuration = config.DcaFrameduration
	options.RawOutput = true
	options.Bitrate = 96
	options.BufferedFrames = config.DcaBufferedFrame
	options.Application = dca.AudioApplicationAudio

	// run player
	for {
		currentSong, err := repo.GetNextSong(client.guildId)
		if currentSong == nil {
			// stop player
			client.Lock()
			client.playerRunning = false
			client.Unlock()
			go client.waitForExit()
			return
		} else if err != nil {
			log.Println(err)
			continue
		} else if currentSong.Playing {
			_ = repo.DeleteSong(currentSong.Id)
			continue
		}
		log.Println(client.guildId, "Playing song: ", currentSong.Title)
		url, err := yt_util.GetYoutubeDownloadUrl(currentSong.YtId)
		if err != nil {
			log.Println(err)
			continue
		}

		client.Lock()
		client.playing = true
		_ = repo.UpdateSongPlaying(currentSong.Id)
		if client.isWaitForExit {
			client.isWaitForExit = false
			client.stopWaitForExit <- struct{}{}
		}
		client.Unlock()

		options.StartTime = int(currentSong.Seek)
		client.playTillComplete(url, options)

		client.Lock()
		if client.destroyed {
			client.Unlock()
			return
		}
		client.skip = false
		client.playing = false
		client.Unlock()
		_ = repo.DeleteSong(currentSong.Id)
	}
}

func (client *client) playTillComplete(url string, options *dca.EncodeOptions) {
	encodeSession, err := dca.EncodeFile(url, options)
	if err != nil {
		return
	}
	defer encodeSession.Cleanup()
	for {
		if client.sendSong(encodeSession) {
			return
		}
	}
}

func (client *client) sendSong(encodeSession *dca.EncodeSession) bool {
	client.udpReadyWait.L.Lock()
	for !client.udpReady {
		client.udpReadyWait.Wait()
	}
	client.udpReadyWait.L.Unlock()

	frameTime := uint32(config.DcaFramerate * config.DcaFrameduration / 1000)
	ticker := time.NewTicker(time.Millisecond * time.Duration(config.DcaFrameduration))

	client.RLock()
	if !client.speaking {
		err := client.connWriteJSON(discord_models.NewVoiceSpeaking(client.udpInfo.SSRC))
		if err != nil {
			log.Println(err)
		}
	}
	client.RUnlock()

	randNum, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		log.Println(err)
	}
	sequenceNumber := uint16(randNum.Uint64())

	randNum, err = rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		log.Println(err)
	}
	timeStamp := uint32(randNum.Uint64())

	header := make([]byte, 12)
	var nonce [24]byte

	header[0] = 0x80
	header[1] = 0x78
	binary.BigEndian.PutUint32(header[8:], client.udpInfo.SSRC)

	if err != nil {
		log.Panicln(err)
	}

	for {
		client.RLock()
		if client.destroyed || client.skip {
			client.RUnlock()
			return true
		}
		client.RUnlock()

		client.pauseWait.L.Lock()
		for client.pausing {
			client.pauseWait.Wait()
		}
		client.pauseWait.L.Unlock()

		frame, err := encodeSession.OpusFrame()
		if err != nil {
			if err != io.EOF {
				log.Println("encodeSession error", err)
			}
			return true
		}

		binary.BigEndian.PutUint16(header[2:], sequenceNumber)
		binary.BigEndian.PutUint32(header[4:], timeStamp)

		copy(nonce[:], header)

		packet := secretbox.Seal(header, frame, &nonce, &client.sessionDescription.SecretKey)

		select {
		case <-client.vcCtx.Done():
			return false
		case <-client.ctx.Done():
			return true
		case <-ticker.C:
		}

		_, err = client.uc.Write(packet)

		if err != nil {
			log.Println(client.guildId, "player udp err", err)
			client.voiceRestart()
			return false

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
