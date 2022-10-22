package command

import (
	"bothoi/bh_context"
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

var commandSkip = discord_models.AppCommand{
	Type:              app_command_type.ChatInput,
	Name:              "skip",
	Description:       "Skip song in the player",
	DefaultPermission: true,
}

func executeSkip(data *discord_models.Interaction) {
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
	clientVoiceChannel := bh_context.GetVoiceClientManager().GetVoiceChannelId(data.GuildId)
	if userVoiceState == nil || *userVoiceState.ChannelId != clientVoiceChannel {
		response = util.BuildPlayerResponse(
			"Skip error",
			fmt.Sprintf("<@%s> not in same voice channel as Bothoi", data.Member.User.Id),
			"error",
			embed_color.Error,
		)
		return
	}
	err := bh_context.GetVoiceClientManager().SkipSong(data.GuildId)
	if err != nil {
		response = util.BuildPlayerResponse(
			"Skip error",
			"Cannot skip",
			"error",
			embed_color.Error,
		)
		return
	}
	response = util.BuildPlayerResponse(
		"Skipped",
		"Skipped by the request of <@"+string(data.Member.User.Id)+">",
		"skipped",
		embed_color.Default,
	)
}
