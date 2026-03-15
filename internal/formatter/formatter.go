package formatter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SanitizeFilename 清理非法字符，生成合法文件名
// 替换: /, :, 空格 -> -
func SanitizeFilename(s string) string {
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, ":", "-")
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

// GenerateFilename 生成输出文件名
// 格式:
//   - Issue: {owner}-{repo}-{number}-{title_slug}.md
//   - PR: {owner}-{repo}-pr-{number}-{title_slug}.md
//   - Discussion: {owner}-{repo}-discussion-{number}-{title_slug}.md
func GenerateFilename(owner, repo, issueType string, number int, title string) string {
	// 清理标题中的非法字符
	titleSlug := SanitizeFilename(title)
	// 转小写
	titleSlug = strings.ToLower(titleSlug)

	var prefix string
	switch issueType {
	case "issue":
		prefix = fmt.Sprintf("%s-%s-%d-", owner, repo, number)
	case "pr":
		prefix = fmt.Sprintf("%s-%s-pr-%d-", owner, repo, number)
	case "discussion":
		prefix = fmt.Sprintf("%s-%s-discussion-%d-", owner, repo, number)
	}

	return prefix + titleSlug + ".md"
}

// ResolveFilenameConflict 解决文件名冲突
// 已存在则添加 _1, _2 后缀
func ResolveFilenameConflict(dir, filename string) string {
	fpath := filepath.Join(dir, filename)

	// 文件不存在，直接返回
	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		return filename
	}

	// 分离文件名和扩展名
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)

	// 直接遍历检查每个编号是否存在
	for i := 1; ; i++ {
		newName := fmt.Sprintf("%s_%d%s", base, i, ext)
		newPath := filepath.Join(dir, newName)
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newName
		}
	}
}
