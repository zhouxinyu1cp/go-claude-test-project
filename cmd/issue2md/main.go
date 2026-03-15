package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zhouxinyu1cp/go-claude-test-project/internal/converter"
	"github.com/zhouxinyu1cp/go-claude-test-project/internal/fetcher"
	"github.com/zhouxinyu1cp/go-claude-test-project/internal/formatter"
	"github.com/zhouxinyu1cp/go-claude-test-project/internal/parser"
)

func main() {
	// 定义命令行 flags
	outputDir := flag.String("o", "./", "Output directory")
	order := flag.String("order", "desc", "Comment order (asc/desc)")
	flag.Parse()

	// 获取 URL 参数
	urls := flag.Args()
	if len(urls) == 0 {
		fmt.Fprintln(os.Stderr, "Error: no URLs provided")
		fmt.Fprintln(os.Stderr, "Usage: issue2md [-o <dir>] [--order <asc|desc>] <url> [<url>...]")
		os.Exit(1)
	}

	// 验证 order 参数
	if *order != "asc" && *order != "desc" {
		fmt.Fprintln(os.Stderr, "Error: order must be 'asc' or 'desc'")
		os.Exit(1)
	}

	// 创建输出目录
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create output directory: %v\n", err)
		os.Exit(1)
	}

	// 获取 GitHub Token
	token := os.Getenv("GITHUB_TOKEN")

	// 初始化组件
	f := fetcher.NewFetcher(token)
	c := converter.NewConverter(*order)

	// 记录成功和失败的数量
	successCount := 0
	failCount := 0

	// 批量处理 URL
	for _, url := range urls {
		if err := processURL(url, *outputDir, f, c); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			failCount++
			continue
		}
		successCount++
	}

	// 输出统计
	fmt.Printf("Completed: %d succeeded, %d failed\n", successCount, failCount)

	if failCount > 0 {
		os.Exit(1)
	}
}

// processURL 处理单个 URL
func processURL(url, outputDir string, f *fetcher.Fetcher, c *converter.Converter) error {
	// 解析 URL
	owner, repo, issueType, number, err := parser.ParseURL(url)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	// 获取 GitHub 内容
	content, err := f.Fetch(owner, repo, issueType, number)
	if err != nil {
		return fmt.Errorf("failed to fetch content: %w", err)
	}

	// 生成文件名
	filename := formatter.GenerateFilename(owner, repo, issueType, number, content.Title)
	filename = formatter.ResolveFilenameConflict(outputDir, filename)

	// 转换为 Markdown
	markdown := c.Convert(content)

	// 写入文件
	fpath := filepath.Join(outputDir, filename)
	if err := os.WriteFile(fpath, []byte(markdown), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Generated: %s\n", filename)
	return nil
}
