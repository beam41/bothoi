package command

import (
	"bothoi/bh_context"
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

const commandSkip = "skip"

func executeSkip(data *discord_models.Interaction) {
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
	userVoiceChannel := repo.GetChannelIdByUserIdAndGuildId(data.Member.User.Id, data.GuildId)
	clientVoiceChannel := repo.GetChannelIdByUserIdAndGuildId(config.BotId, data.GuildId)
	if userVoiceChannel == nil || (clientVoiceChannel == nil && *userVoiceChannel != *clientVoiceChannel) {
		response = util.BuildPlayerResponse(
			"Skip error",
			fmt.Sprintf("<@%d> not in same voice channel as Bothoi", data.Member.User.Id),
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
		fmt.Sprintf("Skipped by the request of <@%d>", uint64(data.Member.User.Id)),
		"skipped",
		embed_color.Default,
	)
}
