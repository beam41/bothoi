package yt_util

import (
	"bothoi/models"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"strings"
)

func DownloadYt(ytID string) ([]byte, error) {
	cmd := exec.Command("youtube-dl", "-g", "-f", "bestaudio", ytID)

	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	urlStr := strings.TrimSpace(string(stdout))

	// i think i need to download to a temp file because this one is somehow super slow
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.Bytes(), nil
}

func isYtVidUrl(testUrl string) bool {
	u, err := url.Parse(testUrl)
	if err == nil && u.Scheme != "" && u.Host != "" {
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

func SearchYt(searchStr string) (models.SongItem, error) {
	var cmd *exec.Cmd
	if !isYtVidUrl(searchStr) {
		// add \" to escape quotes in cmd
		slashedStr := strings.Replace(searchStr, "\"", "\\\"", -1)
		cmd = exec.Command("youtube-dl", "--get-title", "--get-id", "--get-duration", "-f", "bestaudio", fmt.Sprintf("ytsearch:\"%s\"", slashedStr))
	} else {
		cmd = exec.Command("youtube-dl", "--get-title", "--get-id", "--get-duration", "-f", "bestaudio", searchStr)
	}
	stdout, err := cmd.Output()
	if err != nil {
		return models.SongItem{}, err
	}
	output := strings.TrimSpace(string(stdout))
	if output == "" {
		return models.SongItem{}, errors.New("No results found")
	}
	results := strings.Split(output, "\n")
	if len(results) < 3 {
		return models.SongItem{}, errors.New("No results found")
	}
	return models.SongItem{
		Title:       results[0],
		YtID:        results[1],
		Duration:    results[2],
		RequesterID: "",
	}, nil
}
