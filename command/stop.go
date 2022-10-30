package command

import (
	"bothoi/config"
	"bothoi/models/discord_models"
	"bothoi/references/embed_color"
	"bothoi/repo"
	"bothoi/util"
	"log"
	"strconv"
)

func (client *Manager) executeStop(data *discord_models.Interaction) {
	var response discord_models.InteractionResponse
	defer func() { postResponse(data.ID, data.Token, response) }()

	userVoiceChannel := repo.GetChannelIDByUserIDAndGuildID(data.Member.User.ID, data.GuildID)
	clientVoiceChannel := repo.GetChannelIDByUserIDAndGuildID(config.BotID, data.GuildID)
	if pass, res := checkNotSameChannelError(util.BuildPlayerResponse, userVoiceChannel, clientVoiceChannel, "Pause", data.Member.User.ID); pass {
		response = res
		return
	}

	err := client.voiceClientManager.ClientStop(data.GuildID)
	if err != nil {
		log.Println(err)
		response = util.BuildPlayerResponse(
			"Stopped",
			"Cannot be stopped",
			"",
			embed_color.Error,
		)
		return
	}

	response = util.BuildPlayerResponse(
		"Stopped",
		"Stopped by the request of <@"+strconv.FormatUint(uint64(data.Member.User.ID), 10)+">",
		"",
		embed_color.ErrorLow,
	)
}
