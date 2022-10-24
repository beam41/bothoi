package command

import (
	"bothoi/config"
	"bothoi/models/discord_models"
	"bothoi/references/embed_color"
	"bothoi/repo"
	"bothoi/util"
	"bothoi/util/http_util"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const commandQueue = "queue"

func executeQueue(data *discord_models.Interaction) {
	var response discord_models.InteractionResponse
	// do response to interaction
	defer func() {
		url := config.InteractionResponseEndpoint
		url = strings.Replace(url, "<interaction_id>", strconv.FormatUint(uint64(data.Id), 10), 1)
		url = strings.Replace(url, "<interaction_token>", data.Token, 1)

		_, err := http_util.PostJson(url, response)
		if err != nil {
			log.Println(err)
		}
	}()
	var songQ = repo.GetSongQueue(data.GuildId, 10)
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

	if songQ[0].Playing {
		res += fmt.Sprintf("**Currently Playing**\n%s (<@%d>)\n", songQ[0].Title, songQ[0].RequesterId)
		songQ = songQ[1:]
	}

	for i, song := range songQ {
		res += fmt.Sprintf("%d. %s (<@%d>)\n", i+1, song.Title, song.RequesterId)
	}

	response = util.BuildPlayerResponse(
		"Queue",
		res,
		fmt.Sprintf("%d song%s in queue", len(songQ), util.Ternary(len(songQ) == 1, "", "s")),
		embed_color.Default,
	)
}
