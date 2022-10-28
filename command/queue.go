package command

import (
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

const commandQueue = "queue"

func executeQueue(data *discord_models.Interaction) {
	options := util.SliceToMap(data.Data.Options, func(i int, item discord_models.InteractionOption) string { return item.Name })

	// post waiting prevent response timeout
	url := config.InteractionResponseEndpoint
	url = strings.Replace(url, "<interaction_id>", strconv.FormatUint(uint64(data.ID), 10), 1)
	url = strings.Replace(url, "<interaction_token>", data.Token, 1)

	_, err := http_util.PostJson(url, util.BuildPlayerResponse(
		"Queue",
		"Loading...",
		"please wait",
		embed_color.Default,
	))
	if err != nil {
		log.Println(err)
	}

	var response discord_models.InteractionCallbackData
	// do response to interaction
	defer func() {
		url := config.InteractionResponseEditEndpoint
		url = strings.Replace(url, "<application_id>", strconv.FormatUint(uint64(config.BotID), 10), 1)
		url = strings.Replace(url, "<interaction_token>", data.Token, 1)

		_, err := http_util.PatchJson(url, response)
		if err != nil {
			log.Println(err)
		}
	}()

	page := 0
	if opPage, ok := options["page"]; ok {
		page = int(opPage.Value.(float64)) - 1
	}

	var queueSize = repo.GetQueueSize(data.GuildID)
	if queueSize == 0 {
		response = util.BuildPlayerResponseData(
			"No songs in queue",
			"Start playing a song now!",
			"Queue",
			embed_color.Error,
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
	res.WriteString("Song in queue (Requested by)\n")
	for i, song := range songQ {
		res.WriteString(
			fmt.Sprintf(
				"%s %s** | %s | <@%d>**\n",
				util.Ternary(song.Playing, "Playing:", strconv.Itoa(offset+i+1)+"."),
				song.Title,
				util.ConvertSecondsToVidLength(song.Duration),
				song.RequesterID,
			),
		)
	}

	response = util.BuildPlayerResponseData(
		"Queue",
		res.String(),
		fmt.Sprintf("page %d/%d, %d song%s in queue", page+1, maxPage+1, queueSize, util.Ternary(queueSize == 1, "", "s")),
		embed_color.Default,
	)
}
