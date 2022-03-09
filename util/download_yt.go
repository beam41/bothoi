package util

import (
	"bytes"
	"log"

	"github.com/kkdai/youtube/v2"
)

func DownloadYt(ytID string) ([]byte, error) {
	ytClient := youtube.Client{}
	video, err := ytClient.GetVideo(ytID)
	if err != nil {
		log.Panicln(err)
	}

	formats := video.Formats.WithAudioChannels()
	format, err := GetBestAudioYT(formats)
	if err != nil {
		return nil, err
	}

	stream, _, err := ytClient.GetStream(video, format)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes(), nil
}
