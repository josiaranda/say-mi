package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPlayAudio_FileNotFound(t *testing.T) {
	err := PlayAudio("/nonexistent/audio.mp3")
	if err == nil {
		t.Error("PlayAudio() expected error for non-existent file")
	}
}

func TestPlayAudio_UnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file with unsupported extension
	tmpFile := filepath.Join(tmpDir, "audio.txt")
	if err := os.WriteFile(tmpFile, []byte("not audio"), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	err := PlayAudio(tmpFile)
	if err == nil {
		t.Error("PlayAudio() expected error for unsupported format")
	}
}

func TestPlayAudio_InvalidMP3(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a fake mp3 file with invalid content
	tmpFile := filepath.Join(tmpDir, "fake.mp3")
	if err := os.WriteFile(tmpFile, []byte("not a real mp3"), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	err := PlayAudio(tmpFile)
	if err == nil {
		t.Error("PlayAudio() expected error for invalid mp3 content")
	}
}

func TestPlayAudio_InvalidWAV(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a fake wav file with invalid content
	tmpFile := filepath.Join(tmpDir, "fake.wav")
	if err := os.WriteFile(tmpFile, []byte("not a real wav"), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	err := PlayAudio(tmpFile)
	if err == nil {
		t.Error("PlayAudio() expected error for invalid wav content")
	}
}

func TestPlayAudio_InvalidOGG(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a fake ogg file with invalid content
	tmpFile := filepath.Join(tmpDir, "fake.ogg")
	if err := os.WriteFile(tmpFile, []byte("not a real ogg"), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	err := PlayAudio(tmpFile)
	if err == nil {
		t.Error("PlayAudio() expected error for invalid ogg content")
	}
}

func TestPlayAudio_RealFile(t *testing.T) {
	// Skip if test.mp3 doesn't exist
	if _, err := os.Stat("test.mp3"); os.IsNotExist(err) {
		t.Skip("test.mp3 not found, skipping real audio test")
	}

	// This test requires audio output capability
	err := PlayAudio("test.mp3")
	if err != nil {
		t.Logf("PlayAudio with real file: %v (may fail in headless/WSL env)", err)
	}
}
