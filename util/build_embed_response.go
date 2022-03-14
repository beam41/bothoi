package util

import (
	"bothoi/models"
	"bothoi/references/embed_color"
)

func BuildPlayerResponse(title string, desc string, footerText string, color embed_color.EmbedColor) models.InteractionResponse {
	return models.InteractionResponse{
		Type: 4,
		Data: BuildPlayerResponseData(title, desc, footerText, color),
	}
}

func BuildPlayerResponseData(title string, desc string, footerText string, color embed_color.EmbedColor) models.InteractionResponseData {
	return models.InteractionResponseData{
		Embeds: []models.Embed{
			{
				Title:       title,
				Description: desc,
				Footer: &models.EmbedFooter{
					Text: footerText,
				},
				Author: &models.EmbedAuthor{
					Name:    "Bothoi Player",
					IconUrl: "https://avatars.githubusercontent.com/u/1791353?v=4",
				},
				Color: int32(color),
			},
		},
	}
}
