package util

import (
	"fmt"
	"sort"

	"github.com/kkdai/youtube/v2"
)

func GetBestAudioYT(query youtube.FormatList) (*youtube.Format, error) {
	filtered := youtube.FormatList{}
	for _, f := range query {
		if f.Width == 0 && f.Height == 0 {
			filtered = append(filtered, f)
		}
	}
	if len(filtered) == 0 {
		return nil, fmt.Errorf("No audio found")
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Bitrate > filtered[j].Bitrate
	})
	return &filtered[0], nil
}
