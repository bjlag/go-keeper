package style

import (
	"fmt"

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
	HelpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	FocusedButton = FocusedStyle.Render("[ Submit ]")
	BlurredButton = fmt.Sprintf("[ %s ]", BlurredStyle.Render("Submit"))

	SeparatorStyle = lipgloss.NewStyle().
			Margin(1, 0).
			Foreground(lipgloss.Color("245"))
	//BorderStyle(lipgloss.NormalBorder()).Margin(1, 0).
	//BorderForeground(lipgloss.Color("63")).
	//BorderBackground(lipgloss.Color("63")).
	//BorderBottom(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true).
			MarginLeft(0).
			PaddingLeft(0)
)
