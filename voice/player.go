package voice

import (
	"bothoi/config"
	"bothoi/models"
	"bothoi/util"
	"bytes"
	"io"
	"log"

	"github.com/jonas747/dca/v2"
)

// play song if not start already
func (client *VoiceClient) play() {
	// ensure only one player is running at a time
	client.Lock()
	if client.playing {
		client.Unlock()
		return
	}
	client.playing = true
	client.Unlock()

	// encode settings
	options := dca.StdEncodeOptions
	options.FrameRate = config.DCA_FRAMERATE
	options.FrameDuration = config.DCA_FRAMEDURATION
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "audio"

	// run player
	for {
		currentSong := client.songQueue[0]
		currentSong.DownloadLock.Lock()
		if currentSong.SongData == nil {
			sd, err := util.DownloadYt(currentSong.YtID)
			currentSong.SongData = sd
			if err != nil {
				currentSong.DownloadLock.Unlock()
				log.Println("Can't play this song: ", err)

				client.Lock()
				if len(client.songQueue) > 1 {
					client.songQueue = client.songQueue[1:]
					client.Unlock()
					client.downloadUpcoming()
				} else {
					// stop player
					client.songQueue = []*models.SongItemWData{}
					client.playing = false
					client.Unlock()
					go client.waitForExit()
					return
				}
				continue
			}
		}
		currentSong.DownloadLock.Unlock()

		reader := bytes.NewReader(currentSong.SongData)
		log.Println(len(currentSong.SongData))
		client.encodeSong(reader, options)

		client.Lock()
		if len(client.songQueue) > 1 {
			client.songQueue = client.songQueue[1:]
			client.Unlock()
			client.downloadUpcoming()
		} else {
			// stop player
			client.songQueue = []*models.SongItemWData{}
			client.playing = false
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
		client.frameData <- frame

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
		go func() {
			song1.DownloadLock.Lock()
			if song1.SongData == nil {
				sd, _ := util.DownloadYt(song1.YtID)
				song1.SongData = sd
			}
			song1.DownloadLock.Unlock()
		}()
	}
}
