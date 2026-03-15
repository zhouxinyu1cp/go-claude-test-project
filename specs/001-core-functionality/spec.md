# issue2md Specification

## 1. 用户故事

### 1.1 CLI 版本（当前）
作为开发者，我想要通过命令行工具将 GitHub Issue/PR/Discussion 转换为 Markdown 文件，以便我可以离线保存和归档项目讨论。

**典型场景：**
- 整理项目文档时，将相关的 GitHub Issue 讨论内容拉取下来
- 归档重要的 PR 审查过程
- 备份 Decision Record（使用 Discussion）

### 1.2 Web 版（未来规划）
作为非技术用户，我想要通过网页输入 GitHub URL 来生成 Markdown 文件，无需安装命令行工具。

---

## 2. 功能性需求

### 2.1 URL 识别与解析
| 需求ID | 描述 |
|--------|------|
| F001 | 工具必须支持解析三种 GitHub URL 类型：Issue、PR、Discussion |
| F002 | URL 格式：`https://github.com/{owner}/{repo}/issues/{number}`、`https://github.com/{owner}/{repo}/pull/{number}`、`https://github.com/{owner}/{repo}/discussions/{number}` |
| F003 | 工具应能自动识别 URL 类型并调用对应的 GitHub API |
| F004 | 对于不支持的 URL 类型，应返回明确的错误信息 |

### 2.2 GitHub API 交互
| 需求ID | 描述 |
|--------|------|
| F005 | 通过 GitHub REST API 获取 Issue/PR/Discussion 数据 |
| F006 | 获取对应的评论列表（Issue Comments、PR Reviews、Discussion Comments） |
| F007 | 支持通过环境变量 `GITHUB_TOKEN` 传入 Personal Access Token |
| F008 | 仅支持公开仓库（无需认证也可工作） |
| F009 | 透传 GitHub API 的错误信息（如 404、资源不存在、Rate Limit）给用户 |

### 2.3 命令行参数（Flags）
| Flag | 简写 | 描述 | 默认值 |
|------|------|------|--------|
| `--output` | `-o` | 指定输出目录 | 当前目录 (`./`) |
| `--order` | 无 | 评论排序方向：`asc`（正序）或 `desc`（倒序） | `desc` |

| 需求ID | 描述 |
|--------|------|
| F010 | 支持单 URL 处理 |
| F011 | 支持批量处理多个 URL（空格分隔或换行分隔） |
| F012 | 当批量处理时，遇到错误应打印错误并继续处理下一个 URL |

### 2.4 文件名生成
| 需求ID | 描述 |
|--------|------|
| F013 | Issue 文件名格式：`{owner}-{repo}-{issue_number}-{title_slug}.md` |
| F014 | PR 文件名格式：`{owner}-{repo}-pr-{pr_number}-{title_slug}.md` |
| F015 | Discussion 文件名格式：`{owner}-{repo}-discussion-{number}-{title_slug}.md` |
| F016 | 标题中的非法字符（`/`, `:`, 空格等）替换为 `-` |
| F017 | 批量处理时若文件名冲突，在文件名末尾添加 `_1`, `_2` 等序号 |

### 2.5 Markdown 输出格式

#### 2.5.1 Issue 输出结构
```markdown
# {issue_title}

- **Author**: {login}
- **Created**: {created_at}
- **Status**: {state} ({open/closed})
- **Labels**: {label1, label2}
- **URL**: {issue_html_url}

---

{body}

---

## Comments ({comment_count})

### {comment_author} • {created_at}
{comment_body}

### {comment_author} • {created_at}
{comment_body}
...
```

#### 2.5.2 PR 输出结构
```markdown
# PR #{pr_number}: {pr_title}

- **Author**: {login}
- **Created**: {created_at}
- **Status**: {state}
- **Merged**: {merged_at or N/A}
- **URL**: {pr_html_url}

---

{body}

---

## Reviews ({review_count})

### {review_author} • {created_at}
{review_body}

### {review_author} • {created_at}
{review_body}
...
```

#### 2.5.3 Discussion 输出结构
```markdown
# Discussion: {discussion_title}

- **Author**: {login}
- **Created**: {created_at}
- **Status**: {state}
- **Category**: {category_name}
- **URL**: {discussion_html_url}

---

{body}

---

## Comments ({comment_count})

### {comment_author} • {created_at}
{comment_body}

### {comment_author} • {created_at}
{comment_body}
...
```

| 需求ID | 描述 |
|--------|------|
| F018 | 输出文件必须包含：标题、作者、创建时间、状态、主楼内容、所有评论 |
| F019 | 评论默认按**倒序**（最新在前）显示 |
| F020 | 支持通过 `--order asc` 参数改为正序（最早在前） |
| F021 | 保留原始 URL 链接（可点击） |
| F022 | 代码块保持原始 Markdown 格式 |
| F023 | 图片保留 URL 引用（不下载到本地） |
| F024 | @mention 用户直接忽略（不解析） |

---

## 3. 非功能性需求

### 3.1 架构设计
| 需求ID | 描述 |
|--------|------|
| NF001 | 核心逻辑与 I/O 解耦，便于后续扩展 Web 版 |
| NF002 | 使用 Go 标准库 `net/http` 发送 HTTP 请求 |
| NF003 | 使用 Go 标准库 JSON 解析（`encoding/json`） |

### 3.2 错误处理
| 需求ID | 描述 |
|--------|------|
| NF004 | 所有错误信息输出到 stderr |
| NF005 | 404 错误（资源不存在）应明确提示 "Resource not found" |
| NF006 | 网络超时透传 GitHub API 错误信息 |
| NF007 | Rate Limit 错误透传 GitHub API 错误信息（含 Retry-After） |

