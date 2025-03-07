package button

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/bjlag/go-keeper/internal/cli/style"
)

type Button struct {
	text         string
	focus        bool
	FocusedStyle lipgloss.Style
	BlurredStyle lipgloss.Style
}

func NewButton(text string) Button {
	return Button{
		text:         text,
		FocusedStyle: lipgloss.NewStyle(),
		BlurredStyle: lipgloss.NewStyle(),
	}
}

type Option func(m *Button)

func WithFocused() Option {
	return func(m *Button) {
		m.Focus()
	}
}

func CreateDefaultButton(text string, opts ...Option) Button {
	b := NewButton(text)
	b.FocusedStyle = style.FocusedStyle
	b.BlurredStyle = style.BlurredStyle

	for _, opt := range opts {
		opt(&b)
	}

	return b
}

func (b *Button) String() string {
	if b.focus {
		return fmt.Sprintf("[ %s ]", b.FocusedStyle.Render(b.text))
	}
	return fmt.Sprintf("[ %s ]", b.BlurredStyle.Render(b.text))
}

func (b *Button) Focus() {
	b.focus = true
}

func (b *Button) Blur() {
	b.focus = false
}

func (b *Button) Focused() bool {
	return b.focus
}
