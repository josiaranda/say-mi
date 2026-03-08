package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/flac"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/vorbis"
	"github.com/gopxl/beep/wav"
)

// PlayAudio plays an audio file (mp3, wav, ogg, or flac) and blocks until completion
func PlayAudio(filePath string) error {
	// Open the file
	DebugLog("Opening audio file: %s", filePath)
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open audio file: %w", err)
	}
	defer f.Close()

	// Decode based on file extension
	var streamer beep.StreamSeekCloser
	var format beep.Format

	ext := strings.ToLower(filepath.Ext(filePath))
	DebugLog("Audio format detected: %s", ext)

	switch ext {
	case ".mp3":
		streamer, format, err = mp3.Decode(f)
	case ".wav":
		streamer, format, err = wav.Decode(f)
	case ".ogg":
		streamer, format, err = vorbis.Decode(f)
	case ".flac":
		streamer, format, err = flac.Decode(f)
	default:
		return fmt.Errorf("unsupported audio format: %s (supported: mp3, wav, ogg, flac)", ext)
	}

	if err != nil {
		return fmt.Errorf("failed to decode audio: %w", err)
	}
	defer streamer.Close()

	DebugLog("Audio decoded: SampleRate=%d, Channels=%d, Precision=%d, Length=%d samples",
		format.SampleRate, format.NumChannels, format.Precision, streamer.Len())

	// Initialize speaker with the audio's sample rate
	bufferSize := format.SampleRate.N(time.Second / 10)
	DebugLog("Initializing speaker with buffer size: %d", bufferSize)
	if err := speaker.Init(format.SampleRate, bufferSize); err != nil {
		return fmt.Errorf("failed to initialize speaker: %w", err)
	}
	defer speaker.Close()
	DebugLog("Speaker initialized successfully")

	// Create a channel to signal playback completion
	done := make(chan bool)

	// Play the audio
	DebugLog("Starting playback...")
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		DebugLog("Playback callback triggered")
		done <- true
	})))

	// Wait for playback to complete
	<-done
	DebugLog("Playback finished")

	return nil
}
