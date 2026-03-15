package types

import "time"

// GitHubUser GitHub 用户
type GitHubUser struct {
	Login string
}

// GitHubComment 评论结构（Issue/PR/Discussion 通用）
type GitHubComment struct {
	Body      string
	CreatedAt time.Time
	User      GitHubUser
}

// GitHubContent 最终输出的内容结构
type GitHubContent struct {
	Type      string           // "issue" | "pr" | "discussion"
	Title     string
	Body      string
	Author    string
	CreatedAt time.Time
	State     string
	Labels    []string         // 仅 Issue 有
	Category  string           // 仅 Discussion 有
	MergedAt  *time.Time       // 仅 PR 有
	URL       string
	Comments  []GitHubComment
}
