package utils

import "strings"

func IsAudioFileType(path string) bool {
	if strings.HasSuffix(path, ".mp3") || strings.HasSuffix(path, ".ogg") ||
		strings.HasSuffix(path, ".flac") || strings.HasSuffix(path, ".wav") {
		return true
	} else {
		return false
	}
}
