package fetcher

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

// mockTransport 用于模拟 GitHub API 响应
type mockTransport struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func TestFetcher_FetchIssue(t *testing.T) {
	tests := []struct {
		name       string
		owner      string
		repo       string
		number     int
		setupMock  func() *http.Client
		wantTitle  string
		wantAuthor string
		wantErr    bool
	}{
		{
			name:   "fetch issue successfully",
			owner:  "gorilla",
			repo:   "mux",
			number: 123,
			setupMock: func() *http.Client {
				return &http.Client{
					Transport: &mockTransport{
						roundTripFunc: func(req *http.Request) (*http.Response, error) {
							if req.URL.Path == "/repos/gorilla/mux/issues/123" {
								return &http.Response{
									StatusCode: http.StatusOK,
									Body:       io.NopCloser(bytes.NewReader([]byte(`{"title":"Test Issue","user":{"login":"testuser"},"body":"Issue body","state":"open"}`))),
									Header:     http.Header{"Content-Type": []string{"application/json"}},
								}, nil
							}
							if req.URL.Path == "/repos/gorilla/mux/issues/123/comments" {
								return &http.Response{
									StatusCode: http.StatusOK,
									Body:       io.NopCloser(bytes.NewReader([]byte(`[{"body":"comment1","created_at":"2024-01-01T00:00:00Z","user":{"login":"user1"}}]`))),
									Header:     http.Header{"Content-Type": []string{"application/json"}},
								}, nil
							}
							return &http.Response{StatusCode: http.StatusNotFound}, nil
						},
					},
				}
			},
			wantTitle:  "Test Issue",
			wantAuthor: "testuser",
			wantErr:    false,
		},
		{
			name:   "fetch issue 404",
			owner:  "gorilla",
			repo:   "mux",
			number: 999999,
			setupMock: func() *http.Client {
				return &http.Client{
					Transport: &mockTransport{
						roundTripFunc: func(req *http.Request) (*http.Response, error) {
							return &http.Response{StatusCode: http.StatusNotFound}, nil
						},
					},
				}
			},
			wantTitle:  "",
			wantAuthor: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fetcher{
				client: tt.setupMock(),
			}
			content, err := f.FetchIssue(tt.owner, tt.repo, tt.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchIssue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if content.Title != tt.wantTitle {
					t.Errorf("FetchIssue() title = %v, want %v", content.Title, tt.wantTitle)
				}
				if content.Author != tt.wantAuthor {
					t.Errorf("FetchIssue() author = %v, want %v", content.Author, tt.wantAuthor)
				}
			}
		})
	}
}

func TestFetcher_FetchPullRequest(t *testing.T) {
	tests := []struct {
		name       string
		owner      string
		repo       string
		number     int
		setupMock  func() *http.Client
		wantTitle  string
		wantAuthor string
		wantErr    bool
	}{
		{
			name:   "fetch PR successfully",
			owner:  "gorilla",
			repo:   "mux",
			number: 456,
			setupMock: func() *http.Client {
				return &http.Client{
					Transport: &mockTransport{
						roundTripFunc: func(req *http.Request) (*http.Response, error) {
							if req.URL.Path == "/repos/gorilla/mux/pulls/456" {
								return &http.Response{
									StatusCode: http.StatusOK,
									Body:       io.NopCloser(bytes.NewReader([]byte(`{"title":"Test PR","user":{"login":"testuser"},"body":"PR body","state":"open","merged":false}`))),
									Header:     http.Header{"Content-Type": []string{"application/json"}},
								}, nil
							}
							if req.URL.Path == "/repos/gorilla/mux/pulls/456/comments" {
								return &http.Response{
									StatusCode: http.StatusOK,
									Body:       io.NopCloser(bytes.NewReader([]byte(`[{"body":"PR comment1","created_at":"2024-01-01T00:00:00Z","user":{"login":"user1"}}]`))),
									Header:     http.Header{"Content-Type": []string{"application/json"}},
								}, nil
							}
							return &http.Response{StatusCode: http.StatusNotFound}, nil
						},
					},
				}
			},
			wantTitle:  "Test PR",
			wantAuthor: "testuser",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fetcher{
				client: tt.setupMock(),
			}
			content, err := f.FetchPullRequest(tt.owner, tt.repo, tt.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchPullRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if content.Title != tt.wantTitle {
					t.Errorf("FetchPullRequest() title = %v, want %v", content.Title, tt.wantTitle)
				}
				if content.Author != tt.wantAuthor {
					t.Errorf("FetchPullRequest() author = %v, want %v", content.Author, tt.wantAuthor)
				}
			}
		})
	}
}

