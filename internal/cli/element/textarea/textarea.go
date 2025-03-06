package textarea

import (
	"github.com/charmbracelet/bubbles/textarea"
)

type Option func(m *textarea.Model)

func WithFocused() Option {
	return func(m *textarea.Model) {
		m.Focus()
	}
}

func WithValue(value string) Option {
	return func(m *textarea.Model) {
		m.SetValue(value)
	}
}

func CreateDefaultTextArea(placeholder string, opts ...Option) textarea.Model {
	m := textarea.New()
	m.Placeholder = placeholder

	for _, opt := range opts {
		opt(&m)
	}

	return m
}
