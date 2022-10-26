package command

import (
	"bothoi/bh_context"
	"bothoi/config"
	"bothoi/models/discord_models"
	"bothoi/references/embed_color"
	"bothoi/repo"
	"bothoi/util"
	"bothoi/util/http_util"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const commandPause0 = "pause"
const commandPause1 = "resume"

func executePause(data *discord_models.Interaction) {
	var response discord_models.InteractionResponse
	// do response to interaction
	defer func() {
		url := config.InteractionResponseEndpoint
		url = strings.Replace(url, "<interaction_id>", strconv.FormatUint(uint64(data.ID), 10), 1)
		url = strings.Replace(url, "<interaction_token>", data.Token, 1)

		_, err := http_util.PostJson(url, response)
		if err != nil {
			log.Println(err)
		}
	}()
	userVoiceChannel := repo.GetChannelIDByUserIDAndGuildID(data.Member.User.ID, data.GuildID)
	clientVoiceChannel := repo.GetChannelIDByUserIDAndGuildID(config.BotID, data.GuildID)
	if userVoiceChannel == nil || (clientVoiceChannel != nil && *userVoiceChannel != *clientVoiceChannel) {
		response = util.BuildPlayerResponse(
			"Pause error",
			fmt.Sprintf("<@%d> not in same voice channel as Bothoi", data.Member.User.ID),
			"error",
			embed_color.Error,
		)
		return
	}
	pausing, err := bh_context.GetVoiceClientManager().PauseClient(data.GuildID)
	if err != nil {
		response = util.BuildPlayerResponse(
			"Pause error",
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
			embed_color.Default,
		)
	} else {
		response = util.BuildPlayerResponse(
			"Resumed",
			"Resumed by the request of <@"+strconv.FormatUint(uint64(data.Member.User.ID), 10)+">",
			"",
			embed_color.Default,
		)
	}

}
