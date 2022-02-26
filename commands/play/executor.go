package play

import (
	"bothoi/models"
	"bothoi/util"
	"bothoi/util/http_util"
	"log"
	"os"
	"strings"
)

func Execute(data *models.Interaction) {
	options := util.MapInteractionOption(data.Data.Options)

	// response
	response := models.InteractionResponse{
		Type: 4,
		Data: models.InteractionResponseData{
			Embeds: []models.Embed{
				{
					Title:       "Play a song",
					Description: "Playing " + options["song"].Value.(string) + "\n" + "requested by <@" + data.Member.User.ID + ">",
					Footer: &models.EmbedFooter{
						Text: "Playing",
					},
					Author: &models.EmbedAuthor{
						Name:    "Playing",
						IconUrl: "https://avatars.githubusercontent.com/u/1791353?v=4",
					},
				},
			},
			AllowedMentions: &models.AllowedMention{
				Parse: []string{"users"},
			},
		},
	}

	url := os.Getenv("INTERACTION_RESPONSE_ENDPOINT")
	url = strings.Replace(url, "<interaction_id>", data.ID, 1)
	url = strings.Replace(url, "<interaction_token>", data.Token, 1)

	_, err := http_util.PostJson(url, response)
	if err != nil {
		log.Println(err)
	}
}
