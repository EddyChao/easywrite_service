package tools

import (
	"fmt"
	"testing"
)

func TestGetDuration(t *testing.T) {
	audioFilePath := "C:\\Users\\wilinz\\Downloads\\tts.mp3"
	duration, err := GetAudioDuration(audioFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("音频时长:", duration)
}
