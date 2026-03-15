# issue2md 技术方案

## 1. 技术上下文总结

### 1.1 技术选型

| 组件 | 技术选型 | 说明 |
|------|----------|------|
| **语言** | Go >= 1.21.0 | 升级现有 go.mod 中的版本 |
| **Web 框架** | `net/http` (标准库) | 不引入 Gin/Echo 等外部框架 |
| **GitHub API** | `google/go-github` v67+ | 使用 REST API（spec 未要求 GraphQL） |
| **JSON 解析** | `encoding/json` (标准库) | 用于 API 响应解析 |
| **Markdown 输出** | 标准库 strings.Builder | 手动拼接 Markdown，不引入第三方库 |
| **数据库** | 无 | 所有数据通过 API 实时获取 |

### 1.2 依赖评估

根据用户要求使用 `google/go-github` 库：

```go
// go.mod 预计依赖
require (
    github.com/google/go-github/v67 v67.0.0
)
```

---

## 2. "合宪性"审查

### 2.1 简单性原则 (Simplicity First)

| 宪法条款 | 符合情况 | 说明 |
|---------|---------|------|
| 1.1 YAGNI | ✅ | 仅实现 spec.md 中明确要求的功能 |
| 1.2 标准库优先 | ✅ | HTTP 请求使用 go-github（用户指定），Markdown 生成使用标准库 |
| 1.3 反过度工程 | ✅ | 简单函数/结构体优于复杂接口 |

### 2.2 测试先行铁律 (Test-First Imperative)

| 宪法条款 | 符合情况 | 说明 |
|---------|---------|------|
| 2.1 TDD 循环 | ✅ | 每个功能开发遵循 Red-Green-Refactor |
| 2.2 表格驱动 | ✅ | 单元测试采用表格驱动风格 |
| 2.3 拒绝 Mocks | ✅ | 优先使用真实 HTTP 请求，可设置 GitHub API mock 用于测试 |

### 2.3 明确性原则 (Clarity and Explicitness)

| 宪法条款 | 符合情况 | 说明 |
|---------|---------|------|
| 3.1 错误处理 | ✅ | 所有错误使用 `fmt.Errorf("...: %w", err)` 包装 |
| 3.2 无全局变量 | ✅ | 依赖通过结构体成员或函数参数注入 |

---

## 3. 项目结构细化

### 3.1 目录结构

```
issue2md/
├── cmd/
│   └── issue2md/
│       └── main.go           # CLI 入口，参数解析与流程编排
├── internal/
│   ├── types/
│       └── types.go           # 核心数据结构定义
│   ├── parser/
│       └── parser.go         # URL 解析逻辑
│   ├── fetcher/
│       └── fetcher.go        # GitHub API 请求（使用 go-github）
│   ├── converter/
│       └── converter.go      # GitHubContent -> Markdown
│   └── formatter/
│       └── formatter.go      # 文件名生成、冲突处理
├── go.mod
├── go.sum
├── Makefile
└── spec.md
```

### 3.2 各包职责与依赖

| 包 | 职责 | 依赖 |
|---|------|------|
| `cmd/issue2md` | CLI 入口、flag 解析、批量处理流程 | parser, fetcher, converter, formatter, types |
| `internal/types` | 定义跨包共享的数据结构 | 无 |
| `internal/parser` | 解析 GitHub URL | types |
| `internal/fetcher` | 调用 GitHub API 获取数据 | types, github.com/google/go-github/v67 |
| `internal/converter` | 将 GitHubContent 转为 Markdown | types |
| `internal/formatter` | 文件名生成与冲突解决 | 无 |

---

## 4. 核心数据结构

### 4.1 types.go

```go
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
    Type     string // "issue" | "pr" | "discussion"
    Title    string
    Body     string
    Author   string
    CreatedAt time.Time
    State    string
    Labels   []string       // 仅 Issue 有
    Category string         // 仅 Discussion 有
    MergedAt *time.Time     // 仅 PR 有
    URL      string
    Comments []GitHubComment
}
```

---

## 5. 接口设计

### 5.1 internal/parser

```go
// ParseURL 解析 GitHub URL
// 返回: owner, repo, issueType ("issue"|"pr"|"discussion"), number, error
func ParseURL(rawURL string) (owner, repo, issueType string, number int, err error)

// IsValidURL 验证 URL 格式是否有效
func IsValidURL(rawURL string) bool
```

### 5.2 internal/fetcher

