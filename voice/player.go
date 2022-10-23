package voice

import (
	"bothoi/config"
	"bothoi/models"
	"bothoi/util/yt_util.go"
	"io"
	"log"

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
	if client.isWaitForExit {
		client.isWaitForExit = false
		client.stopWaitForExit <- struct{}{}
	}
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
		currentSong := client.songQueue[0]
		log.Println(client.guildId, "Playing song: ", currentSong.Title)
		url, err := yt_util.GetYoutubeDownloadUrl(currentSong.YtId)
		if err != nil {
			log.Println(err)
			continue
		}

		client.Lock()
		client.playing = true
		client.Unlock()

		client.encodeSong(url, options)

		client.Lock()
		if client.destroyed {
			client.Unlock()
			return
		}
		client.skip = false
		client.playing = false
		if len(client.songQueue) > 1 {
			client.songQueue = client.songQueue[1:]
			client.Unlock()
		} else {
			// stop player
			client.songQueue = []*models.SongItem{}
			client.playerRunning = false
			client.Unlock()
			go client.waitForExit()
			return
		}
	}
}

func (client *client) encodeSong(url string, options *dca.EncodeOptions) {
	encodeSession, err := dca.EncodeFile(url, options)
	defer encodeSession.Cleanup()

	if err != nil {
		log.Panicln(err)
	}

	for {
		frame, err := encodeSession.OpusFrame()
		if err != nil {
			if err != io.EOF {
				log.Println("encodeSession error", err)
			}
			return
		}

		client.RLock()
		if client.destroyed || client.skip {
			client.RUnlock()
			return
		}
		// if we unlock before send, it might send a frame after the client has been destroyed
		client.frameData <- frame
		client.RUnlock()

		client.pauseWait.L.Lock()
		for client.pausing {
			client.pauseWait.Wait()
		}
		client.pauseWait.L.Unlock()
	}
}
