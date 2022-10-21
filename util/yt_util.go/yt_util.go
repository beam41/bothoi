package yt_util

import (
	"bothoi/models"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func DownloadYt(ytID string) ([]byte, error) {
	cmd := exec.Command("youtube-dl", "-g", "-f", "bestaudio", "--", ytID)

	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	urlStr := strings.TrimSpace(string(stdout))

	// check content length
	headRes, err := http.Head(urlStr)
	if err != nil || headRes.StatusCode != http.StatusOK {
		return nil, err
	}
	contentLengthStr := headRes.Header.Get("Content-Length")
	contentLength, _ := strconv.ParseInt(contentLengthStr, 10, 64)
	if contentLength == 0 {
		return nil, errors.New("Content-Length is 0")
	}

	bytesArr := make([]byte, contentLength)

	const chunkSize = 100000
	var wg sync.WaitGroup

	// concurrently download chunks
	for pos := int64(0); pos < contentLength; pos += chunkSize {
		wg.Add(1)
		go func(p int64) {
			defer wg.Done()

			req, err := http.NewRequest("GET", urlStr, nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			req.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", p, p+chunkSize-1))
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println(err)
			}
			defer func(body io.ReadCloser) {
				err := body.Close()
				if err != nil {
					log.Println(err)
				}
			}(resp.Body)
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			copy(bytesArr[p:], body)
		}(pos)
	}
	wg.Wait()
	return bytesArr, nil
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

func SearchYt(searchStr string) (models.SongItem, error) {
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
	err := cmd.Run()
	if err != nil {
		return models.SongItem{}, err
	}
	output := strings.TrimSpace(stdout.String())
	if output == "" {
		return models.SongItem{}, errors.New("no results found")
	}
	results := strings.Split(output, "\n")
	if len(results) < 3 {
		return models.SongItem{}, errors.New("no results found")
	}
	return models.SongItem{
		Title:       results[0],
		YtID:        results[1],
		Duration:    results[2],
		RequesterID: "",
	}, nil
}
