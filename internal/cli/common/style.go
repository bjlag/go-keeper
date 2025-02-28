package common

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	"github.com/bjlag/go-keeper/internal/cli/element"
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderBottom(true).
			Foreground(lipgloss.Color("39")).
			Margin(1, 0)

	NoStyle      = lipgloss.NewStyle()
	TextStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	CursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	BlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	ErrorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	BlockStyle      = TextStyle.Margin(1)
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

func CreateDefaultTextInput(placeholder string, limit int) textinput.Model {
	m := textinput.New()

	m.Cursor.Style = CursorStyle
	m.PlaceholderStyle = BlurredStyle
	m.CharLimit = limit
	m.Placeholder = placeholder

	return m
}

func CreateDefaultButton(text string) element.Button {
	b := element.NewButton(text)
	b.FocusedStyle = FocusedStyle
	b.BlurredStyle = BlurredStyle

	return b
}
