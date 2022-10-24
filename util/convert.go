package util

import (
	"fmt"
	"strconv"
	"strings"
)

func ConvertVidLengthToSeconds(vidLength string) (seconds uint32) {
	split := strings.Split(vidLength, ":")
	for i := 0; i < len(split); i++ {
		num, _ := strconv.ParseInt(split[i], 10, 64)
		mul := uint32(1)
		for p := 0; p < len(split)-i-1; p++ {
			mul *= 60
		}
		seconds += uint32(num) * mul
	}
	return
}

func ConvertSecondsToVidLength(seconds uint32) (vidLength string) {
	if seconds < 60 {
		return fmt.Sprintf("0:%02d", seconds)
	}
	for i := 0; i <= 2; i++ {
		v := seconds
		if i != 2 {
			v = v % 60
		}
		str := fmt.Sprintf("%02d", v)
		vidLength = ":" + str + vidLength
		seconds = seconds / 60
		if seconds == 0 {
			break
		}
	}
	vidLength = vidLength[1:]
	if vidLength[0] == '0' {
		vidLength = vidLength[1:]
	}
	return
}
