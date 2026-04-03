# UPCodeCLI

一个用于批量更新 AI 编程助手 CLI 工具的 TUI 应用。

## 功能

- 支持多个 AI 编程助手 CLI 工具的一键更新：
  - Claude Code
  - OpenCode
  - CodeBuddy
  - QoderCLI
- 美观的终端用户界面 (TUI)
- 实时显示更新进度和版本信息

## 安装

```bash
go install github.com/keinx/UPCodeCLI@latest
```

或从源码构建：

```bash
git clone https://github.com/keinx/UPCodeCLI.git
cd UPCodeCLI
go build -o upcodecli
```

## 使用

```bash
upcodecli
```

### 操作说明

| 按键 | 功能 |
|------|------|
| ↑/↓ | 选择工具 |
| Enter | 执行更新 |
| Q / Ctrl+C | 退出 |
| Esc | 返回主菜单 |

## 环境要求

- Go 1.25+
- 已安装需要更新的 CLI 工具

## 依赖

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI 框架
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - 样式定义
- [Bubbles](https://github.com/charmbracelet/bubbles) - UI 组件

## 版本

当前版本：**v0.1**

## 版权

Copyright © 2026 keinx. All rights reserved.

## License

MIT License
