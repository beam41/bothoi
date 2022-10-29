package command

import (
	"bothoi/config"
	"bothoi/gateway"
	"bothoi/models/discord_models"
	"bothoi/references/embed_color"
	"bothoi/repo"
	"bothoi/util"
	"strconv"
)

const commandPause0 = "pause"
const commandPause1 = "resume"

func executePause(gatewayClient *gateway.Client, data *discord_models.Interaction) {
	var response discord_models.InteractionResponse
	defer func() { postResponse(data.ID, data.Token, response) }()

	userVoiceChannel := repo.GetChannelIDByUserIDAndGuildID(data.Member.User.ID, data.GuildID)
	clientVoiceChannel := repo.GetChannelIDByUserIDAndGuildID(config.BotID, data.GuildID)
	if pass, res := checkNotSameChannelError(util.BuildPlayerResponse, userVoiceChannel, clientVoiceChannel, "Pause", data.Member.User.ID); pass {
		response = res
		return
	}

	found, pausing := gatewayClient.VoiceClientPauseSong(data.GuildID)
	if !found {
		response = util.BuildPlayerResponse(
			"Pause Error",
			"Client not found",
			"error",
			embed_color.Error,
		)
		return
	}

	if pausing {
		response = util.BuildPlayerResponse(
			"Paused",
			"Paused by the request of <@"+strconv.FormatUint(uint64(data.Member.User.ID), 10)+">",
			"/resume to resume",
			embed_color.SuccessInterrupt,
		)
	} else {
		response = util.BuildPlayerResponse(
			"Resumed",
			"Resumed by the request of <@"+strconv.FormatUint(uint64(data.Member.User.ID), 10)+">",
			"/pause to pause",
			embed_color.SuccessContinue,
		)
	}
}
