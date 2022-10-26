package command

import (
	"bothoi/bh_context"
	"bothoi/config"
	"bothoi/models/discord_models"
	"bothoi/references/embed_color"
	"bothoi/repo"
	"bothoi/util"
	"bothoi/util/http_util"
	"bothoi/util/yt_util.go"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const commandPlay = "play"

func executePlay(data *discord_models.Interaction) {
	options := util.SliceToMap(data.Data.Options, func(i int, item discord_models.InteractionOption) string { return item.Name })

	// post waiting prevent response timeout
	url := config.InteractionResponseEndpoint
	url = strings.Replace(url, "<interaction_id>", strconv.FormatUint(uint64(data.ID), 10), 1)
	url = strings.Replace(url, "<interaction_token>", data.Token, 1)

	_, err := http_util.PostJson(url, util.BuildPlayerResponse(
		"Play a song",
		"Loading...",
		"please wait",
		embed_color.Playing,
	))
	if err != nil {
		log.Println(err)
	}

	var response discord_models.InteractionCallbackData
	// do response to interaction
	defer func() {
		url := config.InteractionResponseEditEndpoint
		url = strings.Replace(url, "<application_id>", strconv.FormatUint(uint64(config.BotID), 10), 1)
		url = strings.Replace(url, "<interaction_token>", data.Token, 1)

		_, err := http_util.PatchJson(url, response)
		if err != nil {
			log.Println(err)
		}
	}()

	userVoiceChannel := repo.GetChannelIDByUserIDAndGuildID(data.Member.User.ID, data.GuildID)
	if userVoiceChannel == nil {
		response = util.BuildPlayerResponseData(
			"Can't play a song :(",
			fmt.Sprintf("<@%d> not in voice channel", data.Member.User.ID),
			"Error",
			embed_color.Error,
		)
		return
	}

	title, ytID, duration, noResult, _ := yt_util.SearchYt(options["song"].Value.(string))
	if noResult {
		log.Println(err)
		response = util.BuildPlayerResponseData(
			"Can't play a song :(",
			"Song not found",
			"Error",
			embed_color.Error,
		)
		return
	}

	clientVoiceChannel := repo.GetChannelIDByUserIDAndGuildID(data.Member.User.ID, config.BotID)
	if userVoiceChannel == nil || (clientVoiceChannel != nil && *userVoiceChannel != *clientVoiceChannel) {
		response = util.BuildPlayerResponseData(
			"Can't play a song :(",
			fmt.Sprintf("<@%d> not in the same voice channel as bot", data.Member.User.ID),
			"Error",
			embed_color.Error,
		)
		return
	}
	seek := uint32(0)
	if s, ok := options["seek"]; ok {
		seek = util.ConvertVidLengthToSeconds(s.Value.(string))
	}
	err = repo.AddSongToQueue(data.GuildID, data.Member.User.ID, ytID, title, duration, seek)
	if err != nil {
		log.Println(err)
		response = util.BuildPlayerResponseData(
			"Can't play a song :(",
			"Unknown Error",
			"Error",
			embed_color.Error,
		)
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = bh_context.GetVoiceClientManager().StartClient(data.GuildID, *userVoiceChannel)
	if err != nil {
		log.Println(err)
		response = util.BuildPlayerResponseData(
			"Can't play a song :(",
			"Unknown Error",
			"Error",
			embed_color.Error,
		)
		if err != nil {
			log.Println(err)
		}
		return
	}

	queueSize := repo.GetQueueSize(data.GuildID)
	if queueSize == 1 {
		response = util.BuildPlayerResponseData(
			"Play a song",
			fmt.Sprintf("Playing %s\n%s\nrequested by <@%d>", title, util.ConvertSecondsToVidLength(duration), data.Member.User.ID),
			"Playing",
			embed_color.Playing,
		)
	} else {
		response = util.BuildPlayerResponseData(
			"Play a song",
			fmt.Sprintf("Added %s\n%s\nrequested by <@%d>", title, util.ConvertSecondsToVidLength(duration), data.Member.User.ID),
			fmt.Sprintf("#%d in queue", queueSize),
			embed_color.Playing,
		)
	}
}
