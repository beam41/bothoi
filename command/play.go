package command

import (
	"bothoi/config"
	"bothoi/models/discord_models"
	"bothoi/references/embed_color"
	"bothoi/repo"
	"bothoi/util"
	"bothoi/util/yt_util.go"
	"fmt"
	"log"
)

func (client *Manager) executePlay(data *discord_models.Interaction) {
	postLoading(data.ID, data.Token, "Play")

	options := util.SliceToMap(data.Data.Options, func(i int, item discord_models.InteractionOption) string { return item.Name })

	var response discord_models.InteractionCallbackData
	defer func() { patchResponseAfterLoading(data.Token, response) }()

	userVoiceChannel := repo.GetChannelIDByUserIDAndGuildID(data.Member.User.ID, data.GuildID)
	clientVoiceChannel := repo.GetChannelIDByUserIDAndGuildID(data.Member.User.ID, config.BotID)
	if pass, res := checkNotSameChannelError(util.BuildPlayerResponseData, userVoiceChannel, clientVoiceChannel, "Play", data.Member.User.ID); pass {
		response = res
		return
	}

	result, playlistInfo, _ := yt_util.SearchYt(options["song"].Value.(string))
	if len(result) == 0 {
		response = util.BuildPlayerResponseData(
			"Play Error",
			"Song not found",
			"Error",
			embed_color.Error,
		)
		return
	}

	if len(result) == 1 {
		seek := uint32(0)
		if s, ok := options["seek"]; ok {
			seek = util.ConvertVidLengthToSeconds(s.Value.(string))
		}
		queueSize := repo.GetQueueSize(data.GuildID)
		err := repo.AddSongToQueue(data.GuildID, data.Member.User.ID, data.ChannelID, result[0].ID, result[0].Title, uint32(result[0].Duration), seek, queueSize != 0)
		if err != nil {
			log.Println(err)
			response = util.BuildPlayerResponseData(
				"Play Error",
				"Can't add a song",
				"",
				embed_color.Error,
			)
			return
		}
	} else {
		err := repo.AddSongToQueueMultiple(data.GuildID, data.Member.User.ID, data.ChannelID, result)
		if err != nil {
			log.Println(err)
			response = util.BuildPlayerResponseData(
				"Play Error",
				"Can't add songs",
				"",
				embed_color.Error,
			)
			return
		}
	}

	err := client.voiceClientManager.ClientStart(data.GuildID, *userVoiceChannel)
	if err != nil {
		log.Println(err)
		response = util.BuildPlayerResponseData(
			"Play Error",
			"Can't start player",
			"",
			embed_color.Error,
		)
		return
	}

	queueSize := repo.GetQueueSize(data.GuildID)
	if len(result) > 1 {
		response = util.BuildPlayerResponseData(
			"Play a song",
			fmt.Sprintf(
				"Added %d song%s from [%s](%s)\nRequested by <@%d>",
				len(result),
				util.Ternary(len(result) > 1, "s", ""),
				playlistInfo.Title,
				playlistInfo.WebpageURL,
				data.Member.User.ID,
			),
			fmt.Sprintf("%d song%s in queue", queueSize, util.Ternary(queueSize > 1, "s", "")),
			embed_color.SuccessScheduled,
		)
	} else if queueSize == 1 {
		response = util.BuildPlayerResponseData(
			"Play a song",
			fmt.Sprintf(
				"Playing [%s](https://youtu.be/%s) | `%s`\nRequested by <@%d>",
				result[0].Title,
				result[0].ID,
				util.ConvertSecondsToVidLength(uint32(result[0].Duration)),
				data.Member.User.ID,
			),
			"Playing",
			embed_color.SuccessContinue,
		)
	} else {
		response = util.BuildPlayerResponseData(
			"Play a song",
			fmt.Sprintf(
				"Added [%s](https://youtu.be/%s) | `%s`\nRequested by <@%d>",
				result[0].Title,
				result[0].ID,
				util.ConvertSecondsToVidLength(uint32(result[0].Duration)),
				data.Member.User.ID,
			),
			fmt.Sprintf("#%d in queue", queueSize),
			embed_color.SuccessScheduled,
		)
	}
}
