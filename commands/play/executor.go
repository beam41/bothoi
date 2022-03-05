package play

import (
	"bothoi/config"
	"bothoi/models"
	"bothoi/references/embed_color"
	"bothoi/states"
	"bothoi/util"
	"bothoi/util/http_util"
	"fmt"
	"log"
	"strings"
)

func Execute(data *models.Interaction) {
	options := util.MapInteractionOption(data.Data.Options)
	userVoiceState := states.VoiceState[data.Member.User.ID]

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
	var songQ = states.SongQueue[data.GuildID]
	if songQ != nil {
		if songQ.CurrVoiceChannelID != userVoiceState.ChannelID {
			response = util.BuildPlayerResponse(
				"Can't play a song :(",
				fmt.Sprintf("<@%s> not in the same voice channel as bot", data.Member.User.Username),
				"Error",
				embed_color.Error,
			)
			return
		}
		_ = append(songQ.Songs, models.SongItem{
			YtID:        options["song"].Value.(string),
			Title:       options["song"].Value.(string),
			Duration:    0,
			RequesterID: data.Member.User.ID,
		})
	} else {
		newSongQ := &models.SongQueue{
			Songs: []models.SongItem{
				{
					YtID:        options["song"].Value.(string),
					Title:       options["song"].Value.(string),
					Duration:    0,
					RequesterID: data.Member.User.ID,
				},
			},
			CurrVoiceChannelID: userVoiceState.ChannelID,
			GuildID:            data.GuildID,
		}
		states.SongQueue[data.GuildID] = newSongQ
	}
	log.Println(states.SongQueue[data.GuildID])
	response = util.BuildPlayerResponse(
		"Play a song",
		fmt.Sprintf("Playing %s\nrequested by <@%s>", options["song"].Value.(string), data.Member.User.ID),
		"Playing",
		embed_color.Playing,
	)
}
