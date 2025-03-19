package textinput

import (
	"github.com/charmbracelet/bubbles/textinput"

	"github.com/bjlag/go-keeper/internal/cli/style"
)

// Option описывает тип опции элемента.
type Option func(m *textinput.Model)

// WithFocused настроит элемент в фокусе.
func WithFocused() Option {
	return func(m *textinput.Model) {
		m.Focus()
	}
}

// WithValue настроит элемент с определенным значением.
func WithValue(value string) Option {
	return func(m *textinput.Model) {
		m.SetValue(value)
	}
}

// CreateDefaultTextInput создает [textinput.Model] с заранее определенными настройками.
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