### 3.3 依赖管理
| 需求ID | 描述 |
|--------|------|
| NF008 | 严格遵循"标准库优先"原则，不引入非必要第三方依赖 |
| NF009 | 如需依赖外部库，须在 `constitution.md` 框架下评估必要性 |

---

## 4. 验收标准

### 4.1 功能验收

| 测试 Case | 输入 | 预期输出 |
|-----------|------|----------|
| TC001 | 单个 Issue URL | 生成正确的 Issue Markdown 文件 |
| TC002 | 单个 PR URL | 生成正确的 PR Markdown 文件 |
| TC003 | 单个 Discussion URL | 生成正确的 Discussion Markdown 文件 |
| TC004 | 批量 3 个 URL（1 Issue + 1 PR + 1 Discussion） | 生成 3 个正确的 Markdown 文件 |
| TC005 | 使用 `--output /tmp` 指定目录 | 文件保存到 `/tmp` 目录 |
| TC006 | 使用 `--order asc` | 评论按正序显示 |
| TC007 | 默认（无 `--order`） | 评论按倒序显示 |
| TC008 | 批量处理时第 2 个 URL 错误 | 打印错误，继续处理第 3 个 URL |
| TC009 | Issue 标题包含 `/` 和空格 | 文件名中 `/` 和空格替换为 `-` |
| TC010 | 批量处理时文件名冲突 | 自动添加 `_1`, `_2` 后缀 |
| TC011 | 不存在的 Issue URL | 打印 "Resource not found" 错误并退出 |
| TC012 | 不使用 `GITHUB_TOKEN`（公开仓库） | 正常工作 |
| TC013 | 使用 `GITHUB_TOKEN`（私有仓库） | 正常获取内容 |

### 4.2 输出格式验收

| 测试 Case | 检查点 |
|-----------|--------|
| TC014 | Issue 输出包含：标题、Author、Created、Status、Labels、URL |
| TC015 | PR 输出包含：标题、Author、Created、Status、Merged、URL |
| TC016 | Discussion 输出包含：标题、Author、Created、Status、Category、URL |
| TC017 | 评论内容保留原始代码块格式 |
| TC018 | 评论中的图片保留 URL 引用 |
| TC019 | 主楼和评论之间用 `---` 分隔 |

---

## 5. 输出格式示例

### 5.1 Issue 示例

**输入：**
```
https://github.com/gorilla/mux/issues/123
```

**输出文件：`gorilla-mux-issue-123-refactor-handler.md`：**
```markdown
# Refactor handler registration

- **Author**: johndoe
- **Created**: 2024-01-15T10:30:00Z
- **Status**: closed
- **Labels**: enhancement, needs-review
- **URL**: https://github.com/gorilla/mux/issues/123

---

We should refactor the handler registration to support more flexible patterns.

```go
func Register(pattern string, handler http.Handler) error {
    // ...
}
```

---

## Comments (3)

### alice • 2024-01-15T11:00:00Z
This looks great! Could you also add unit tests?

### bob • 2024-01-15T14:30:00Z
+1, would love to see this merged.

### johndoe • 2024-01-16T09:00:00Z
Added tests in the latest commit.
```

### 5.2 PR 示例

**输入：**
```
https://github.com/gorilla/mux/pull/456
```

**输出文件：`gorilla-mux-pr-456-add-middleware.md`：**
```markdown
# PR #456: Add middleware support

- **Author**: janedoe
- **Created**: 2024-02-01T08:00:00Z
- **Status**: merged
- **Merged**: 2024-02-03T16:00:00Z
- **URL**: https://github.com/gorilla/mux/pull/456

---

This PR adds middleware support to the router.

---

## Reviews (2)

### reviewer1 • 2024-02-01T10:00:00Z
Please add tests for the middleware.

### reviewer2 • 2024-02-02T12:00:00Z
LGTM! Ship it.
```

### 5.3 Discussion 示例

**输入：**
```
https://github.com/gorilla/mux/discussions/789
```

**输出文件：`gorilla-mux-discussion-789-roadmap-2024.md`：**
```markdown
# Discussion: Roadmap 2024

- **Author**: maintainer
- **Created**: 2024-01-01T00:00:00Z
- **Status**: closed
- **Category**: Q&A
- **URL**: https://github.com/gorilla/mux/discussions/789

---

What features should we prioritize for 2024?

---

## Comments (2)

### contributor1 • 2024-01-02T10:00:00Z
We should focus on performance improvements.

### maintainer • 2024-01-03T12:00:00Z
Agreed, let's create a dedicated issue for this.
```

---

## 6. 技术实现提示

### 6.1 GitHub API 端点
- Issue: `GET /repos/{owner}/{repo}/issues/{issue_number}`
- PR: `GET /repos/{owner}/{repo}/pulls/{pull_number}`
- Discussion: `GET /repos/{owner}/{repo}/discussions/{discussion_number}`
- Issue Comments: `GET /repos/{owner}/{repo}/issues/{issue_number}/comments`
- PR Review Comments: `GET /repos/{owner}/{repo}/pulls/{pull_number}/comments`
- Discussion Comments: `GET /repos/{owner}/{repo}/discussions/{discussion_number}/comments`

### 6.2 目录结构建议
```
issue2md/
├── cmd/
│   └── issue2md/
│       └── main.go
├── internal/
│   ├── fetcher/
│   │   └── fetcher.go      # GitHub API 请求逻辑
│   ├── parser/
│   │   └── parser.go       # URL 解析逻辑
│   ├── converter/
│   │   └── converter.go    # 转换为 Markdown
│   └── formatter/
│       └── formatter.go   # 文件名生成、字符串处理
└── spec.md
```
