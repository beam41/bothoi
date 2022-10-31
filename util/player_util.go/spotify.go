package player_util

import (
	"bothoi/config"
	"bothoi/util"
	"bothoi/util/http_util"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type spotifyToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func IsSpotifyUrl(testUrl string) bool {
	u, err := url.Parse(testUrl)
	if err != nil {
		return false
	}
	return strings.Contains(u.Hostname(), "open.spotify.com")
}

func spotifyAuth() (spotifyToken, error) {
	authEncoded := base64.StdEncoding.EncodeToString([]byte(config.SpotifyClientID + ":" + config.SpotifyClientSecret))
	header := map[string]string{
		"Authorization": "Basic " + authEncoded,
		"Content-Type":  "application/x-www-form-urlencoded",
	}
	body := url.Values{
		"grant_type": []string{"client_credentials"},
	}
	post, err := http_util.Post("https://accounts.spotify.com/api/token", strings.NewReader(body.Encode()), header)
	if err != nil {
		return spotifyToken{}, err
	}
	var tk spotifyToken
	err = json.Unmarshal(post, &tk)
	return tk, err
}

type spotifyType byte

const (
	spotifyTrack spotifyType = iota + 1
	spotifyAlbum
	spotifyPlaylist
)

func extractSpotifyUrl(sUrl string) (spotifyType, string) {
	u, err := url.Parse(sUrl)
	if err != nil {
		return 0, ""
	}
	path := strings.Split(u.Path[1:], "/")
	if len(path) < 2 {
		return 0, ""
	}
	switch path[0] {
	case "track":
		return spotifyTrack, path[1]
	case "album":
		return spotifyAlbum, path[1]
	case "playlist":
		return spotifyPlaylist, path[1]
	default:
		return 0, ""
	}
}

func ExtractSpotifyUrl(sUrl string) (searchStr []string, name string, link string, err error) {
	urlType, id := extractSpotifyUrl(sUrl)
	if urlType < spotifyTrack || urlType > spotifyPlaylist {
		err = errors.New("unknown url")
		return
	}
	var typeRoute string
	switch urlType {
	case spotifyTrack:
		typeRoute = "tracks"
	case spotifyAlbum:
		typeRoute = "albums"
	case spotifyPlaylist:
		typeRoute = "playlists"
	}
	auth, err := spotifyAuth()
	if err != nil {
		return
	}
	url_ := fmt.Sprintf("https://api.spotify.com/v1/%s/%s", typeRoute, id)
	header := map[string]string{
		"Authorization": auth.TokenType + " " + auth.AccessToken,
	}
	resultJson, err := http_util.Get(url_, header)
	if err != nil {
		return
	}
	switch urlType {
	case spotifyTrack:
		var result spotifyTrackResult
		err = json.Unmarshal(resultJson, &result)
		if err != nil {
			return
		}
		name = result.Name
		link = result.ExternalUrls.Spotify
		searchStr = []string{buildTrackSearchStr(result)}
		return
	case spotifyAlbum:
		var result spotifyAlbumResult
		err = json.Unmarshal(resultJson, &result)
		if err != nil {
			return
		}
		next := result.Tracks.Next
		for next != nil {
			var nextResult spotifyAlbumPagination
			resultJson, err = http_util.Get(*next, header)
			if err != nil {
				return
			}
			err = json.Unmarshal(resultJson, &nextResult)
			if err != nil {
				return
			}
			result.Tracks.Items = append(result.Tracks.Items, nextResult.Items...)
			next = nextResult.Next
		}
		name = result.Name
		link = result.ExternalUrls.Spotify
		searchStr = make([]string, len(result.Tracks.Items))
		for i, track := range result.Tracks.Items {
			searchStr[i] = buildTrackSearchStr(track)
		}
		return
	case spotifyPlaylist:
		var result spotifyPlaylistResult
		err = json.Unmarshal(resultJson, &result)
		if err != nil {
			return
		}
		next := result.Tracks.Next
		for next != nil {
			var nextResult spotifyPlaylistPagination
			resultJson, err = http_util.Get(*next, header)
			if err != nil {
				return
			}
			err = json.Unmarshal(resultJson, &nextResult)
			if err != nil {
				return
			}
			result.Tracks.Items = append(result.Tracks.Items, nextResult.Items...)
			next = nextResult.Next
		}
		name = result.Name
		link = result.ExternalUrls.Spotify
		searchStr = make([]string, len(result.Tracks.Items))
		for i, track := range result.Tracks.Items {
			searchStr[i] = buildTrackSearchStr(track.Track)
		}
		return
	}
	return
}

func buildTrackSearchStr(track spotifyTrackResult) string {
	mapArtist := util.Map(track.Artists, func(_ int, artist spotifyArtistResult) string { return artist.Name })
	return strings.Join(mapArtist, ", ") + " - " + track.Name
}

type spotifyAlbumResult struct {
	Name         string                 `json:"name"`
	ExternalUrls externalUrlsResult     `json:"external_urls"`
	Tracks       spotifyAlbumPagination `json:"tracks"`
}

type spotifyAlbumPagination struct {
	Items    []spotifyTrackResult `json:"items"`
	Limit    int                  `json:"limit"`
	Next     *string              `json:"next"`
	Offset   int                  `json:"offset"`
	Previous *string              `json:"previous"`
	Total    int                  `json:"total"`
}

type spotifyPlaylistResult struct {
	Name         string                    `json:"name"`
	ExternalUrls externalUrlsResult        `json:"external_urls"`
	Tracks       spotifyPlaylistPagination `json:"tracks"`
}

type spotifyPlaylistPagination struct {
	Items []struct {
		Track spotifyTrackResult `json:"track"`
	} `json:"items"`
	Limit    int     `json:"limit"`
	Next     *string `json:"next"`
	Offset   int     `json:"offset"`
	Previous *string `json:"previous"`
	Total    int     `json:"total"`
}

type spotifyTrackResult struct {
	Artists      []spotifyArtistResult `json:"artists"`
	Name         string                `json:"name"`
	ExternalUrls externalUrlsResult    `json:"external_urls"`
}

type spotifyArtistResult struct {
	Name string `json:"name"`
}

type externalUrlsResult struct {
	Spotify string `json:"spotify"`
}
