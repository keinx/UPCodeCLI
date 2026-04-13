package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CLI 工具配置
type CLITool struct {
	Name        string
	Command     string
	UpdateCmd   string
	VersionFlag string
	// 如果 UpdateCommand 不为空，则使用它作为完整的更新命令（通过 shell 执行）
	// 否则使用 Command + UpdateCmd 的方式
	UpdateCommand string
	// 更新成功后执行的额外命令（如环境变量刷新等）
	PostUpdateCommand string
}

var cliTools = []CLITool{
	{Name: "1. ClaudeCode", Command: "claude", UpdateCmd: "update", VersionFlag: "--version"},
	{Name: "2. OpenCode", Command: "opencode", UpdateCmd: "upgrade", VersionFlag: "--version"},
	{Name: "3. CodeBuddy", Command: "codebuddy", UpdateCmd: "update", VersionFlag: "--version"},
	{Name: "4. QoderCLI", Command: "qodercli", UpdateCmd: "update", VersionFlag: "--version"},
	{Name: "+ QwenCodeCLI", Command: "qwen", UpdateCmd: "", VersionFlag: "--version", UpdateCommand: "npm install -g @qwen-code/qwen-code@latest --registry https://registry.npmmirror.com"},
	{Name: "+ GeminiCLI", Command: "gemini", UpdateCmd: "", VersionFlag: "--version", UpdateCommand: "npm install -g @google/gemini-cli@latest --registry https://registry.npmmirror.com"},
	{Name: "+ GrokCLI", Command: "grok", UpdateCmd: "", VersionFlag: "--version", UpdateCommand: "npm install -g @vibe-kit/grok-cli@latest --registry https://registry.npmmirror.com"},
	{Name: "* Git", Command: "git", UpdateCmd: "update-git-for-windows", VersionFlag: "--version"},
	{Name: "* FNM", Command: "fnm", UpdateCmd: "", VersionFlag: "--version", UpdateCommand: "winget upgrade fnm", PostUpdateCommand: "fnm env --use-on-cd | Out-String | Invoke-Expression"},
}

// 应用信息
const (
	AppVersion = "v0.2"
	AppAuthor  = "keinx"
)

// 样式定义
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5F5FFF")).
			Bold(true)

	versionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F5F")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5F5FFF"))

	boxStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Margin(1, 2)
)

// 应用状态
type state int

const (
	stateSelect state = iota
	stateUpdating
	stateDone
)

// 消息类型
type updateCompleteMsg struct {
	tool       CLITool
	success    bool
	message    string
	oldVersion string
}

// 获取版本消息
type versionCheckedMsg struct {
	tool    CLITool
	version string
}

// 模型
type model struct {
	list         list.Model
	state        state
	tool         CLITool
	message      string
	success      bool
	currentVer   string
}

func initialModel() model {
	items := make([]list.Item, len(cliTools))
	for i, tool := range cliTools {
		items[i] = item{tool: tool}
	}

	// 自定义委托样式
	delegate := list.NewDefaultDelegate()
	delegate.SetHeight(2) // 每个项目占2行（标题+描述）
	// 选中状态：明亮鲜艳的颜色
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#00D7FF")).
		Bold(true)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("#5FD7FF"))
	// 未选中状态：置灰
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(lipgloss.Color("#666666"))
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.
		Foreground(lipgloss.Color("#444444"))

	l := list.New(items, delegate, 40, 30)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	l.Title = "CLI 工具更新器"
	l.Styles.Title = titleStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return model{
		list:  l,
		state: stateSelect,
	}
}

// 列表项
type item struct {
	tool CLITool
}

func (i item) FilterValue() string {
	return i.tool.Name
}

func (i item) Title() string {
	return i.tool.Name
}

func (i item) Description() string {
	return fmt.Sprintf("命令: %s", i.tool.Command)
}

// 执行命令并获取输出
func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

// 获取版本号
func getVersion(tool CLITool) string {
	output, err := runCommand(tool.Command, tool.VersionFlag)
	if err != nil {
		return "未安装或无法获取版本"
	}
	return output
}

