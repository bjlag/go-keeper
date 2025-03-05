package element

import (
	"github.com/charmbracelet/bubbles/textinput"

	"github.com/bjlag/go-keeper/internal/cli/style"
)

type TextInputOption func(m *textinput.Model)

func WithFocused() TextInputOption {
	return func(m *textinput.Model) {
		m.Focus()
	}
}

func CreateDefaultTextInput(placeholder string, limit int, opts ...TextInputOption) textinput.Model {
	m := textinput.New()

	m.Cursor.Style = style.CursorStyle
	m.PlaceholderStyle = style.BlurredStyle
	m.CharLimit = limit
	m.Placeholder = placeholder

	for _, o := range opts {
		o(&m)
	}

	return m
}
