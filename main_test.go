package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestMainIntegration tests the CLI behavior by running the compiled binary
func TestMainIntegration(t *testing.T) {
	// Build the binary with test tag to avoid ALSA dependency
	binPath := filepath.Join(t.TempDir(), "say-mi-test")
	buildCmd := exec.Command("go", "build", "-tags", "test", "-o", binPath, ".")
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, output)
	}

	// Create test config (no fallback so missing categories actually fail)
	testConfig := `voices:
  hello:
    - "test.mp3"
  male:
    sato:
      hello:
        - "test.mp3"
`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.yaml")
	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	tests := []struct {
		name     string
		args     []string
		wantExit int
	}{
		{
			name:     "no arguments shows usage",
			args:     []string{},
			wantExit: 1,
		},
		{
			name:     "missing category without error-exit flag exits 0",
			args:     []string{"-c", configPath, "nonexistent"},
			wantExit: 0,
		},
		{
			name:     "missing category with error-exit flag exits 1",
			args:     []string{"-c", configPath, "-e", "nonexistent"},
			wantExit: 1,
		},
		{
			name:     "nested missing category with error-exit flag exits 1",
			args:     []string{"-c", configPath, "-e", "male.nonexistent"},
			wantExit: 1,
		},
		{
			name:     "partial path (not a leaf) exits 1 with -e",
			args:     []string{"-c", configPath, "-e", "male"},
			wantExit: 1,
		},
		{
			name:     "config file not found",
			args:     []string{"-c", "/nonexistent/config.yaml", "hello"},
			wantExit: 1,
		},
		{
			name:     "long form flags",
			args:     []string{"--config", configPath, "--error-exit", "nonexistent"},
			wantExit: 1,
		},
		{
			name: "nested category found",
			args: []string{"-c", configPath, "male.sato.hello"},
			// Note: Will try to play audio, which may fail in headless env
			// But category lookup should succeed
			wantExit: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binPath, tt.args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			exitCode := 0
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCode = exitErr.ExitCode()
				} else {
					t.Fatalf("Failed to run command: %v", err)
				}
			}

			if exitCode != tt.wantExit {
				t.Errorf("Exit code = %d, want %d", exitCode, tt.wantExit)
			}
		})
	}
}

func TestHelpFlag(t *testing.T) {
	binPath := filepath.Join(t.TempDir(), "say-mi-help-test")
	buildCmd := exec.Command("go", "build", "-tags", "test", "-o", binPath, ".")
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, output)
	}

	cmd := exec.Command(binPath, "-h")
	output, err = cmd.CombinedOutput()

	// -h should exit with 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != 0 {
				t.Errorf("-h flag should exit with 0, got %d", exitErr.ExitCode())
			}
		}
	}

	// Check that help text contains expected content
	helpText := string(output)
	expectedStrings := []string{
		"Usage:",
		"category",
		"-config",
		"-error-exit",
		"kyoko.greeting",
	}

	for _, expected := range expectedStrings {
		if !containsStr(helpText, expected) {
			t.Errorf("Help output missing expected string: %q", expected)
		}
	}
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
