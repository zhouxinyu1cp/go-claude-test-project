package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseURL 解析 GitHub URL
// 返回: owner, repo, issueType ("issue"|"pr"|"discussion"), number, error
func ParseURL(rawURL string) (owner, repo, issueType string, number int, err error) {
	if rawURL == "" {
		return "", "", "", 0, fmt.Errorf("empty URL")
	}

	// 移除协议前缀
	urlWithoutProto := strings.TrimPrefix(rawURL, "https://")
	urlWithoutProto = strings.TrimPrefix(urlWithoutProto, "http://")

	// 按 / 分割
	parts := strings.Split(urlWithoutProto, "/")
	if len(parts) < 5 {
		return "", "", "", 0, fmt.Errorf("invalid URL format")
	}

	// parts[0] 是域名 (github.com 或 github.mycompany.com)
	// parts[1] 是 owner
	// parts[2] 是 repo
	// parts[3] 是类型 (issues, pull, discussions)
	// parts[4] 是编号

	host := parts[0]
	// 简单验证：必须是 github.com 或包含 github 的域名
	if !strings.Contains(host, "github") {
		return "", "", "", 0, fmt.Errorf("not a GitHub URL")
	}

	owner = parts[1]
	repo = parts[2]
	typePart := parts[3]
	numberStr := parts[4]

	// 验证类型
	switch typePart {
	case "issues":
		issueType = "issue"
	case "pull":
		issueType = "pr"
	case "discussions":
		issueType = "discussion"
	default:
		return "", "", "", 0, fmt.Errorf("invalid type: %s", typePart)
	}

	// 解析编号
	number, err = strconv.Atoi(numberStr)
	if err != nil {
		return "", "", "", 0, fmt.Errorf("invalid number: %w", err)
	}

	return owner, repo, issueType, number, nil
}

// IsValidURL 验证 URL 格式是否有效
func IsValidURL(rawURL string) bool {
	_, _, _, _, err := ParseURL(rawURL)
	return err == nil
}
