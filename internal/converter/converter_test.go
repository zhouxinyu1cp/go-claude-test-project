package converter

import (
	"strings"
	"testing"
	"time"

	"github.com/zhouxinyu1cp/go-claude-test-project/internal/types"
)

func TestConverter_Convert(t *testing.T) {
	createdAt := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name        string
		content     *types.GitHubContent
		order       string
		checkFunc   func(t *testing.T, output string)
	}{
		{
			name: "convert issue",
			content: &types.GitHubContent{
				Type:      "issue",
				Title:     "Test Issue",
				Body:      "This is the issue body",
				Author:    "testuser",
				CreatedAt: createdAt,
				State:     "open",
				Labels:    []string{"bug", "priority"},
				URL:       "https://github.com/owner/repo/issues/123",
			},
			order: "desc",
			checkFunc: func(t *testing.T, output string) {
				if !strings.Contains(output, "# Test Issue") {
					t.Error("output should contain title")
				}
				if !strings.Contains(output, "**Author:** testuser") {
					t.Error("output should contain author")
				}
				if !strings.Contains(output, "**State:** open") {
					t.Error("output should contain state")
				}
				if !strings.Contains(output, "**Labels:** bug, priority") {
					t.Error("output should contain labels")
				}
				if !strings.Contains(output, "This is the issue body") {
					t.Error("output should contain body")
				}
			},
		},
		{
			name: "convert PR",
			content: &types.GitHubContent{
				Type:      "pr",
				Title:     "Fix Bug",
				Body:      "PR description",
				Author:    "prauthor",
				CreatedAt: createdAt,
				State:     "closed",
				URL:       "https://github.com/owner/repo/pull/456",
			},
			order: "desc",
			checkFunc: func(t *testing.T, output string) {
				if !strings.Contains(output, "# Fix Bug") {
					t.Error("output should contain title")
				}
				if !strings.Contains(output, "**Author:** prauthor") {
					t.Error("output should contain author")
				}
				if !strings.Contains(output, "**State:** closed") {
					t.Error("output should contain state")
				}
			},
		},
		{
			name: "convert discussion",
			content: &types.GitHubContent{
				Type:      "discussion",
				Title:     "Question about usage",
				Body:      "How do I use this?",
				Author:    "questionuser",
				CreatedAt: createdAt,
				State:     "open",
				Category:  "General",
				URL:       "https://github.com/owner/repo/discussions/789",
			},
			order: "desc",
			checkFunc: func(t *testing.T, output string) {
				if !strings.Contains(output, "# Question about usage") {
					t.Error("output should contain title")
				}
				if !strings.Contains(output, "**Author:** questionuser") {
					t.Error("output should contain author")
				}
				if !strings.Contains(output, "**Category:** General") {
					t.Error("output should contain category")
				}
			},
		},
		{
			name: "issue with comments",
			content: &types.GitHubContent{
				Type:      "issue",
				Title:     "Issue with Comments",
				Body:      "Issue body",
				Author:    "author",
				CreatedAt: createdAt,
				State:     "open",
				Comments: []types.GitHubComment{
					{
						Body:      "First comment",
						CreatedAt: time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
						User:      types.GitHubUser{Login: "user1"},
					},
					{
						Body:      "Second comment",
						CreatedAt: time.Date(2024, 1, 17, 10, 0, 0, 0, time.UTC),
						User:      types.GitHubUser{Login: "user2"},
					},
				},
			},
			order: "desc",
			checkFunc: func(t *testing.T, output string) {
				if !strings.Contains(output, "## Comments") {
					t.Error("output should contain comments section")
				}
				if !strings.Contains(output, "First comment") {
					t.Error("output should contain first comment")
				}
				if !strings.Contains(output, "Second comment") {
					t.Error("output should contain second comment")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConverter(tt.order)
			output := c.Convert(tt.content)
			tt.checkFunc(t, output)
		})
	}
}

func TestConverter_CommentSorting(t *testing.T) {
	createdAt := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name      string
		order     string
		wantFirst string // 第一个评论的用户
		wantLast  string // 最后一个评论的用户
	}{
		{
			name:      "desc order",
			order:     "desc",
			wantFirst: "user3", // 最新的在最前
			wantLast:  "user1",
		},
		{
			name:      "asc order",
			order:     "asc",
			wantFirst: "user1", // 最早的在前
			wantLast:  "user3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := &types.GitHubContent{
				Type:      "issue",
				Title:     "Test",
				Body:      "Body",
				Author:    "author",
				CreatedAt: createdAt,
				State:     "open",
				Comments: []types.GitHubComment{
					{
						Body:      "comment1",
						CreatedAt: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
						User:      types.GitHubUser{Login: "user1"},
					},
					{
						Body:      "comment2",
						CreatedAt: time.Date(2024, 1, 16, 12, 0, 0, 0, time.UTC),
						User:      types.GitHubUser{Login: "user2"},
					},
					{
						Body:      "comment3",
						CreatedAt: time.Date(2024, 1, 17, 12, 0, 0, 0, time.UTC),
						User:      types.GitHubUser{Login: "user3"},
					},
				},
			}

			c := NewConverter(tt.order)
			output := c.Convert(content)

			lines := strings.Split(output, "\n")
			var commentLines []string
			for _, line := range lines {
				if strings.Contains(line, "**User:**") {
					commentLines = append(commentLines, line)
				}
			}

			if len(commentLines) < 2 {
				t.Fatalf("expected at least 2 comment lines, got %d", len(commentLines))
			}

			firstLine := commentLines[0]
			lastLine := commentLines[len(commentLines)-1]

			if tt.order == "desc" {
				// 最新在前
				if !strings.Contains(firstLine, tt.wantFirst) {
					t.Errorf("first comment should be %s, got: %s", tt.wantFirst, firstLine)
				}
				if !strings.Contains(lastLine, tt.wantLast) {
					t.Errorf("last comment should be %s, got: %s", tt.wantLast, lastLine)
				}
			} else {
				// 最早在前
				if !strings.Contains(firstLine, tt.wantFirst) {
					t.Errorf("first comment should be %s, got: %s", tt.wantFirst, firstLine)
				}
				if !strings.Contains(lastLine, tt.wantLast) {
					t.Errorf("last comment should be %s, got: %s", tt.wantLast, lastLine)
				}
			}
		})
	}
}
