package voice

import (
	"bothoi/config"
	"bothoi/models"
	"bothoi/util/yt_util.go"
	"bytes"
	"io"
	"log"

	"github.com/jonas747/dca/v2"
)

// play song if not start already
func (client *VoiceClient) play() {
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
	options.FrameRate = config.DCA_FRAMERATE
	options.FrameDuration = config.DCA_FRAMEDURATION
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "audio"

	// run player
	for {
		currentSong := client.songQueue[0]
		log.Println(client.guildID, "Playing song: ", currentSong.Title)
		currentSong.DownloadLock.Lock()
		if currentSong.SongData == nil {
			sd, err := yt_util.DownloadYt(currentSong.YtID)
			currentSong.SongData = sd
			if err != nil {
				currentSong.DownloadLock.Unlock()
				log.Println(client.guildID, "Can't play this song: ", err)

				client.Lock()
				if len(client.songQueue) > 1 {
					client.songQueue = client.songQueue[1:]
					client.Unlock()
					client.downloadUpcoming()
				} else {
					// stop player
					client.songQueue = []*models.SongItemWData{}
					client.playerRunning = false
					client.Unlock()
					go client.waitForExit()
					return
				}
				continue
			}
		}
		currentSong.DownloadLock.Unlock()

		client.Lock()
		client.playing = true
		client.Unlock()

		reader := bytes.NewReader(currentSong.SongData)
		client.encodeSong(reader, options)

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
			client.downloadUpcoming()
		} else {
			// stop player
			client.songQueue = []*models.SongItemWData{}
			client.playerRunning = false
			client.Unlock()
			go client.waitForExit()
			return
		}
	}
}

func (client *VoiceClient) encodeSong(reader io.Reader, options *dca.EncodeOptions) {
	encodeSession, err := dca.EncodeMem(reader, options)
	defer encodeSession.Cleanup()

	if err != nil {
		log.Panicln(err)
	}

	for {
		frame, err := encodeSession.OpusFrame()
		if err != nil {
			if err == io.EOF {
				return
			} else {
				log.Panicln(err)
				return
			}
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

func (client *VoiceClient) downloadUpcoming() {
	client.Lock()
	defer client.Unlock()
	if len(client.songQueue) >= 2 {
		song1 := client.songQueue[1]
		if song1.Downloading {
			// prevent multiple goroutine waiting on lock while download
			return
		}
		song1.Downloading = true

		go func() {
			song1.DownloadLock.Lock()
			if song1.SongData == nil {
				sd, _ := yt_util.DownloadYt(song1.YtID)
				song1.SongData = sd
			}
			song1.DownloadLock.Unlock()
		}()
	}
}
