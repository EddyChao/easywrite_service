package tools

import (
	"errors"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetAudioDuration(audioFilePath string) (time.Duration, error) {
	cmd := exec.Command("ffmpeg", "-i", audioFilePath)
	output, _ := cmd.CombinedOutput()
	duration, err := parseDuration(output)
	return duration, err
}

func parseDuration(output []byte) (time.Duration, error) {
	outputStr := string(output)

	durationLine := ""
	for _, line := range strings.Split(outputStr, "\n") {
		if strings.Contains(line, "Duration") {
			durationLine = line
			break
		}
	}

	if durationLine == "" {
		return 0, errors.New("duration not found in output")
	}

	duration := extractDuration(durationLine)

	return duration, nil
}

func extractDuration(line string) time.Duration {
	re := regexp.MustCompile(`Duration: (\d+):(\d+):(\d+)\.(\d+)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) == 5 {
		hours, _ := strconv.Atoi(matches[1])
		minutes, _ := strconv.Atoi(matches[2])
		seconds, _ := strconv.Atoi(matches[3])
		millisecond, _ := strconv.Atoi(matches[4])
		durationMillisecond := (hours*3600+minutes*60+seconds)*1000 + millisecond
		return time.Duration(durationMillisecond) * time.Millisecond
	}
	return 0
}
