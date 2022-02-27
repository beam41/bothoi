package util

import "bothoi/models"

func BuildBothoiPlayerResponse(title string, desc string, footerText string, color int32) models.InteractionResponse {
	return models.InteractionResponse{
		Type: 4,
		Data: models.InteractionResponseData{
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
					Color: color,
				},
			},
		},
	}
}
