// Package style содержит стили CLI приложения, которые используются для вывода UI.
package style

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderBottom(true).
			Foreground(lipgloss.Color("39")).
			Margin(1, 0)

	NoStyle      = lipgloss.NewStyle()
	CursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	BlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	ErrorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	ListTitleStyle        = lipgloss.NewStyle().MarginLeft(2)
	ListItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	SelectedListItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	ListPaginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	ListHelpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)

	ErrorBlockStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("245")).
			Padding(0, 1).
			MarginTop(1).
			Foreground(lipgloss.Color("196"))
)

func SetNoStyle(input textinput.Model) textinput.Model {
	input.PromptStyle = NoStyle
	input.TextStyle = NoStyle
	return input
}

func SetFocusStyle(input textinput.Model) textinput.Model {
	input.PromptStyle = FocusedStyle
	input.TextStyle = FocusedStyle
	return input
}

func SetErrorStyle(input textinput.Model) textinput.Model {
	input.PromptStyle = ErrorStyle
	return input
}
