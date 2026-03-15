package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/zhouxinyu1cp/go-claude-test-project/internal/types"
)

// Fetcher GitHub API 客户端
type Fetcher struct {
	client *http.Client
	token  string
}

// NewFetcher 创建 Fetcher 实例
// token: 可选 GitHub Token (环境变量 GITHUB_TOKEN)
func NewFetcher(token string) *Fetcher {
	return &Fetcher{
		client: &http.Client{},
		token:  token,
	}
}

// FetchIssue 获取 Issue 及评论
func (f *Fetcher) FetchIssue(owner, repo string, number int) (*types.GitHubContent, error) {
	ctx := context.Background()

	// 获取 Issue
	issueURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d", owner, repo, number)
	issueResp, err := f.doRequest(ctx, issueURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issue: %w", err)
	}
	defer issueResp.Body.Close()

	if issueResp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("issue not found")
	}

	var issue issueJSON
	if err := json.NewDecoder(issueResp.Body).Decode(&issue); err != nil {
		return nil, fmt.Errorf("failed to decode issue: %w", err)
	}

	// 获取评论
	commentsURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d/comments", owner, repo, number)
	commentsResp, err := f.doRequest(ctx, commentsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}
	defer commentsResp.Body.Close()

	var comments []commentJSON
	if err := json.NewDecoder(commentsResp.Body).Decode(&comments); err != nil {
		return nil, fmt.Errorf("failed to decode comments: %w", err)
	}

	return &types.GitHubContent{
		Type:      "issue",
		Title:     issue.Title,
		Body:      issue.Body,
		Author:    issue.User.Login,
		CreatedAt: issue.CreatedAt.Time,
		State:     issue.State,
		Labels:    extractLabelNames(issue.Labels),
		URL:       issue.HTMLURL,
		Comments:  convertComments(comments),
	}, nil
}

// FetchPullRequest 获取 PR 及评论
func (f *Fetcher) FetchPullRequest(owner, repo string, number int) (*types.GitHubContent, error) {
	ctx := context.Background()

	// 获取 PR
	prURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d", owner, repo, number)
	prResp, err := f.doRequest(ctx, prURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PR: %w", err)
	}
	defer prResp.Body.Close()

	if prResp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("PR not found")
	}

	var pr prJSON
	if err := json.NewDecoder(prResp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("failed to decode PR: %w", err)
	}

	// 获取评论
	commentsURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/comments", owner, repo, number)
	commentsResp, err := f.doRequest(ctx, commentsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}
	defer commentsResp.Body.Close()

	var comments []commentJSON
	if err := json.NewDecoder(commentsResp.Body).Decode(&comments); err != nil {
		return nil, fmt.Errorf("failed to decode comments: %w", err)
	}

	var mergedAt *time.Time
	if pr.MergedAt != nil {
		mergedAt = &pr.MergedAt.Time
	}

	return &types.GitHubContent{
		Type:      "pr",
		Title:     pr.Title,
		Body:      pr.Body,
		Author:    pr.User.Login,
		CreatedAt: pr.CreatedAt.Time,
		State:     pr.State,
		MergedAt:  mergedAt,
		URL:       pr.HTMLURL,
		Comments:  convertComments(comments),
	}, nil
}

// FetchDiscussion 获取 Discussion 及评论
func (f *Fetcher) FetchDiscussion(owner, repo string, number int) (*types.GitHubContent, error) {
	ctx := context.Background()

	// 获取 Discussion
	discussionURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/discussions/%d", owner, repo, number)
	discussionResp, err := f.doRequest(ctx, discussionURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch discussion: %w", err)
	}
	defer discussionResp.Body.Close()

	if discussionResp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("discussion not found")
	}

	var discussion discussionJSON
	if err := json.NewDecoder(discussionResp.Body).Decode(&discussion); err != nil {
		return nil, fmt.Errorf("failed to decode discussion: %w", err)
	}

	// 获取评论
	commentsURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/discussions/%d/comments", owner, repo, number)
	commentsResp, err := f.doRequest(ctx, commentsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}
	defer commentsResp.Body.Close()

	var comments []discussionCommentJSON
	if err := json.NewDecoder(commentsResp.Body).Decode(&comments); err != nil {
		return nil, fmt.Errorf("failed to decode comments: %w", err)
	}

	category := ""
	if discussion.Category != nil {
		category = discussion.Category.Name
	}

	return &types.GitHubContent{
		Type:      "discussion",
		Title:     discussion.Title,
		Body:      discussion.Body,
		Author:    discussion.User.Login,
		CreatedAt: discussion.CreatedAt.Time,
		State:     "open",
		Category:  category,
		URL:       discussion.HTMLURL,
		Comments:  convertDiscussionComments(comments),
	}, nil
}

