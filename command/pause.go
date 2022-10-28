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

const commandPause0 = "pause"
const commandPause1 = "resume"

func executePause(data *discord_models.Interaction) {
	var response discord_models.InteractionResponse
	defer responseNoLoading(data.ID, data.Token, response)

	userVoiceChannel := repo.GetChannelIDByUserIDAndGuildID(data.Member.User.ID, data.GuildID)
	clientVoiceChannel := repo.GetChannelIDByUserIDAndGuildID(config.BotID, data.GuildID)
	if pass, res := checkNotSameChannelError(util.BuildPlayerResponse, userVoiceChannel, clientVoiceChannel, "Pause", data.Member.User.ID); pass {
		response = res
		return
	}

	pausing, err := bh_context.GetVoiceClientManager().PauseClient(data.GuildID)
	if err != nil {
		response = util.BuildPlayerResponse(
			"Pause Error",
			"Cannot be paused/resumed",
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
