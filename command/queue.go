package command

import (
	"bothoi/bh_context"
	"bothoi/config"
	"bothoi/models/discord_models"
	"bothoi/references/app_command_type"
	"bothoi/references/embed_color"
	"bothoi/util"
	"bothoi/util/http_util"
	"fmt"
	"log"
	"strings"
)

var commandQueue = discord_models.AppCommand{
	Type:              app_command_type.ChatInput,
	Name:              "queue",
	Description:       "List the music player queue",
	DefaultPermission: true,
}

func executeQueue(cm *commandManager, data *discord_models.Interaction) {
	var response discord_models.InteractionResponse
	// do response to interaction
	defer func() {
		url := config.InteractionResponseEndpoint
		url = strings.Replace(url, "<interaction_id>", string(data.Id), 1)
		url = strings.Replace(url, "<interaction_token>", data.Token, 1)

		_, err := http_util.PostJson(url, response)
		if err != nil {
			log.Println(err)
		}
	}()
	var playing, songQ = bh_context.Ctx.VoiceClientManager.GetSongQueue(data.GuildId, 0, 10)
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
		res += fmt.Sprintf("**Currently Playing**\n%s (<@%s>)\n", songQ[0].Title, songQ[0].RequesterId)
		songQ = songQ[1:]
	}

	for i, song := range songQ {
		res += fmt.Sprintf("%d. %s (<@%s>)\n", i+1, song.Title, song.RequesterId)
	}

	response = util.BuildPlayerResponse(
		"Queue",
		res,
		fmt.Sprintf("%d %s in queue", len(songQ), util.Ternary(len(songQ) == 1, "song", "songs")),
		embed_color.Default,
	)
}
