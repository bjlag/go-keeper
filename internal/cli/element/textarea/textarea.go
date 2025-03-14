package textarea

import (
	"github.com/charmbracelet/bubbles/textarea"
)

// Option описывает тип опции элемента.
type Option func(m *textarea.Model)

// WithFocused настроит элемент в фокусе.
func WithFocused() Option {
	return func(m *textarea.Model) {
		m.Focus()
	}
}

// WithValue настроит элемент с определенным значением.
func WithValue(value string) Option {
	return func(m *textarea.Model) {
		m.SetValue(value)
	}
}

// CreateDefaultTextArea создает [textarea.Model] с заранее определенными настройками.
func CreateDefaultTextArea(placeholder string, opts ...Option) textarea.Model {
	m := textarea.New()
	m.Placeholder = placeholder

	for _, opt := range opts {
		opt(&m)
	}

	return m
}
