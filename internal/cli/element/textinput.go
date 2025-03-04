package element

import (
	"github.com/charmbracelet/bubbles/textinput"

	"github.com/bjlag/go-keeper/internal/cli/style"
)

func CreateDefaultTextInput(placeholder string, limit int) textinput.Model {
	m := textinput.New()

	m.Cursor.Style = style.CursorStyle
	m.PlaceholderStyle = style.BlurredStyle
	m.CharLimit = limit
	m.Placeholder = placeholder

	return m
}
