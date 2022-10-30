package yt_util

type Playlist struct {
	Entries    []Video `json:"entries"`
	Title      string  `json:"title"`
	WebpageURL string  `json:"webpage_url"`
}

type PlaylistInfo struct {
	Title      string
	WebpageURL string
}

type Video struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Duration float64 `json:"duration"`
	// for check if bestaudio exist "none" if not
	Acodec string `json:"acodec"`
}
