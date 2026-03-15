---
description: 分析git diff，生成符合Conventional Commits规范的提交信息并提交。
allowed-tools: Bash(git diff:*), Bash(git commit:*), Bash(git add:*)
---
1. 执行 `git add .` 添加当前项目的所有变更文件。
2. 执行 `git diff --staged` 获取暂存区的变更。
3. 根据变更内容，生成一条遵循 `CLAUDE.md` 中 **Conventional Commits** 规范的 Commit Message。
4. 向用户展示生成的 Message，并询问是否确认提交。
5. 如果确认，执行 `git commit -m "..."`。
