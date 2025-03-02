package main

import (
	"fmt"
	"io"
	"os"

	"github.com/hajimehoshi/go-mp3"
)

// Function to detect silence (dummy implementation, needs amplitude analysis)
func detectSilence(data []int16, sampleRate int, silenceThreshold int16, minSilenceDuration int) []int {
	silenceStart := -1
	silencePositions := []int{}

	for i := 0; i < len(data); i++ {
		if abs(data[i]) < silenceThreshold {
			if silenceStart == -1 {
				silenceStart = i
			}
			if i-silenceStart >= minSilenceDuration*sampleRate {
				silencePositions = append(silencePositions, silenceStart)
				silenceStart = -1
			}
		} else {
			silenceStart = -1
		}
	}
	return silencePositions
}

// Helper function to get absolute value of int16
func abs(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}

// Function to trim leading/trailing silence
func trimSilence(data []int16, silenceThreshold int16) []int16 {
	start, end := 0, len(data)-1

	// Trim leading silence
	for start < len(data) && abs(data[start]) < silenceThreshold {
		start++
	}

	// Trim trailing silence
	for end > start && abs(data[end]) < silenceThreshold {
		end--
	}

	return data[start : end+1]
}

func main() {
	// Open MP3 file
	file, err := os.Open("input.mp3")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Decode MP3 to PCM
	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		fmt.Println("Error decoding MP3:", err)
		return
	}

	// Read PCM data
	var pcmData []int16
	buf := make([]byte, 1024)
	for {
		n, err := decoder.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading PCM data:", err)
			return
		}
		for i := 0; i < n; i += 2 {
			pcmData = append(pcmData, int16(buf[i])|(int16(buf[i+1])<<8))
		}
	}

	// Detect silence positions (threshold -40dB, 2s silence)
	silencePositions := detectSilence(pcmData, 44100, 500, 2)
	fmt.Println("Silence Positions:", silencePositions)

	// Trim silence from each segment
	trimmedAudio := trimSilence(pcmData, 500)

	// TODO: Repeat each trimmed segment 3x and re-encode as MP3

	fmt.Println("Processing complete!")
}