// Fetch 根据 issueType 自动分发
func (f *Fetcher) Fetch(owner, repo, issueType string, number int) (*types.GitHubContent, error) {
	switch issueType {
	case "issue":
		return f.FetchIssue(owner, repo, number)
	case "pr":
		return f.FetchPullRequest(owner, repo, number)
	case "discussion":
		return f.FetchDiscussion(owner, repo, number)
	default:
		return nil, fmt.Errorf("invalid issue type: %s", issueType)
	}
}

// doRequest 执行 HTTP 请求
func (f *Fetcher) doRequest(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	if f.token != "" {
		req.Header.Set("Authorization", "token "+f.token)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// JSON 辅助结构

type issueJSON struct {
	Title     string         `json:"title"`
	Body      string         `json:"body"`
	User      userJSON       `json:"user"`
	State     string         `json:"state"`
	Labels    []labelJSON    `json:"labels"`
	CreatedAt githubTime     `json:"created_at"`
	HTMLURL   string         `json:"html_url"`
}

type prJSON struct {
	Title     string         `json:"title"`
	Body      string         `json:"body"`
	User      userJSON       `json:"user"`
	State     string         `json:"state"`
	MergedAt  *githubTime    `json:"merged_at"`
	CreatedAt githubTime     `json:"created_at"`
	HTMLURL   string         `json:"html_url"`
}

type discussionJSON struct {
	Title     string                `json:"title"`
	Body      string                `json:"body"`
	User      userJSON              `json:"user"`
	Category  *discussionCategory  `json:"category"`
	CreatedAt githubTime            `json:"created_at"`
	HTMLURL   string                `json:"html_url"`
}

type userJSON struct {
	Login string `json:"login"`
}

type labelJSON struct {
	Name string `json:"name"`
}

type discussionCategory struct {
	Name string `json:"name"`
}

type commentJSON struct {
	Body      string     `json:"body"`
	CreatedAt githubTime `json:"created_at"`
	User      userJSON   `json:"user"`
}

type discussionCommentJSON struct {
	Body      string     `json:"body"`
	CreatedAt githubTime `json:"created_at"`
	Author    userJSON   `json:"author"`
}

type githubTime struct {
	Time time.Time
}

func (gt *githubTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("failed to unmarshal time: %w", err)
	}
	t, err := time.Parse("2006-01-02T15:04:05Z07:00", s)
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}
	gt.Time = t
	return nil
}

func extractLabelNames(labels []labelJSON) []string {
	result := make([]string, len(labels))
	for i, label := range labels {
		result[i] = label.Name
	}
	return result
}

func convertComments(comments []commentJSON) []types.GitHubComment {
	result := make([]types.GitHubComment, len(comments))
	for i, c := range comments {
		result[i] = types.GitHubComment{
			Body:      c.Body,
			CreatedAt: c.CreatedAt.Time,
			User:      types.GitHubUser{Login: c.User.Login},
		}
	}
	return result
}

func convertDiscussionComments(comments []discussionCommentJSON) []types.GitHubComment {
	result := make([]types.GitHubComment, len(comments))
	for i, c := range comments {
		result[i] = types.GitHubComment{
			Body:      c.Body,
			CreatedAt: c.CreatedAt.Time,
			User:      types.GitHubUser{Login: c.Author.Login},
		}
	}
	return result
}

