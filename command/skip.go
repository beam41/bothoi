package command

import (
	"bothoi/bh_context"
	"bothoi/config"
	"bothoi/models/discord_models"
	"bothoi/references/embed_color"
	"bothoi/repo"
	"bothoi/util"
	"fmt"
)

const commandSkip = "skip"

func executeSkip(data *discord_models.Interaction) {
	var response discord_models.InteractionResponse
	defer responseNoLoading(data.ID, data.Token, response)

	userVoiceChannel := repo.GetChannelIDByUserIDAndGuildID(data.Member.User.ID, data.GuildID)
	clientVoiceChannel := repo.GetChannelIDByUserIDAndGuildID(config.BotID, data.GuildID)
	if pass, res := checkNotSameChannelError(util.BuildPlayerResponse, userVoiceChannel, clientVoiceChannel, "Pause", data.Member.User.ID); pass {
		response = res
		return
	}

	err := bh_context.GetVoiceClientManager().SkipSong(data.GuildID)
	if err != nil {
		response = util.BuildPlayerResponse(
			"Skip Error",
			"Cannot skip",
			"Error",
			embed_color.Error,
		)
		return
	}

	response = util.BuildPlayerResponse(
		"Skipped",
		fmt.Sprintf("Skipped by the request of <@%d>", uint64(data.Member.User.ID)),
		"Skipped",
		embed_color.SuccessInterrupt,
	)
}
