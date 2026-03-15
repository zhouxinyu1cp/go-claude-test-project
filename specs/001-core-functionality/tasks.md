# issue2md 任务列表

## Phase 1: Foundation (数据结构定义)

### 1.1 项目初始化
- [ ] **Task 1.1.1 [P]** 更新 go.mod Go 版本至 1.21.0
- [ ] **Task 1.1.2 [P]** 添加 google/go-github/v67 依赖 (执行 `go get github.com/google/go-github/v67@v67.0.0`)

### 1.2 核心数据结构
- [ ] **Task 1.2.1** 创建 `internal/types/types.go` - 定义核心数据结构
  - 创建 `GitHubUser` 结构体 (Login 字段)
  - 创建 `GitHubComment` 结构体 (Body, CreatedAt, User 字段)
  - 创建 `GitHubContent` 结构体 (Type, Title, Body, Author, CreatedAt, State, Labels, Category, MergedAt, URL, Comments 字段)

---

## Phase 2: GitHub Fetcher (API交互逻辑，TDD)

### 2.1 URL 解析 (parser)
- [ ] **Task 2.1.1** 创建 `internal/parser/parser_test.go` - 表格驱动测试
  - 测试 ParseURL 解析 Issue URL
  - 测试 ParseURL 解析 PR URL
  - 测试 ParseURL 解析 Discussion URL
  - 测试 ParseURL 无效 URL
  - 测试 IsValidURL 各种 URL 格式
- [ ] **Task 2.1.2** 创建 `internal/parser/parser.go` - 实现 URL 解析
  - 实现 `ParseURL(rawURL string)` 函数
  - 实现 `IsValidURL(rawURL string)` 函数

### 2.2 GitHub API 获取 (fetcher)
- [ ] **Task 2.2.1** 创建 `internal/fetcher/fetcher_test.go` - 表格驱动测试
  - 测试 FetchIssue 获取 Issue 及评论
  - 测试 FetchPullRequest 获取 PR 及评论
  - 测试 FetchDiscussion 获取 Discussion 及评论
  - 测试 Fetch 自动分发逻辑
  - 测试 404 错误处理
- [ ] **Task 2.2.2** 创建 `internal/fetcher/fetcher.go` - 实现 API 请求
  - 创建 `Fetcher` 结构体 (client, token 字段)
  - 实现 `NewFetcher(token string)` 构造函数
  - 实现 `FetchIssue(owner, repo string, number int)` 方法
  - 实现 `FetchPullRequest(owner, repo string, number int)` 方法
  - 实现 `FetchDiscussion(owner, repo string, number int)` 方法
  - 实现 `Fetch(owner, repo, issueType string, number int)` 自动分发方法

---

## Phase 3: Markdown Converter (转换逻辑，TDD)

### 3.1 文件名格式化 (formatter)
- [ ] **Task 3.1.1** 创建 `internal/formatter/formatter_test.go` - 表格驱动测试
  - 测试 SanitizeFilename 清理非法字符
  - 测试 GenerateFilename 生成 Issue/PR/Discussion 文件名
  - 测试 ResolveFilenameConflict 解决文件名冲突
- [ ] **Task 3.1.2** 创建 `internal/formatter/formatter.go` - 实现格式化
  - 实现 `SanitizeFilename(s string)` 函数
  - 实现 `GenerateFilename(owner, repo, issueType string, number int, title string)` 函数
  - 实现 `ResolveFilenameConflict(dir, filename string)` 函数

### 3.2 Markdown 转换 (converter)
- [ ] **Task 3.2.1** 创建 `internal/converter/converter_test.go` - 表格驱动测试
  - 测试 Convert Issue 转 Markdown
  - 测试 Convert PR 转 Markdown
  - 测试 Convert Discussion 转 Markdown
  - 测试评论排序 (asc/desc)
- [ ] **Task 3.2.2** 创建 `internal/converter/converter.go` - 实现转换
  - 创建 `Converter` 结构体 (Order 字段)
  - 实现 `NewConverter(order string)` 构造函数
  - 实现 `Convert(content *types.GitHubContent)` 方法

---

## Phase 4: CLI Assembly (命令行入口集成)

### 4.1 CLI 入口
- [ ] **Task 4.1.1** 创建 `cmd/issue2md/main.go` - CLI 主程序
  - 定义命令行 flags (--output/-o, --order)
  - 实现 URL 批量处理逻辑
  - 集成 parser, fetcher, converter, formatter
  - 实现错误处理和 stderr 输出
  - 实现批量处理时单个失败继续执行的逻辑

---

## 验证阶段

### 手动验证 (按照 spec.md 验收标准)
- [ ] **Task V1** TC001: 单个 Issue URL - 执行 CLI，检查输出文件内容
- [ ] **Task V2** TC002: 单个 PR URL - 执行 CLI，检查输出文件内容
- [ ] **Task V3** TC003: 单个 Discussion URL - 执行 CLI，检查输出文件内容
- [ ] **Task V4** TC004: 批量 3 个 URL - 执行 CLI，检查生成 3 个文件
- [ ] **Task V5** TC005: 指定输出目录 - 使用 `-o /tmp`，检查文件位置
- [ ] **Task V6** TC006: --order asc - 检查评论顺序
- [ ] **Task V7** TC007: 默认 order - 检查评论倒序
- [ ] **Task V8** TC008: 批量错误处理 - 故意输入错误 URL，检查错误打印
- [ ] **Task V9** TC009: 文件名特殊字符 - Issue 标题含 `/` 和空格，检查文件名
- [ ] **Task V10** TC010: 文件名冲突 - 同一 URL 执行两次，检查 `_1` 后缀
- [ ] **Task V11** TC011: 404 错误 - 输入不存在 Issue，检查错误信息
- [ ] **Task V12** TC012: 公开仓库无 Token - 不设置 GITHUB_TOKEN，检查正常获取
- [ ] **Task V13** TC013: 私有仓库有 Token - 设置 GITHUB_TOKEN，检查正常获取

---

## 依赖关系图

```
Phase 1
  └── Task 1.2.1: types.go (无依赖)

Phase 2
  ├── Task 2.1.1: parser_test.go (依赖 Task 1.2.1)
  ├── Task 2.1.2: parser.go (依赖 Task 2.1.1)
  ├── Task 2.2.1: fetcher_test.go (依赖 Task 1.2.1)
  └── Task 2.2.2: fetcher.go (依赖 Task 2.1.2, Task 2.2.1)

Phase 3
  ├── Task 3.1.1: formatter_test.go (无依赖)
  ├── Task 3.1.2: formatter.go (依赖 Task 3.1.1)
  ├── Task 3.2.1: converter_test.go (依赖 Task 1.2.1)
  └── Task 3.2.2: converter.go (依赖 Task 3.2.1)

Phase 4
  └── Task 4.1.1: main.go (依赖 Task 2.1.2, 2.2.2, 3.1.2, 3.2.2)

并行任务标记 [P]: Task 1.1.1, 1.1.2 可以并行执行
```

---

## 执行顺序建议

1. 首先执行 Phase 1 任务 (1.1.1, 1.1.2, 1.2.1)
2. 然后执行 Phase 2 任务 (按依赖顺序: parser 测试 -> parser 实现 -> fetcher 测试 -> fetcher 实现)
3. 执行 Phase 3 任务 (formatter 和 converter 可以并行进行)
4. 执行 Phase 4 任务 (main.go)
5. 最后执行验证任务
