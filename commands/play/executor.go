package play

import (
	"bothoi/config"
	"bothoi/models/discord_models"
	"bothoi/references/embed_color"
	"bothoi/repo"
	"bothoi/util"
	"bothoi/util/http_util"
	"bothoi/util/yt_util.go"
	"bothoi/voice"
	"fmt"
	"log"
	"strings"
)

func Execute(data *discord_models.Interaction) {
	options := util.MapInteractionOption(data.Data.Options)
	userVoiceState := repo.GetVoiceState(data.Member.User.Id)

	// post waiting prevent response timeout
	url := config.InteractionResponseEndpoint
	url = strings.Replace(url, "<interaction_id>", string(data.Id), 1)
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
		url := config.InteractionResponseEndpoint
		url = strings.Replace(url, "<application_id>", config.BotId, 1)
		url = strings.Replace(url, "<interaction_token>", data.Token, 1)

		_, err := http_util.PatchJson(url, response)
		if err != nil {
			log.Println(err)
		}
	}()

	if userVoiceState == nil || userVoiceState.GuildId != data.GuildId || *userVoiceState.ChannelId == "" {
		response = util.BuildPlayerResponseData(
			"Can't play a song :(",
			fmt.Sprintf("<@%s> not in voice channel", data.Member.User.Id),
			"Error",
			embed_color.Error,
		)
		return
	}

	song, err := yt_util.SearchYt(options["song"].Value.(string))
	if err != nil {
		log.Println(err)
		response = util.BuildPlayerResponseData(
			"Can't play a song :(",
			"Song not found",
			"Error",
			embed_color.Error,
		)
		return
	}
	song.RequesterId = data.Member.User.Id

	err = voice.StartClient(data.GuildId, *userVoiceState.ChannelId)
	if err != nil {
		if err.Error() == "already in a different voice channel" {
			response = util.BuildPlayerResponseData(
				"Can't play a song :(",
				fmt.Sprintf("<@%s> not in the same voice channel as bot", data.Member.User.Id),
				"Error",
				embed_color.Error,
			)
			return
		} else {
			log.Println(err)
			response = util.BuildPlayerResponseData(
				"Can't play a song :(",
				"Unknown Error",
				"Error",
				embed_color.Error,
			)
			err := voice.StopClient(data.GuildId)
			if err != nil {
				log.Println(err)
			}
			return
		}
	}
	log.Println("Starting client", song)
	queueSize := voice.AppendSongToSongQueue(data.GuildId, song)
	if queueSize == 1 {
		response = util.BuildPlayerResponseData(
			"Play a song",
			fmt.Sprintf("Playing %s\nrequested by <@%s>", song.Title, data.Member.User.Id),
			"Playing",
			embed_color.Playing,
		)
	} else {
		response = util.BuildPlayerResponseData(
			"Play a song",
			fmt.Sprintf("Added %s\nrequested by <@%s> to queue", song.Title, data.Member.User.Id),
			fmt.Sprintf("#%d in queue", queueSize),
			embed_color.Playing,
		)
	}
}
