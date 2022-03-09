package play

import (
	"bothoi/config"
	"bothoi/models"
	"bothoi/references/embed_color"
	"bothoi/states"
	"bothoi/util"
	"bothoi/util/http_util"
	"bothoi/voice"
	"fmt"
	"log"
	"strings"
)

func Execute(data *models.Interaction) {
	options := util.MapInteractionOption(data.Data.Options)
	userVoiceState := states.GetVoiceState(data.Member.User.ID)

	var response models.InteractionResponse
	// do response to interaction
	defer func() {
		url := config.INTERACTION_RESPONSE_ENDPOINT
		url = strings.Replace(url, "<interaction_id>", data.ID, 1)
		url = strings.Replace(url, "<interaction_token>", data.Token, 1)

		_, err := http_util.PostJson(url, response)
		if err != nil {
			log.Println(err)
		}
	}()

	if userVoiceState == nil || userVoiceState.GuildID != data.GuildID || userVoiceState.ChannelID == "" {
		response = util.BuildPlayerResponse(
			"Can't play a song :(",
			fmt.Sprintf("<@%s> not in voice channel", data.Member.User.Username),
			"Error",
			embed_color.Error,
		)
		return
	}
	err := voice.StartClient(data.GuildID, userVoiceState.ChannelID)
	if err != nil {
		if err.Error() == "Already in a different voice channel" {
			response = util.BuildPlayerResponse(
				"Can't play a song :(",
				fmt.Sprintf("<@%s> not in the same voice channel as bot", data.Member.User.Username),
				"Error",
				embed_color.Error,
			)
			return
		} else {
			log.Println(err)
			response = util.BuildPlayerResponse(
				"Can't play a song :(",
				"Unknown Error",
				"Error",
				embed_color.Error,
			)
			voice.StopClient(data.GuildID)
			return
		}
	}
	queueSize := voice.AppendSongToSongQueue(data.GuildID, models.SongItem{
		YtID:        options["song"].Value.(string),
		Title:       options["song"].Value.(string),
		Duration:    0,
		RequesterID: data.Member.User.ID,
	})
	if queueSize == 1 {
		response = util.BuildPlayerResponse(
			"Play a song",
			fmt.Sprintf("Playing %s\nrequested by <@%s>", options["song"].Value.(string), data.Member.User.ID),
			"Playing",
			embed_color.Playing,
		)
	} else {
		response = util.BuildPlayerResponse(
			"Play a song",
			fmt.Sprintf("Added %s\nrequested by <@%s> to queue", options["song"].Value.(string), data.Member.User.ID),
			fmt.Sprintf("#%d in queue", queueSize),
			embed_color.Playing,
		)
	}
}
