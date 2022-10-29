package util

import (
	"bothoi/models/discord_models"
	"bothoi/references/embed_color"
)

func BuildPlayerResponse(title string, desc string, footerText string, color embed_color.EmbedColor) discord_models.InteractionResponse {
	return discord_models.InteractionResponse{
		Type: 4,
		Data: BuildPlayerResponseData(title, desc, footerText, color),
	}
}

func BuildPlayerResponseData(title string, desc string, footerText string, color embed_color.EmbedColor) discord_models.InteractionCallbackData {
	return discord_models.InteractionCallbackData{
		Embeds: []discord_models.Embed{
			{
				Title:       title,
				Description: desc,
				Footer: &discord_models.EmbedFooter{
					Text: footerText,
				},
				Author: &discord_models.EmbedAuthor{
					Name:    "Bothoi Player",
					IconUrl: "https://avatars.githubusercontent.com/u/1791353?v=4",
				},
				Color: uint32(color),
			},
		},
	}
}
