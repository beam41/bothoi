package command

import (
	"bothoi/command/command_interface"
	"bothoi/config"
	"bothoi/models/discord_models"
	"bothoi/models/types"
	"bothoi/references/embed_color"
	"bothoi/util"
	"bothoi/util/http_util"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type commandManager struct {
	executorList map[string]func(*discord_models.Interaction)
}

func NewCommandManager() command_interface.CommandManagerInterface {
	return &commandManager{
		executorList: map[string]func(*discord_models.Interaction){
			commandPlay:   executePlay,
			commandQueue:  executeQueue,
			commandPause0: executePause,
			commandPause1: executePause,
			commandStop:   executeStop,
			commandSkip:   executeSkip,
		},
	}
}

func (cm *commandManager) MapInteractionExecute(data *discord_models.Interaction) {
	if interaction, ok := cm.executorList[data.Data.Name]; ok {
		interaction(data)
	}
}

func checkNotSameChannelError[InteractionResponse discord_models.InteractionResponse | discord_models.InteractionCallbackData](
	builder func(title string, desc string, footerText string, color embed_color.EmbedColor) InteractionResponse,
	userVoiceChannel *types.Snowflake,
	clientVoiceChannel *types.Snowflake,
	cmd string,
	userID types.Snowflake,
) (bool, InteractionResponse) {
	desc := ""
	if userVoiceChannel == nil {
		desc = "<@%d> not in a voice channel"
	} else if clientVoiceChannel != nil && *userVoiceChannel != *clientVoiceChannel {
		desc = "<@%d> not in same voice channel as Bothoi"
	}

	if desc != "" {
		return true, builder(
			cmd+" Error",
			fmt.Sprintf(desc, userID),
			"Error",
			embed_color.Error,
		)
	}
	return false, nil
}

func responseNoLoading(id types.Snowflake, token string, response discord_models.InteractionResponse) {
	url := config.InteractionResponseEndpoint
	url = strings.Replace(url, "<interaction_id>", strconv.FormatUint(uint64(id), 10), 1)
	url = strings.Replace(url, "<interaction_token>", token, 1)

	_, err := http_util.PostJson(url, response)
	if err != nil {
		log.Println(err)
	}
}

func postLoading(id types.Snowflake, token string, cmd string) {
	// post waiting prevent response timeout
	url := config.InteractionResponseEndpoint
	url = strings.Replace(url, "<interaction_id>", strconv.FormatUint(uint64(id), 10), 1)
	url = strings.Replace(url, "<interaction_token>", token, 1)

	_, err := http_util.PostJson(url, util.BuildPlayerResponse(
		cmd,
		"Loading...",
		"Please Wait",
		embed_color.Playing,
	))
	if err != nil {
		log.Println(err)
	}
}

func responseAfterLoading(token string, response discord_models.InteractionCallbackData) {
	url := config.InteractionResponseEditEndpoint
	url = strings.Replace(url, "<application_id>", strconv.FormatUint(uint64(config.BotID), 10), 1)
	url = strings.Replace(url, "<interaction_token>", token, 1)

	_, err := http_util.PatchJson(url, response)
	if err != nil {
		log.Println(err)
	}
}
