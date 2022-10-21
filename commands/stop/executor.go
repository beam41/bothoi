package stop

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
			"Cannot stop",
			fmt.Sprintf("<@%s> not in same voice channel as Bothoi", data.Member.User.ID),
			"Error",
			embed_color.Error,
		)
		return
	}
	err := voice.StopClient(data.GuildID)
	if err != nil {
		response = util.BuildPlayerResponse(
			"Stopped",
			"Cannot be stopped",
			"error",
			embed_color.Error,
		)
		return
	}
	response = util.BuildPlayerResponse(
		"Stopped",
		"Stopped by the request of <@"+data.Member.User.ID+">",
		"goodbye",
		embed_color.Default,
	)
}
