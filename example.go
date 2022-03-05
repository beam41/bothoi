package main

import (
	"fmt"
	"io"
	"os"

	"github.com/jonas747/dca/v2"
	"github.com/kkdai/youtube/v2"
)

func ExampleClient() {
	videoID := "JMmUW4d3Noc"
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		panic(err)
	}

	formats := video.Formats.WithAudioChannels() // only get videos with audio
	for _, format := range formats {
		fmt.Println(format.AudioSampleRate, format.AudioQuality, format.Bitrate, format.QualityLabel, format.AudioChannels, format.Width, format.Height)
	}
	x, err := client.GetStreamURL(video, &formats[3])

	encodeSession, err := dca.EncodeFile(x, dca.StdEncodeOptions)
	// Make sure everything is cleaned up, that for example the encoding process if any issues happened isnt lingering around
	defer encodeSession.Cleanup()
	
	output, err := os.Create("output.dca")
	if err != nil {
		// Handle the error
	}
	io.Copy(output, encodeSession)
}
