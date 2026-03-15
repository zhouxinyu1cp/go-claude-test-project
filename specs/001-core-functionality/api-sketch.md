# API Sketch

本文档描述 `internal` 目录下各包的对外接口，作为开发参考。

---

## 1. internal/types

定义跨包共享的数据结构。

### 类型定义

```go
// GitHubIssue GitHub Issue 数据结构
type GitHubIssue struct {
    Number    int
    Title     string
    Body      string
    State     string // "open" | "closed"
    CreatedAt time.Time
    UpdatedAt time.Time
    HTMLURL   string
    User      GitHubUser
    Labels    []string
}

// GitHubPullRequest GitHub PR 数据结构
type GitHubPullRequest struct {
    Number    int
    Title     string
    Body      string
    State     string // "open" | "closed" | "merged"
    CreatedAt time.Time
    UpdatedAt time.Time
    MergedAt  *time.Time // 可能为 nil
    HTMLURL   string
    User      GitHubUser
}

// GitHubDiscussion GitHub Discussion 数据结构
type GitHubDiscussion struct {
    Number    int
    Title     string
    Body      string
    State     string // "open" | "closed"
    CreatedAt time.Time
    UpdatedAt time.Time
    HTMLURL   string
    Category  string
    User      GitHubUser
}

// GitHubUser GitHub 用户
type GitHubUser struct {
    Login string
}

// GitHubComment 评论结构（Issue Comment、PR Review Comment、Discussion Comment 通用）
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

## 2. internal/parser

负责解析 GitHub URL，提取 owner、repo、类型和编号。

### 函数签名

```go
// ParseURL 解析 GitHub URL，返回 owner、repo、类型和编号
// 支持的 URL 类型：
//   - https://github.com/owner/repo/issues/123
//   - https://github.com/owner/repo/pull/456
//   - https://github.com/owner/repo/discussions/789
//
// 返回值：
//   - owner: 仓库所有者
//   - repo: 仓库名
//   - issueType: "issue" | "pr" | "discussion"
//   - number: Issue/PR/Discussion 编号
//   - err: 解析错误
func ParseURL(rawURL string) (owner, repo, issueType string, number int, err error)

// IsValidURL 简单验证 URL 格式是否有效
func IsValidURL(rawURL string) bool
```

---

## 3. internal/fetcher

负责与 GitHub API 交互，获取 Issue/PR/Discussion 数据及评论。

### 函数签名

```go
// NewFetcher 创建 Fetcher 实例
// token: 可选的 GitHub Personal Access Token（用于私有仓库或提高 rate limit）
func NewFetcher(token string) *Fetcher

// FetchIssue 获取 Issue 及其评论
func (f *Fetcher) FetchIssue(owner, repo string, number int) (*types.GitHubContent, error)

// FetchPullRequest 获取 PR 及其评论
func (f *Fetcher) FetchPullRequest(owner, repo string, number int) (*types.GitHubContent, error)

// FetchDiscussion 获取 Discussion 及其评论
func (f *Fetcher) FetchDiscussion(owner, repo string, number int) (*types.GitHubContent, error)

// Fetch 根据类型自动获取 Issue/PR/Discussion
func (f *Fetcher) Fetch(owner, repo, issueType string, number int) (*types.GitHubContent, error)
```

---

## 4. internal/formatter

负责文件名生成、字符串清理等辅助功能。

### 函数签名

```go
// SanitizeFilename 清理字符串，生成合法的文件名
// 将非法字符（/、:、空格等）替换为 -
func SanitizeFilename(s string) string

// GenerateFilename 生成输出文件名
// 参数：
//   - owner: 仓库所有者
//   - repo: 仓库名
//   - issueType: "issue" | "pr" | "discussion"
//   - number: 编号
//   - title: 标题（会经过 SanitizeFilename 处理）
//
// 返回值例如：
//   - "gorilla-mux-issue-123-refactor-handler.md"
//   - "gorilla-mux-pr-456-add-middleware.md"
//   - "gorilla-mux-discussion-789-roadmap-2024.md"
func GenerateFilename(owner, repo, issueType string, number int, title string) string

// ResolveFilenameConflict 解决文件名冲突
// 如果文件已存在，在文件名末尾添加 _1, _2 等序号
func ResolveFilenameConflict(dir, filename string) string
```

---

## 5. internal/converter

负责将 `types.GitHubContent` 转换为 Markdown 格式字符串。

### 函数签名

```go
// Converter 评论排序方向配置
type Converter struct {
    Order string // "asc" | "desc"
}

// NewConverter 创建 Converter 实例
func NewConverter(order string) *Converter

// Convert 将 GitHubContent 转换为 Markdown 字符串
func (c *Converter) Convert(content *types.GitHubContent) string

// FormatComment 格式化单条评论
func (c *Converter) FormatComment(comment types.GitHubComment) string
```

---

## 6. cmd/issue2md

CLI 入口，负责命令行参数解析、流程编排。

### 主流程

```
1. 解析命令行参数（-o, --order）
2. 读取 URL（单个或批量）
3. 对每个 URL：
   a. parser.ParseURL() 解析 URL
   b. fetcher.Fetch() 获取数据
   c. converter.Convert() 转换为 Markdown
   d. formatter.GenerateFilename() 生成文件名
   e. formatter.ResolveFilenameConflict() 处理冲突
   f. 写入文件
4. 打印成功信息
```

---

## 7. 包依赖关系

```
cmd/issue2md
    ├── parser
    ├── fetcher
    ├── converter
    ├── formatter
    └── types (被所有包引用)

fetcher ──► types
parser ──► types
converter ──► types
formatter ──► (无外部依赖)
```

---

## 8. 错误处理约定

所有返回 error 的函数遵循以下约定：
- 使用 `fmt.Errorf("...: %w", err)` 包装底层错误
- 错误信息清晰描述失败原因（如 "failed to fetch issue: 404 Not Found"）
- 不使用全局错误变量
