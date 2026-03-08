package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"
)

// Version is set at build time via -ldflags
var Version = "dev"

func main() {
	// Define flags
	configPath := flag.String("config", "./config.yaml", "Path to YAML config file")
	configPathShort := flag.String("c", "", "Path to YAML config file (shorthand)")
	errorExit := flag.Bool("error-exit", false, "Exit with code 1 on file/category not found")
	errorExitShort := flag.Bool("e", false, "Exit with code 1 on file/category not found (shorthand)")
	showVersion := flag.Bool("version", false, "Show version")
	showVersionShort := flag.Bool("v", false, "Show version (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "say-mi %s\n\n", Version)
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <category>\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Play a random sound from the specified category.")
		fmt.Fprintln(os.Stderr, "Categories use dot notation for nested paths (e.g., female.sayako.email).")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintf(os.Stderr, "  %s hello\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s female.sayako.email\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -c ./sounds.yaml permission\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -e nonexistent  # exits 1\n", os.Args[0])
	}

	flag.Parse()

	// Merge short and long flags
	if *configPathShort != "" {
		*configPath = *configPathShort
	}
	if *errorExitShort {
		*errorExit = true
	}
	if *showVersionShort {
		*showVersion = true
	}

	// Handle --version
	if *showVersion {
		fmt.Printf("say-mi %s (%s/%s)\n", Version, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	// Get category from positional args
	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	category := args[0]

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Load config
	config, err := LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Find audio files for category (supports dot notation)
	audioFiles, exists := config.FindCategory(category)
	var audioPath string

	if !exists || len(audioFiles) == 0 {
		// Use fallback if available
		if config.Fallback != "" {
			audioPath = config.Fallback
		} else {
			// No fallback, handle based on error-exit flag
			if *errorExit {
				os.Exit(1)
			}
			os.Exit(0)
		}
	} else {
		// Pick random audio from category
		audioPath = audioFiles[rand.Intn(len(audioFiles))]
	}

	// Play the audio
	if err := PlayAudio(audioPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error playing audio: %v\n", err)
		if *errorExit {
			os.Exit(1)
		}
	}

	os.Exit(0)
}
