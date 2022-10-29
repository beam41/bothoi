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

func GetYoutubeDownloadUrl(ytID string) (string, error) {
	cmd := exec.Command("youtube-dl", "-g", "-f", "bestaudio", "--", ytID)

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

	if strings.Contains(u.Hostname(), "youtube.com") {
		if u.Path == "/watch" {
			v := u.Query().Get("v")
			// yt vid id is 11 characters long
			match, err := regexp.MatchString("^[a-zA-Z0-9_-]{11}$", v)
			if err != nil {
				return false
			}
			return match
		} else if u.Path == "/playlist" {
			return true
		}
	} else if strings.Contains(u.Hostname(), "youtu.be") {
		match, err := regexp.MatchString("^/[a-zA-Z0-9_-]{11}$", u.Path)
		if err != nil {
			return false
		}
		return match
	}

	return false
}

type SearchResult struct {
	Title    string
	YtID     string
	Duration uint32
}

func SearchYt(searchStr string) ([]SearchResult, error) {
	var cmd *exec.Cmd
	if !isYtVidUrl(searchStr) {
		// add \" to escape quotes in cmd
		slashedStr := strings.Replace(searchStr, "\"", "\\\"", -1)
		// bestaudio may not be available, so I have to include it
		cmd = exec.Command("youtube-dl", "--get-title", "--get-id", "--get-duration", "-f", "bestaudio", "--ignore-errors", fmt.Sprintf("ytsearch:\"%s\"", slashedStr))
	} else {
		cmd = exec.Command("youtube-dl", "--get-title", "--get-id", "--get-duration", "-f", "bestaudio", "--ignore-errors", searchStr)
	}
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	// even when we ignore-errors, youtube-dl still exit with non-zero if there are errors
	err := cmd.Run()
	output := strings.TrimSpace(stdout.String())
	if output == "" {
		return []SearchResult{}, err
	}
	strResults := strings.Split(output, "\n")
	if len(strResults)%3 != 0 {
		return []SearchResult{}, err
	}
	results := make([]SearchResult, len(strResults)/3)
	for i, j := 0, 0; i < len(results); i, j = i+1, j+3 {
		results[i] = SearchResult{
			Title:    strResults[j],
			YtID:     strResults[j+1],
			Duration: util.ConvertVidLengthToSeconds(strResults[j+2]),
		}
	}

	return results, err
}
