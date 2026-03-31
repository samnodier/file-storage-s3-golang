package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}
	var data ffprobeOutput
	if err := json.Unmarshal(stdout.Bytes(), &data); err != nil {
		return "", err
	}
	if len(data.Streams) == 0 {
		return "", fmt.Errorf("no streams found in file")
	}

	w := float64(data.Streams[0].Width)
	h := float64(data.Streams[0].Height)

	if h == 0 {
		return "", fmt.Errorf("invalid height")
	}

	ratio := w / h
	epsilon := 0.02

	if ratio > (1.777-epsilon) && ratio < (1.777+epsilon) {
		return "16:9", nil
	}
	if ratio > (0.562-epsilon) && ratio < (0.562+epsilon) {
		return "9:16", nil
	}
	return "other", nil
}

func processVideoForFastStart(filePath string) (string, error) {
	outputPath := filePath + ".processing"
	cmd := exec.Command("ffmpeg",
		"-i", filePath,
		"-c", "copy",
		"-movflags", "faststart",
		"-f", "mp4",
		outputPath,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg error: %s", stderr.String())
	}
	return outputPath, nil
}

type ffprobeOutput struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}
