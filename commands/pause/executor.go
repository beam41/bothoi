package pause

import (
	"bothoi/config"
	"bothoi/models"
	"bothoi/references/embed_color"
	"bothoi/repo"
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
	userVoiceState := repo.GetVoiceState(data.Member.User.ID)
	clientVoiceChannel := voice.GetVoiceChannelID(data.GuildID)
	if userVoiceState == nil || userVoiceState.ChannelID != clientVoiceChannel {
		response = util.BuildPlayerResponse(
			"Pause error",
			fmt.Sprintf("<@%s> not in same voice channel as Bothoi", data.Member.User.ID),
			"error",
			embed_color.Error,
		)
		return
	}
	pausing, err := voice.PauseClient(data.GuildID)
	if err != nil {
		response = util.BuildPlayerResponse(
			"Pause error",
			"Cannot be paused/resumed",
			"error",
			embed_color.Error,
		)
		return
	}
	if pausing {
		response = util.BuildPlayerResponse(
			"Paused",
			"Paused by the request of <@"+data.Member.User.ID+">",
			"/resume to resume",
			embed_color.Default,
		)
	} else {
		response = util.BuildPlayerResponse(
			"Resumed",
			"Resumed by the request of <@"+data.Member.User.ID+">",
			"",
			embed_color.Default,
		)
	}

}
