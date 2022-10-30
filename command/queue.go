package command

import (
	"bothoi/models/discord_models"
	"bothoi/references/embed_color"
	"bothoi/repo"
	"bothoi/util"
	"fmt"
	"strconv"
	"strings"
)

func (client *Manager) executeQueue(data *discord_models.Interaction) {
	postLoading(data.ID, data.Token, "Queue")

	options := util.SliceToMap(data.Data.Options, func(i int, item discord_models.InteractionOption) string { return item.Name })

	var response discord_models.InteractionCallbackData
	defer func() { patchResponseAfterLoading(data.Token, response) }()

	page := 0
	if opPage, ok := options["page"]; ok {
		page = int(opPage.Value.(float64)) - 1
	}

	var queueSize = repo.GetQueueSize(data.GuildID)
	if queueSize == 0 {
		response = util.BuildPlayerResponseData(
			"Queue",
			"No song in queue",
			"Start playing some songs now!",
			embed_color.ErrorLow,
		)
		return
	}

	var maxPage = int(queueSize-1) / 10
	if maxPage < page {
		page = maxPage
	}
	offset := page * 10

	var songQ = repo.GetSongQueue(data.GuildID, offset, 10)
	var res strings.Builder
	for i, song := range songQ {
		res.WriteString(
			fmt.Sprintf(
				"%s [%s](https://youtu.be/%s) | `%s` | <@%d>\n",
				util.Ternary(song.Playing, "**Playing:**", strconv.Itoa(offset+i+1)+"."),
				song.Title,
				song.YtID,
				util.ConvertSecondsToVidLength(song.Duration),
				song.RequesterID,
			),
		)
	}

	response = util.BuildPlayerResponseData(
		"Queue",
		res.String(),
		fmt.Sprintf("Page %d/%d | %d song%s in queue", page+1, maxPage+1, queueSize, util.Ternary(queueSize == 1, "", "s")),
		embed_color.Info,
	)
}
