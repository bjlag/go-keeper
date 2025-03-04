package element

import (
	"fmt"
	"github.com/bjlag/go-keeper/internal/cli/style"
	"github.com/charmbracelet/lipgloss"
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

func CreateDefaultButton(text string) Button {
	b := NewButton(text)
	b.FocusedStyle = style.FocusedStyle
	b.BlurredStyle = style.BlurredStyle

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