```go
// Fetcher GitHub API 客户端
type Fetcher struct {
    client  *github.Client
    token   string
}

// NewFetcher 创建 Fetcher 实例
// token: 可选 GitHub Token (环境变量 GITHUB_TOKEN)
func NewFetcher(token string) *Fetcher

// FetchIssue 获取 Issue 及评论
func (f *Fetcher) FetchIssue(owner, repo string, number int) (*types.GitHubContent, error)

// FetchPullRequest 获取 PR 及评论
func (f *Fetcher) FetchPullRequest(owner, repo string, number int) (*types.GitHubContent, error)

// FetchDiscussion 获取 Discussion 及评论
func (f *Fetcher) FetchDiscussion(owner, repo string, number int) (*types.GitHubContent, error)

// Fetch 根据 issueType 自动分发
func (f *Fetcher) Fetch(owner, repo, issueType string, number int) (*types.GitHubContent, error)
```

### 5.3 internal/formatter

```go
// SanitizeFilename 清理非法字符，生成合法文件名
// 替换: /, :, 空格 -> -
func SanitizeFilename(s string) string

// GenerateFilename 生成输出文件名
// 格式:
//   - Issue: {owner}-{repo}-{number}-{title_slug}.md
//   - PR: {owner}-{repo}-pr-{number}-{title_slug}.md
//   - Discussion: {owner}-{repo}-discussion-{number}-{title_slug}.md
func GenerateFilename(owner, repo, issueType string, number int, title string) string

// ResolveFilenameConflict 解决文件名冲突
// 已存在则添加 _1, _2 后缀
func ResolveFilenameConflict(dir, filename string) string
```

### 5.4 internal/converter

```go
// Converter Markdown 转换器
type Converter struct {
    Order string // "asc" | "desc"
}

// NewConverter 创建 Converter 实例
func NewConverter(order string) *Converter

// Convert 将 GitHubContent 转换为 Markdown 字符串
func (c *Converter) Convert(content *types.GitHubContent) string
```

---

## 6. CLI 参数设计

### 6.1 命令行 Flags

| Flag | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| `--output` | `-o` | `./` | 输出目录 |
| `--order` | 无 | `desc` | 评论排序方向 (`asc`/`desc`) |

### 6.2 使用方式

```bash
# 单个 URL
./issue2md https://github.com/gorilla/mux/issues/123

# 批量 URL (空格分隔)
./issue2md https://github.com/gorilla/mux/issues/123 https://github.com/gorilla/mux/pull/456

# 指定输出目录
./issue2md -o /tmp https://github.com/gorilla/mux/issues/123

# 指定评论排序
./issue2md --order asc https://github.com/gorilla/mux/issues/123
```

---

## 7. 错误处理规范

| 错误场景 | 处理方式 |
|---------|---------|
| URL 解析失败 | `fmt.Errorf("invalid URL: %w", err)` |
| GitHub API 404 | `fmt.Errorf("resource not found: %w", err)` |
| Rate Limit | 透传 GitHub 错误信息（含 Retry-After） |
| 网络超时 | 透传错误信息 |
| 文件写入失败 | `fmt.Errorf("failed to write file: %w", err)` |
| 批量处理时单个失败 | 打印错误到 stderr，继续处理下一个 |

---

## 8. 验证方案

### 8.1 单元测试

```bash
# 运行所有测试
make test

# 测试覆盖目标
- parser: URL 解析各种边界情况
- formatter: 文件名生成、冲突处理
- converter: Markdown 格式输出
```

### 8.2 手动验证

| 测试 Case | 验证方式 |
|-----------|---------|
| TC001 单个 Issue URL | 执行 CLI，检查输出文件内容 |
| TC002 单个 PR URL | 执行 CLI，检查输出文件内容 |
| TC003 单个 Discussion URL | 执行 CLI，检查输出文件内容 |
| TC004 批量 3 个 URL | 执行 CLI，检查生成 3 个文件 |
| TC005 指定输出目录 | 使用 `-o /tmp`，检查文件位置 |
| TC006 --order asc | 检查评论顺序 |
| TC007 默认 order | 检查评论倒序 |
| TC008 批量错误处理 | 故意输入错误 URL，检查错误打印 |
| TC009 文件名特殊字符 | Issue 标题含 `/` 和空格，检查文件名 |
| TC010 文件名冲突 | 同一 URL 执行两次，检查 `_1` 后缀 |
| TC011 404 错误 | 输入不存在 Issue，检查错误信息 |
| TC012 公开仓库无 Token | 不设置 GITHUB_TOKEN，检查正常获取 |
| TC013 私有仓库有 Token | 设置 GITHUB_TOKEN，检查正常获取 |

---

## 9. 下一步

技术方案评审通过后，进入实现阶段：

1. 更新 `go.mod` Go 版本至 1.21.0
2. 添加 `google/go-github/v67` 依赖
3. 按照 TDD 循环逐步实现各包
