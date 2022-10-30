package yt_util

import (
	"bytes"
	"encoding/json"
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

var ytIDRegex = regexp.MustCompile("^[a-zA-Z0-9_-]{11}$")

func isYtVidUrl(testUrl string) (vid bool, playlist bool) {
	u, err := url.Parse(testUrl)
	if err != nil {
		return false, false
	}

	if strings.Contains(u.Hostname(), "youtube.com") {
		if u.Path == "/watch" {
			v := u.Query().Get("v")
			// yt vid id is 11 characters long
			match := ytIDRegex.MatchString(v)
			return match, false
		} else if u.Path == "/playlist" {
			return true, true
		}
	} else if strings.Contains(u.Hostname(), "youtu.be") {
		match := ytIDRegex.MatchString(u.Path)
		return match, false
	}

	return false, false
}

func SearchYt(searchStr string) ([]Video, *PlaylistInfo, error) {
	var cmd *exec.Cmd
	isVid, isPlaylist := isYtVidUrl(searchStr)
	if isPlaylist {
		cmd = exec.Command("youtube-dl", "--ignore-errors", "--flat-playlist", "--dump-single-json", searchStr)
	} else if isVid {
		cmd = exec.Command("youtube-dl", "--ignore-errors", "-f", "bestaudio", "--dump-json", searchStr)
	} else {
		// add \" to escape quotes in cmd
		slashedStr := strings.Replace(searchStr, "\"", "\\\"", -1)
		// bestaudio may not be available, so I have to include it
		cmd = exec.Command("youtube-dl", "--ignore-errors", "-f", "bestaudio", "--dump-json", fmt.Sprintf("ytsearch:\"%s\"", slashedStr))
	}
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	// even when we ignore-errors, youtube-dl still exit with non-zero if there are errors
	_ = cmd.Run()
	if isPlaylist {
		var result Playlist
		err := json.Unmarshal(stdout.Bytes(), &result)
		if err != nil {
			return []Video{}, nil, err
		}
		return result.Entries, &PlaylistInfo{result.Title, result.WebpageURL}, nil
	}
	var result Video
	err := json.Unmarshal(stdout.Bytes(), &result)
	if err != nil {
		return []Video{}, nil, err
	}
	return []Video{result}, nil, nil
}
