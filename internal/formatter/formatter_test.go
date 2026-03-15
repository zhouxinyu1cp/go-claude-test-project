package formatter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no special characters",
			input:    "simpletitle",
			expected: "simpletitle",
		},
		{
			name:     "slash replaced",
			input:    "title/with/slash",
			expected: "title-with-slash",
		},
		{
			name:     "colon replaced",
			input:    "title:with:colon",
			expected: "title-with-colon",
		},
		{
			name:     "spaces replaced",
			input:    "title with spaces",
			expected: "title-with-spaces",
		},
		{
			name:     "multiple special characters",
			input:    "title/with:multiple special/chars",
			expected: "title-with-multiple-special-chars",
		},
		{
			name:     "only special characters",
			input:    "///",
			expected: "---",
		},
		{
			name:     "trailing spaces",
			input:    "title  ",
			expected: "title--",
		},
		{
			name:     "leading spaces",
			input:    "  title",
			expected: "--title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeFilename(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeFilename(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGenerateFilename(t *testing.T) {
	tests := []struct {
		name       string
		owner      string
		repo       string
		issueType  string
		number     int
		title      string
		wantPrefix string
		wantExt    string
	}{
		{
			name:       "issue filename",
			owner:      "gorilla",
			repo:      "mux",
			issueType: "issue",
			number:    123,
			title:     "Test Issue",
			wantPrefix: "gorilla-mux-123-",
			wantExt:    ".md",
		},
		{
			name:       "pr filename",
			owner:      "gorilla",
			repo:      "mux",
			issueType: "pr",
			number:    456,
			title:     "Fix bug",
			wantPrefix: "gorilla-mux-pr-456-",
			wantExt:    ".md",
		},
		{
			name:       "discussion filename",
			owner:      "gorilla",
			repo:      "mux",
			issueType: "discussion",
			number:    789,
			title:     "Question",
			wantPrefix: "gorilla-mux-discussion-789-",
			wantExt:    ".md",
		},
		{
			name:       "title with special chars",
			owner:      "owner",
			repo:      "repo",
			issueType: "issue",
			number:    1,
			title:     "Fix: bug in /api",
			wantPrefix: "owner-repo-1-",
			wantExt:    ".md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateFilename(tt.owner, tt.repo, tt.issueType, tt.number, tt.title)
			if len(result) == 0 {
				t.Errorf("GenerateFilename() returned empty string")
				return
			}
			if result[:len(tt.wantPrefix)] != tt.wantPrefix {
				t.Errorf("GenerateFilename() prefix = %q, want prefix %q", result[:len(tt.wantPrefix)], tt.wantPrefix)
			}
			if result[len(result)-len(tt.wantExt):] != tt.wantExt {
				t.Errorf("GenerateFilename() extension = %q, want extension %q", result[len(result)-len(tt.wantExt):], tt.wantExt)
			}
		})
	}
}

func TestResolveFilenameConflict(t *testing.T) {
	tests := []struct {
		name         string
		filename     string
		existing     []string
		wantConflict bool
		wantResult   string
	}{
		{
			name:         "no conflict",
			filename:     "newfile.md",
			existing:     []string{"existing.md"},
			wantConflict: false,
			wantResult:   "newfile.md",
		},
		{
			name:         "file exists",
			filename:     "file.md",
			existing:     []string{"file.md"},
			wantConflict: true,
			wantResult:   "file_1.md",
		},
		{
			name:         "multiple conflicts",
			filename:     "file.md",
			existing:     []string{"file.md", "file_1.md", "file_2.md"},
			wantConflict: true,
			wantResult:   "file_3.md",
		},
		{
			name:         "gap in numbering",
			filename:     "file.md",
			existing:     []string{"file.md", "file_2.md"},
			wantConflict: true,
			wantResult:   "file_1.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 每个子测试使用独立的临时目录
			tmpDir := t.TempDir()

			// 创建已存在的文件
			for _, fname := range tt.existing {
				fpath := filepath.Join(tmpDir, fname)
				if err := os.WriteFile(fpath, []byte("test"), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			}

			result := ResolveFilenameConflict(tmpDir, tt.filename)

			// 检查结果是否正确
			fpath := filepath.Join(tmpDir, result)
			if _, err := os.Stat(fpath); err == nil {
				t.Errorf("ResolveFilenameConflict() created file that should not exist yet")
			}

			if result != tt.wantResult {
				t.Errorf("ResolveFilenameConflict() = %q, want %q", result, tt.wantResult)
			}
		})
	}
}