// 更新命令消息
func checkVersion(tool CLITool) tea.Cmd {
	return func() tea.Msg {
		version := getVersion(tool)
		return versionCheckedMsg{
			tool:    tool,
			version: version,
		}
	}
}

func doUpdate(tool CLITool, oldVersion string) tea.Cmd {
	return func() tea.Msg {
		message, success := performUpdateWithVersion(tool, oldVersion)
		return updateCompleteMsg{
			tool:    tool,
			success: success,
			message: message,
			oldVersion: oldVersion,
		}
	}
}

func runShellCommand(command string) (string, error) {
	cmd := exec.Command("cmd", "/c", command)
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func performUpdateWithVersion(tool CLITool, oldVersion string) (string, bool) {
	if strings.Contains(oldVersion, "未安装") {
		return oldVersion, false
	}

	// 执行更新命令
	var updateOutput string
	var err error
	if tool.UpdateCommand != "" {
		updateOutput, err = runShellCommand(tool.UpdateCommand)
	} else {
		updateOutput, err = runCommand(tool.Command, tool.UpdateCmd)
	}
	if err != nil {
		// winget 等工具在已是最新版本时也会返回错误码
		if strings.Contains(updateOutput, "较新的包版本") || strings.Contains(updateOutput, "no newer") {
			return fmt.Sprintf("已是最新版本\n当前版本: %s", oldVersion), true
		}
		return fmt.Sprintf("更新失败: %s\n错误: %s", updateOutput, err.Error()), false
	}

	// 获取更新后版本
	newVersion := getVersion(tool)

	// 执行更新后命令（如环境变量刷新等）
	if tool.PostUpdateCommand != "" {
		runShellCommand(tool.PostUpdateCommand)
	}

	if oldVersion == newVersion {
		return fmt.Sprintf("已是最新版本\n当前版本: %s", oldVersion), true
	}

	return fmt.Sprintf("更新成功!\n更新前: %s\n更新后: %s", oldVersion, newVersion), true
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case stateSelect:
			switch msg.String() {
			case "enter":
				i, ok := m.list.SelectedItem().(item)
				if ok {
					m.state = stateUpdating
					m.tool = i.tool
					return m, checkVersion(i.tool)
				}
			case "ctrl+c", "q":
				return m, tea.Quit
			}

		case stateDone:
			switch msg.String() {
			case "enter", "esc":
				// 返回选择列表
				m.state = stateSelect
				m.message = ""
				m.currentVer = ""
				m.list.Select(0)
				return m, nil
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		}

	case versionCheckedMsg:
		m.currentVer = msg.version
		return m, doUpdate(msg.tool, msg.version)

	case updateCompleteMsg:
		m.state = stateDone
		m.message = msg.message
		m.success = msg.success
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	switch m.state {
	case stateUpdating:
		var lines []string
		lines = append(lines, fmt.Sprintf("正在更新 %s，请稍候...", m.tool.Name))
		if m.currentVer != "" {
			lines = append(lines, fmt.Sprintf("当前版本: %s", m.currentVer))
		}
		return boxStyle.Render(strings.Join(lines, "\n"))

	case stateDone:
		var lines []string
		lines = append(lines, fmt.Sprintf("📦 %s", m.tool.Name))
		lines = append(lines, "")
		if m.success {
			lines = append(lines, statusStyle.Render("✓ "+m.message))
		} else {
			lines = append(lines, errorStyle.Render("✗ "+m.message))
		}
		lines = append(lines, "")
		lines = append(lines, infoStyle.Render("按 Enter 返回主菜单，按 Q 退出"))
		return boxStyle.Render(strings.Join(lines, "\n"))

	default:
		header := fmt.Sprintf("版本 %s | 版权 %s", AppVersion, AppAuthor)
		return versionStyle.Render(header) + "\n" + m.list.View() + "\n" + infoStyle.Render("按 Enter 选择更新，按 Q 退出")
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
