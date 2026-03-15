package converter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/zhouxinyu1cp/go-claude-test-project/internal/types"
)

// Converter Markdown 转换器
type Converter struct {
	Order string // "asc" | "desc"
}

// NewConverter 创建 Converter 实例
func NewConverter(order string) *Converter {
	return &Converter{Order: order}
}

// Convert 将 GitHubContent 转换为 Markdown 字符串
func (c *Converter) Convert(content *types.GitHubContent) string {
	var sb strings.Builder

	// 标题
	sb.WriteString("# ")
	sb.WriteString(content.Title)
	sb.WriteString("\n\n")

	// 元信息
	sb.WriteString("**Author:** ")
	sb.WriteString(content.Author)
	sb.WriteString("\n")

	sb.WriteString("**State:** ")
	sb.WriteString(content.State)
	sb.WriteString("\n")

	sb.WriteString("**Created:** ")
	sb.WriteString(content.CreatedAt.Format("2006-01-02 15:04:05"))
	sb.WriteString("\n")

	// Issue 特有字段
	if len(content.Labels) > 0 {
		sb.WriteString("**Labels:** ")
		sb.WriteString(strings.Join(content.Labels, ", "))
		sb.WriteString("\n")
	}

	// Discussion 特有字段
	if content.Category != "" {
		sb.WriteString("**Category:** ")
		sb.WriteString(content.Category)
		sb.WriteString("\n")
	}

	// PR 特有字段
	if content.MergedAt != nil && !content.MergedAt.IsZero() {
		sb.WriteString("**Merged:** ")
		sb.WriteString(content.MergedAt.Format("2006-01-02 15:04:05"))
		sb.WriteString("\n")
	}

	sb.WriteString("\n**URL:** ")
	sb.WriteString(content.URL)
	sb.WriteString("\n\n")

	// 正文
	sb.WriteString("---\n\n")
	sb.WriteString(content.Body)
	sb.WriteString("\n\n")

	// 评论
	if len(content.Comments) > 0 {
		sb.WriteString("## Comments\n\n")

		// 排序
		comments := make([]types.GitHubComment, len(content.Comments))
		copy(comments, content.Comments)

		if c.Order == "asc" {
			sort.Slice(comments, func(i, j int) bool {
				return comments[i].CreatedAt.Before(comments[j].CreatedAt)
			})
		} else {
			// 默认 desc
			sort.Slice(comments, func(i, j int) bool {
				return comments[i].CreatedAt.After(comments[j].CreatedAt)
			})
		}

		for _, comment := range comments {
			sb.WriteString(fmt.Sprintf("**User:** %s | **Date:** %s\n\n",
				comment.User.Login,
				comment.CreatedAt.Format("2006-01-02 15:04:05")))
			sb.WriteString(comment.Body)
			sb.WriteString("\n\n---\n\n")
		}
	}

	return sb.String()
}
