package queue

import (
	"bothoi/config"
	"bothoi/models"
	"bothoi/references/embed_color"
	"bothoi/util"
	"bothoi/util/http_util"
	"bothoi/voice"
	"fmt"
	"log"
	"strings"
)

func Execute(data *models.Interaction) {
	var response models.InteractionResponse
	// do response to interaction
	defer func() {
		url := config.InteractionResponseEndpoint
		url = strings.Replace(url, "<interaction_id>", data.ID, 1)
		url = strings.Replace(url, "<interaction_token>", data.Token, 1)

		_, err := http_util.PostJson(url, response)
		if err != nil {
			log.Println(err)
		}
	}()
	var playing, songQ = voice.GetSongQueue(data.GuildID, 0, 10)
	if songQ == nil || len(songQ) == 0 {
		response = util.BuildPlayerResponse(
			"No songs in queue",
			"Start playing a song now!",
			"Queue",
			embed_color.Error,
		)
		return
	}

	res := "Song in queue (Requested by)\n"

	if playing {
		res += fmt.Sprintf("**Currently Playing**\n%s (<@%s>)\n", songQ[0].Title, songQ[0].RequesterID)
		songQ = songQ[1:]
	}

	for i, song := range songQ {
		res += fmt.Sprintf("%d. %s (<@%s>)\n", i+1, song.Title, song.RequesterID)
	}

	response = util.BuildPlayerResponse(
		"Queue",
		res,
		fmt.Sprintf("%d %s in queue", len(songQ), util.Ternary(len(songQ) == 1, "song", "songs")),
		embed_color.Default,
	)
}
