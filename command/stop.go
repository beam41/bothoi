package command

import (
	"bothoi/config"
	"bothoi/models/discord_models"
	"bothoi/references/app_command_type"
	"bothoi/references/embed_color"
	"bothoi/repo"
	"bothoi/util"
	"bothoi/util/http_util"
	"fmt"
	"log"
	"strings"
)

var commandStop = discord_models.AppCommand{
	Type:              app_command_type.ChatInput,
	Name:              "stop",
	Description:       "Stop the player and leave the voice channel",
	DefaultPermission: true,
}

func executeStop(cm *commandManager, data *discord_models.Interaction) {
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
	userVoiceState := repo.GetVoiceState(data.Member.User.Id)
	clientVoiceChannel := cm.voiceClientManager.GetVoiceChannelId(data.GuildId)
	if userVoiceState == nil || *userVoiceState.ChannelId != clientVoiceChannel {
		response = util.BuildPlayerResponse(
			"Cannot stop",
			fmt.Sprintf("<@%s> not in same voice channel as Bothoi", data.Member.User.Id),
			"Error",
			embed_color.Error,
		)
		return
	}
	err := cm.voiceClientManager.StopClient(data.GuildId)
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
		"Stopped by the request of <@"+string(data.Member.User.Id)+">",
		"goodbye",
		embed_color.Default,
	)
}