func TestFetcher_FetchDiscussion(t *testing.T) {
	tests := []struct {
		name       string
		owner      string
		repo       string
		number     int
		setupMock  func() *http.Client
		wantTitle  string
		wantAuthor string
		wantErr    bool
	}{
		{
			name:   "fetch discussion successfully",
			owner:  "gorilla",
			repo:   "mux",
			number: 789,
			setupMock: func() *http.Client {
				return &http.Client{
					Transport: &mockTransport{
						roundTripFunc: func(req *http.Request) (*http.Response, error) {
							// GitHub Discussions API 路径
							if req.URL.Path == "/repos/gorilla/mux/discussions/789" {
								return &http.Response{
									StatusCode: http.StatusOK,
									Body:       io.NopCloser(bytes.NewReader([]byte(`{"title":"Test Discussion","user":{"login":"testuser"},"body":"Discussion body","category":{"name":"General"}}`))),
									Header:     http.Header{"Content-Type": []string{"application/json"}},
								}, nil
							}
							// Discussions comments 路径
							if req.URL.Path == "/repos/gorilla/mux/discussions/789/comments" {
								return &http.Response{
									StatusCode: http.StatusOK,
									Body:       io.NopCloser(bytes.NewReader([]byte(`[{"body":"Discussion comment1","created_at":"2024-01-01T00:00:00Z","user":{"login":"user1"}}]`))),
									Header:     http.Header{"Content-Type": []string{"application/json"}},
								}, nil
							}
							return &http.Response{StatusCode: http.StatusNotFound}, nil
						},
					},
				}
			},
			wantTitle:  "Test Discussion",
			wantAuthor: "testuser",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fetcher{
				client: tt.setupMock(),
			}
			content, err := f.FetchDiscussion(tt.owner, tt.repo, tt.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchDiscussion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if content.Title != tt.wantTitle {
					t.Errorf("FetchDiscussion() title = %v, want %v", content.Title, tt.wantTitle)
				}
				if content.Author != tt.wantAuthor {
					t.Errorf("FetchDiscussion() author = %v, want %v", content.Author, tt.wantAuthor)
				}
			}
		})
	}
}

func TestFetcher_Fetch(t *testing.T) {
	tests := []struct {
		name       string
		owner      string
		repo       string
		issueType  string
		number     int
		setupMock  func() *http.Client
		wantType   string
		wantTitle  string
		wantErr    bool
	}{
		{
			name:      "fetch issue via Fetch",
			owner:     "gorilla",
			repo:      "mux",
			issueType: "issue",
			number:    123,
			setupMock: func() *http.Client {
				return &http.Client{
					Transport: &mockTransport{
						roundTripFunc: func(req *http.Request) (*http.Response, error) {
							if req.URL.Path == "/repos/gorilla/mux/issues/123" {
								return &http.Response{
									StatusCode: http.StatusOK,
									Body:       io.NopCloser(bytes.NewReader([]byte(`{"title":"Test Issue","user":{"login":"testuser"},"body":"Issue body","state":"open"}`))),
									Header:     http.Header{"Content-Type": []string{"application/json"}},
								}, nil
							}
							if req.URL.Path == "/repos/gorilla/mux/issues/123/comments" {
								return &http.Response{
									StatusCode: http.StatusOK,
									Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
									Header:     http.Header{"Content-Type": []string{"application/json"}},
								}, nil
							}
							return &http.Response{StatusCode: http.StatusNotFound}, nil
						},
					},
				}
			},
			wantType:  "issue",
			wantTitle: "Test Issue",
			wantErr:   false,
		},
		{
			name:      "fetch PR via Fetch",
			owner:     "gorilla",
			repo:      "mux",
			issueType: "pr",
			number:    456,
			setupMock: func() *http.Client {
				return &http.Client{
					Transport: &mockTransport{
						roundTripFunc: func(req *http.Request) (*http.Response, error) {
							if req.URL.Path == "/repos/gorilla/mux/pulls/456" {
								return &http.Response{
									StatusCode: http.StatusOK,
									Body:       io.NopCloser(bytes.NewReader([]byte(`{"title":"Test PR","user":{"login":"testuser"},"body":"PR body","state":"open","merged":false}`))),
									Header:     http.Header{"Content-Type": []string{"application/json"}},
								}, nil
							}
							if req.URL.Path == "/repos/gorilla/mux/pulls/456/comments" {
								return &http.Response{
									StatusCode: http.StatusOK,
									Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
									Header:     http.Header{"Content-Type": []string{"application/json"}},
								}, nil
							}
							return &http.Response{StatusCode: http.StatusNotFound}, nil
						},
					},
				}
			},
			wantType:  "pr",
			wantTitle: "Test PR",
			wantErr:   false,
		},
		{
			name:      "fetch discussion via Fetch",
			owner:     "gorilla",
			repo:      "mux",
			issueType: "discussion",
			number:    789,
			setupMock: func() *http.Client {
				return &http.Client{
					Transport: &mockTransport{
						roundTripFunc: func(req *http.Request) (*http.Response, error) {
							if req.URL.Path == "/repos/gorilla/mux/discussions/789" {
								return &http.Response{
									StatusCode: http.StatusOK,
									Body:       io.NopCloser(bytes.NewReader([]byte(`{"title":"Test Discussion","user":{"login":"testuser"},"body":"Discussion body","category":{"name":"General"}}`))),
									Header:     http.Header{"Content-Type": []string{"application/json"}},
								}, nil
							}
							if req.URL.Path == "/repos/gorilla/mux/discussions/789/comments" {
								return &http.Response{
									StatusCode: http.StatusOK,
									Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
									Header:     http.Header{"Content-Type": []string{"application/json"}},
								}, nil
							}
							return &http.Response{StatusCode: http.StatusNotFound}, nil
						},
					},
				}
			},
			wantType:  "discussion",
			wantTitle: "Test Discussion",
			wantErr:   false,
		},
		{
			name:       "invalid issue type",
			owner:      "gorilla",
			repo:       "mux",
			issueType:  "invalid",
			number:     123,
			setupMock:  func() *http.Client { return &http.Client{} },
			wantType:   "",
			wantTitle:  "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fetcher{
				client: tt.setupMock(),
			}
			content, err := f.Fetch(tt.owner, tt.repo, tt.issueType, tt.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if content.Type != tt.wantType {
					t.Errorf("Fetch() type = %v, want %v", content.Type, tt.wantType)
				}
				if content.Title != tt.wantTitle {
					t.Errorf("Fetch() title = %v, want %v", content.Title, tt.wantTitle)
				}
			}
		})
	}
}
