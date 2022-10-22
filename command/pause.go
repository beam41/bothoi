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

var commandPause = []discord_models.AppCommand{
	{
		Type:              app_command_type.ChatInput,
		Name:              "pause",
		Description:       "Pause/Resume the player",
		DefaultPermission: true,
	},
	{
		Type:              app_command_type.ChatInput,
		Name:              "resume",
		Description:       "Pause/Resume the player",
		DefaultPermission: true,
	},
}

func executePause(cm *commandManager, data *discord_models.Interaction) {
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
			"Pause error",
			fmt.Sprintf("<@%s> not in same voice channel as Bothoi", data.Member.User.Id),
			"error",
			embed_color.Error,
		)
		return
	}
	pausing, err := cm.voiceClientManager.PauseClient(data.GuildId)
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
			"Paused by the request of <@"+string(data.Member.User.Id)+">",
			"/resume to resume",
			embed_color.Default,
		)
	} else {
		response = util.BuildPlayerResponse(
			"Resumed",
			"Resumed by the request of <@"+string(data.Member.User.Id)+">",
			"",
			embed_color.Default,
		)
	}

}
