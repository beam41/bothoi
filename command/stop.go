package command

import (
	"bothoi/bh_context"
	"bothoi/config"
	"bothoi/models/discord_models"
	"bothoi/references/embed_color"
	"bothoi/repo"
	"bothoi/util"
	"strconv"
)

const commandStop = "stop"

func executeStop(data *discord_models.Interaction) {
	var response discord_models.InteractionResponse
	defer responseNoLoading(data.ID, data.Token, response)

	userVoiceChannel := repo.GetChannelIDByUserIDAndGuildID(data.Member.User.ID, data.GuildID)
	clientVoiceChannel := repo.GetChannelIDByUserIDAndGuildID(config.BotID, data.GuildID)
	if pass, res := checkNotSameChannelError(util.BuildPlayerResponse, userVoiceChannel, clientVoiceChannel, "Pause", data.Member.User.ID); pass {
		response = res
		return
	}

	err := bh_context.GetVoiceClientManager().StopClient(data.GuildID)
	if err != nil {
		response = util.BuildPlayerResponse(
			"Stopped",
			"Cannot be stopped",
			"Error",
			embed_color.Error,
		)
		return
	}

	response = util.BuildPlayerResponse(
		"Stopped",
		"Stopped by the request of <@"+strconv.FormatUint(uint64(data.Member.User.ID), 10)+">",
		"Stopped",
		embed_color.EmbedColor(0xff1744),
	)
}
