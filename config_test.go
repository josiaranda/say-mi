package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantErr  bool
		fallback string
	}{
		{
			name: "valid config with fallback",
			content: `fallback: "audio/error.mp3"
categories:
  hello:
    - "audio/hello.mp3"
`,
			wantErr:  false,
			fallback: "audio/error.mp3",
		},
		{
			name: "valid config without fallback",
			content: `categories:
  hello:
    - "audio/hello.mp3"
`,
			wantErr:  false,
			fallback: "",
		},
		{
			name: "empty config",
			content: `categories: {}
`,
			wantErr:  false,
			fallback: "",
		},
		{
			name:    "invalid yaml",
			content: `invalid: [yaml: content`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "config.yaml")
			if err := os.WriteFile(tmpFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}

			config, err := LoadConfig(tmpFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && config.Fallback != tt.fallback {
				t.Errorf("LoadConfig() fallback = %v, want %v", config.Fallback, tt.fallback)
			}
		})
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("LoadConfig() expected error for non-existent file")
	}
}

func TestLoadConfig_RealFile(t *testing.T) {
	// Test loading the actual config.yaml
	config, err := LoadConfig("config.yaml")
	if err != nil {
		t.Fatalf("Failed to load config.yaml: %v", err)
	}

	if config.Fallback == "" {
		t.Error("Expected fallback to be set in config.yaml")
	}

	if config.Categories == nil {
		t.Error("Expected categories to be set in config.yaml")
	}
}

func TestFindCategory(t *testing.T) {
	content := `fallback: "test.mp3"
categories:
  hello:
    - "test.mp3"
  male:
    sato:
      hello:
        - "test.mp3"
      email:
        - "test.mp3"
    tanaka:
      hello:
        - "test.mp3"
  female:
    sayako:
      email:
        - "test.mp3"
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	config, err := LoadConfig(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		wantLen  int
		wantFile string
		exists   bool
	}{
		{
			name:     "flat category",
			path:     "hello",
			wantLen:  1,
			wantFile: "test.mp3",
			exists:   true,
		},
		{
			name:     "nested two levels - male.sato.hello",
			path:     "male.sato.hello",
			wantLen:  1,
			wantFile: "test.mp3",
			exists:   true,
		},
		{
			name:     "nested two levels - male.sato.email",
			path:     "male.sato.email",
			wantLen:  1,
			wantFile: "test.mp3",
			exists:   true,
		},
		{
			name:     "nested two levels - female.sayako.email",
			path:     "female.sayako.email",
			wantLen:  1,
			wantFile: "test.mp3",
			exists:   true,
		},
		{
			name:   "nonexistent top level",
			path:   "nonexistent",
			exists: false,
		},
		{
			name:   "nonexistent nested path",
			path:   "male.nonexistent.hello",
			exists: false,
		},
		{
			name:   "partial path (not a leaf)",
			path:   "male.sato",
			exists: false,
		},
		{
			name:   "empty path",
			path:   "",
			exists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, exists := config.FindCategory(tt.path)

			if exists != tt.exists {
				t.Errorf("FindCategory(%q) exists = %v, want %v", tt.path, exists, tt.exists)
				return
			}

			if tt.exists {
				if len(files) != tt.wantLen {
					t.Errorf("FindCategory(%q) got %d files, want %d", tt.path, len(files), tt.wantLen)
				}

				if tt.wantFile != "" && len(files) > 0 {
					// Check if the expected file is in the list
					found := false
					for _, f := range files {
						if f == tt.wantFile {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("FindCategory(%q) missing expected file %q", tt.path, tt.wantFile)
					}
				}
			}
		})
	}
}

func TestFindCategory_RealConfig(t *testing.T) {
	config, err := LoadConfig("config.yaml")
	if err != nil {
		t.Fatalf("Failed to load config.yaml: %v", err)
	}

	// Test flat category
	files, exists := config.FindCategory("hello")
	if !exists {
		t.Error("FindCategory(\"hello\") should exist")
	}
	if len(files) == 0 {
		t.Error("FindCategory(\"hello\") should have files")
	}

	// Test nested category
	files, exists = config.FindCategory("male.sato.hello")
	if !exists {
		t.Error("FindCategory(\"male.sato.hello\") should exist")
	}
	if len(files) == 0 {
		t.Error("FindCategory(\"male.sato.hello\") should have files")
	}

	// Test deep nested category
	files, exists = config.FindCategory("female.sayako.email")
	if !exists {
		t.Error("FindCategory(\"female.sayako.email\") should exist")
	}
	if len(files) == 0 {
		t.Error("FindCategory(\"female.sayako.email\") should have files")
	}
}

func TestFindCategory_EmptyCategories(t *testing.T) {
	content := `categories: {}`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	config, err := LoadConfig(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	_, exists := config.FindCategory("anything")
	if exists {
		t.Error("FindCategory() should return false for empty categories")
	}
}

func TestExtractStringSlice(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		want   []string
		wantOk bool
	}{
		{
			name:   "[]interface{} with strings",
			input:  []interface{}{"a.mp3", "b.wav", "c.ogg"},
			want:   []string{"a.mp3", "b.wav", "c.ogg"},
			wantOk: true,
		},
		{
			name:   "[]string directly",
			input:  []string{"x.mp3", "y.wav"},
			want:   []string{"x.mp3", "y.wav"},
			wantOk: true,
		},
		{
			name:   "[]interface{} with mixed types (only strings extracted)",
			input:  []interface{}{"a.mp3", 123, "b.wav"},
			want:   []string{"a.mp3", "b.wav"},
			wantOk: true,
		},
		{
			name:   "empty slice",
			input:  []interface{}{},
			want:   nil,
			wantOk: false,
		},
		{
			name:   "non-slice",
			input:  "not a slice",
			want:   nil,
			wantOk: false,
		},
		{
			name:   "map instead of slice",
			input:  map[string]string{"key": "value"},
			want:   nil,
			wantOk: false,
		},
		{
			name: "[]interface{} with objects containing audio field",
			input: []interface{}{
				map[string]interface{}{"text": "Hello", "audio": "kyoko/say/greeting/0.mp3"},
				map[string]interface{}{"text": "Hi", "audio": "kyoko/say/greeting/1.mp3"},
			},
			want:   []string{"kyoko/say/greeting/0.mp3", "kyoko/say/greeting/1.mp3"},
			wantOk: true,
		},
		{
			name: "[]interface{} with mixed objects and strings",
			input: []interface{}{
				"plain_string.mp3",
				map[string]interface{}{"text": "Hello", "audio": "obj.mp3"},
			},
			want:   []string{"plain_string.mp3", "obj.mp3"},
			wantOk: true,
		},
		{
			name: "[]interface{} with objects missing or empty audio",
			input: []interface{}{
				map[string]interface{}{"text": "No audio"},
				map[string]interface{}{"text": "Empty audio", "audio": ""},
			},
			want:   nil,
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := extractStringSlice(tt.input)
			if ok != tt.wantOk {
				t.Errorf("extractStringSlice() ok = %v, want %v", ok, tt.wantOk)
				return
			}

			if tt.wantOk {
				if len(got) != len(tt.want) {
					t.Errorf("extractStringSlice() got %v, want %v", got, tt.want)
					return
				}
				for i := range got {
					if got[i] != tt.want[i] {
						t.Errorf("extractStringSlice()[%d] = %v, want %v", i, got[i], tt.want[i])
					}
				}
			}
		})
	}
}
