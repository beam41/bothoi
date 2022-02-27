package play

import (
	"bothoi/models"
	"bothoi/references/embed_color"
	"bothoi/states"
	"bothoi/util"
	"bothoi/util/http_util"
	"fmt"
	"log"
	"os"
	"strings"
)

func Execute(data *models.Interaction) {
	options := util.MapInteractionOption(data.Data.Options)
	userVoiceState := states.VoiceState[data.Member.User.ID]

	var response models.InteractionResponse
	if userVoiceState == nil || userVoiceState.GuildID != data.GuildID || userVoiceState.ChannelID == "" {
		response = util.BuildBothoiPlayerResponse(
			"Can't play a song :(",
			fmt.Sprintf("<@%s> not in voice channel", data.Member.User.Username),
			"Error",
			embed_color.Error,
		)
	} else {
		response = util.BuildBothoiPlayerResponse(
			"Play a song",
			fmt.Sprintf("Playing %s\nrequested by <@%s>", options["song"].Value.(string), data.Member.User.ID),
			"Playing",
			embed_color.Playing,
		)
	}

	url := os.Getenv("INTERACTION_RESPONSE_ENDPOINT")
	url = strings.Replace(url, "<interaction_id>", data.ID, 1)
	url = strings.Replace(url, "<interaction_token>", data.Token, 1)

	_, err := http_util.PostJson(url, response)
	if err != nil {
		log.Println(err)
	}
}
