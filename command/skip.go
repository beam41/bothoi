package command

import (
	"bothoi/config"
	"bothoi/models/discord_models"
	"bothoi/references/embed_color"
	"bothoi/repo"
	"bothoi/util"
	"fmt"
)

func (client *Manager) executeSkip(data *discord_models.Interaction) {
	var response discord_models.InteractionResponse
	defer func() { postResponse(data.ID, data.Token, response) }()

	userVoiceChannel := repo.GetChannelIDByUserIDAndGuildID(data.Member.User.ID, data.GuildID)
	clientVoiceChannel := repo.GetChannelIDByUserIDAndGuildID(config.BotID, data.GuildID)
	if pass, res := checkNotSameChannelError(util.BuildPlayerResponse, userVoiceChannel, clientVoiceChannel, "Pause", data.Member.User.ID); pass {
		response = res
		return
	}

	success := client.voiceClientManager.ClientSkipSong(data.GuildID)
	if !success {
		response = util.BuildPlayerResponse(
			"Skip Error",
			"Cannot skip",
			"",
			embed_color.Error,
		)
		return
	}

	response = util.BuildPlayerResponse(
		"Skipped",
		fmt.Sprintf("Skipped by the request of <@%d>", uint64(data.Member.User.ID)),
		"",
		embed_color.SuccessInterrupt,
	)
}
