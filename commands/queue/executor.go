package queue

import (
	"bothoi/config"
	"bothoi/models"
	"bothoi/references/embed_color"
	"bothoi/states"
	"bothoi/util"
	"bothoi/util/http_util"
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

func Execute(data *models.Interaction, _ *websocket.Conn) {
	var response models.InteractionResponse
	// do response to interaction
	defer func() {
		url := config.INTERACTION_RESPONSE_ENDPOINT
		url = strings.Replace(url, "<interaction_id>", data.ID, 1)
		url = strings.Replace(url, "<interaction_token>", data.Token, 1)

		_, err := http_util.PostJson(url, response)
		if err != nil {
			log.Println(err)
		}
	}()
	var songQ = states.SongQueue[data.GuildID]
	if songQ == nil {
		response = util.BuildPlayerResponse(
			"No songs in queue",
			"Start playing a song now!",
			"Queue",
			embed_color.Error,
		)
		return
	}

	res := "Song in queue (Requested by)\n"
	for i, song := range songQ.Songs {
		res += fmt.Sprintf("%d. %s (<@%s>)\n", i+1, song.Title, song.RequesterID)
	}

	response = util.BuildPlayerResponse(
		"Queue",
		res,
		fmt.Sprintf("%d %s in queue", len(songQ.Songs), util.Ternary(len(songQ.Songs) == 1, "song", "songs")),
		embed_color.Default,
	)
}
