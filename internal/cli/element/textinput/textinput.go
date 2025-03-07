package textinput

import (
	"github.com/charmbracelet/bubbles/textinput"

	"github.com/bjlag/go-keeper/internal/cli/style"
)

type Option func(m *textinput.Model)

func WithFocused() Option {
	return func(m *textinput.Model) {
		m.Focus()
	}
}

func WithValue(value string) Option {
	return func(m *textinput.Model) {
		m.SetValue(value)
	}
}

func CreateDefaultTextInput(placeholder string, opts ...Option) textinput.Model {
	m := textinput.New()

	m.Cursor.Style = style.CursorStyle
	m.PlaceholderStyle = style.BlurredStyle
	m.CharLimit = 50
	m.Placeholder = placeholder

	for _, opt := range opts {
		opt(&m)
	}

	return m
}
