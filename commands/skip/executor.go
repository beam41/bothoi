package skip

import (
	"bothoi/config"
	"bothoi/models/discord_models"
	"bothoi/references/embed_color"
	"bothoi/repo"
	"bothoi/util"
	"bothoi/util/http_util"
	"bothoi/voice"
	"fmt"
	"log"
	"strings"
)

func Execute(data *discord_models.Interaction) {
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
	clientVoiceChannel := voice.GetVoiceChannelId(data.GuildId)
	if userVoiceState == nil || *userVoiceState.ChannelId != clientVoiceChannel {
		response = util.BuildPlayerResponse(
			"Skip error",
			fmt.Sprintf("<@%s> not in same voice channel as Bothoi", data.Member.User.Id),
			"error",
			embed_color.Error,
		)
		return
	}
	err := voice.SkipSong(data.GuildId)
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
