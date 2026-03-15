package parser

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name       string
		rawURL     string
		wantOwner  string
		wantRepo   string
		wantType   string
		wantNumber int
		wantErr    bool
	}{
		// Issue URL
		{
			name:       "parse issue URL",
			rawURL:     "https://github.com/gorilla/mux/issues/123",
			wantOwner:  "gorilla",
			wantRepo:   "mux",
			wantType:   "issue",
			wantNumber: 123,
			wantErr:    false,
		},
		// PR URL
		{
			name:       "parse PR URL",
			rawURL:     "https://github.com/gorilla/mux/pull/456",
			wantOwner:  "gorilla",
			wantRepo:   "mux",
			wantType:   "pr",
			wantNumber: 456,
			wantErr:    false,
		},
		// Discussion URL
		{
			name:       "parse discussion URL",
			rawURL:     "https://github.com/gorilla/mux/discussions/789",
			wantOwner:  "gorilla",
			wantRepo:   "mux",
			wantType:   "discussion",
			wantNumber: 789,
			wantErr:    false,
		},
		// Invalid URL - not github.com
		{
			name:       "invalid URL - not github.com",
			rawURL:     "https://gitlab.com/gorilla/mux/issues/123",
			wantOwner:  "",
			wantRepo:   "",
			wantType:   "",
			wantNumber: 0,
			wantErr:    true,
		},
		// Invalid URL - wrong path
		{
			name:       "invalid URL - wrong path",
			rawURL:     "https://github.com/gorilla/mux/abc/123",
			wantOwner:  "",
			wantRepo:   "",
			wantType:   "",
			wantNumber: 0,
			wantErr:    true,
		},
		// Invalid URL - empty
		{
			name:       "invalid URL - empty",
			rawURL:     "",
			wantOwner:  "",
			wantRepo:   "",
			wantType:   "",
			wantNumber: 0,
			wantErr:    true,
		},
		// GitHub Enterprise URL
		{
			name:       "parse enterprise URL",
			rawURL:     "https://github.mycompany.com/owner/repo/issues/1",
			wantOwner:  "owner",
			wantRepo:   "repo",
			wantType:   "issue",
			wantNumber: 1,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, issueType, number, err := ParseURL(tt.rawURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if owner != tt.wantOwner {
				t.Errorf("ParseURL() owner = %v, want %v", owner, tt.wantOwner)
			}
			if repo != tt.wantRepo {
				t.Errorf("ParseURL() repo = %v, want %v", repo, tt.wantRepo)
			}
			if issueType != tt.wantType {
				t.Errorf("ParseURL() issueType = %v, want %v", issueType, tt.wantType)
			}
			if number != tt.wantNumber {
				t.Errorf("ParseURL() number = %v, want %v", number, tt.wantNumber)
			}
		})
	}
}

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		want    bool
	}{
		{
			name:   "valid issue URL",
			rawURL: "https://github.com/gorilla/mux/issues/123",
			want:   true,
		},
		{
			name:   "valid PR URL",
			rawURL: "https://github.com/gorilla/mux/pull/456",
			want:   true,
		},
		{
			name:   "valid discussion URL",
			rawURL: "https://github.com/gorilla/mux/discussions/789",
			want:   true,
		},
		{
			name:   "invalid - not github.com",
			rawURL: "https://gitlab.com/gorilla/mux/issues/123",
			want:   false,
		},
		{
			name:   "invalid - wrong path",
			rawURL: "https://github.com/gorilla/mux/abc/123",
			want:   false,
		},
		{
			name:   "invalid - empty",
			rawURL: "",
			want:   false,
		},
		{
			name:   "valid enterprise URL",
			rawURL: "https://github.mycompany.com/owner/repo/issues/1",
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidURL(tt.rawURL); got != tt.want {
				t.Errorf("IsValidURL(%q) = %v, want %v", tt.rawURL, got, tt.want)
			}
		})
	}
}
