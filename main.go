package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"
)

// Version is set at build time via -ldflags
var Version = "dev"

// Verbose controls debug logging
var Verbose = false

// DebugLog prints a message if verbose mode is enabled
func DebugLog(format string, args ...interface{}) {
	if Verbose {
		log.Printf("[DEBUG] "+format, args...)
	}
}

func main() {
	// Define flags
	configPath := flag.String("config", "./config.yaml", "Path to YAML config file")
	configPathShort := flag.String("c", "", "Path to YAML config file (shorthand)")
	errorExit := flag.Bool("error-exit", false, "Exit with code 1 on file/category not found")
	errorExitShort := flag.Bool("e", false, "Exit with code 1 on file/category not found (shorthand)")
	showVersion := flag.Bool("version", false, "Show version")
	showVersionShort := flag.Bool("v", false, "Show version (shorthand)")
	verbose := flag.Bool("verbose", false, "Enable verbose/debug output")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "say-mi %s\n\n", Version)
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <category>\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Play a random sound from the specified category.")
		fmt.Fprintln(os.Stderr, "Categories use dot notation for nested paths (e.g., kyoko.greeting).")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintf(os.Stderr, "  %s greeting\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s kyoko.greeting\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -c ./voices.yaml permission\n", os.Args[0])
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
	if *verbose {
		Verbose = true
	}

	DebugLog("Starting say-mi %s", Version)

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
	DebugLog("Category requested: %s", category)

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Load config
	DebugLog("Loading config from: %s", *configPath)
	config, err := LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}
	DebugLog("Config loaded successfully")

	// Find audio files for category (supports dot notation)
	audioFiles, exists := config.FindCategory(category)
	DebugLog("Found %d audio files for category '%s' (exists: %v)", len(audioFiles), category, exists)
	var audioPath string

	if !exists || len(audioFiles) == 0 {
		// Use fallback if available
		if config.Fallback != "" {
			audioPath = config.Fallback
			DebugLog("Using fallback audio: %s", audioPath)
		} else {
			// No fallback, handle based on error-exit flag
			DebugLog("No audio files found and no fallback configured")
			if *errorExit {
				os.Exit(1)
			}
			os.Exit(0)
		}
	} else {
		// Pick random audio from category
		audioPath = audioFiles[rand.Intn(len(audioFiles))]
		DebugLog("Selected audio file: %s", audioPath)
	}

	// Check if file exists
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: audio file not found: %s\n", audioPath)
		if *errorExit {
			os.Exit(1)
		}
		os.Exit(0)
	}
	DebugLog("Audio file exists: %s", audioPath)

	// Play the audio
	DebugLog("Starting audio playback...")
	if err := PlayAudio(audioPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error playing audio: %v\n", err)
		if *errorExit {
			os.Exit(1)
		}
	}
	DebugLog("Audio playback completed")

	os.Exit(0)
}
