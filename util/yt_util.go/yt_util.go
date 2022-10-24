package yt_util

import (
	"bothoi/util"
	"bytes"
	"fmt"
	"net/url"
	"os/exec"
	"regexp"
	"strings"
)

func GetYoutubeDownloadUrl(ytId string) (string, error) {
	cmd := exec.Command("youtube-dl", "-g", "-f", "bestaudio", "--", ytId)

	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	urlStr := strings.TrimSpace(string(stdout))
	return urlStr, nil
}

func isYtVidUrl(testUrl string) bool {
	u, err := url.Parse(testUrl)
	if err != nil {
		return false
	}

	if u.Host == "www.youtube.com" && u.Path == "/watch" {
		v := u.Query().Get("v")
		// yt vid id is 11 characters long
		match, err := regexp.MatchString("^[a-zA-Z0-9_-]{11}$", v)
		if err != nil {
			return false
		}
		return match
	} else if u.Host == "youtu.be" {
		match, err := regexp.MatchString("^/[a-zA-Z0-9_-]{11}$", u.Path)
		if err != nil {
			return false
		}
		return match
	}
	return false
}

func SearchYt(searchStr string) (title string, ytId string, duration uint32, noResult bool, err error) {
	var cmd *exec.Cmd
	if !isYtVidUrl(searchStr) {
		// add \" to escape quotes in cmd
		slashedStr := strings.Replace(searchStr, "\"", "\\\"", -1)
		cmd = exec.Command("youtube-dl", "--get-title", "--get-id", "--get-duration", "-f", "bestaudio", fmt.Sprintf("ytsearch:\"%s\"", slashedStr))
	} else {
		cmd = exec.Command("youtube-dl", "--get-title", "--get-id", "--get-duration", "-f", "bestaudio", searchStr)
	}
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err = cmd.Run()
	if err != nil {
		noResult = true
		return
	}
	output := strings.TrimSpace(stdout.String())
	if output == "" {
		noResult = true
		return
	}
	results := strings.Split(output, "\n")
	if len(results) < 3 {
		noResult = true
		return
	}
	title = results[0]
	ytId = results[1]
	duration = util.ConvertVidLengthToSeconds(results[2])
	return
}
