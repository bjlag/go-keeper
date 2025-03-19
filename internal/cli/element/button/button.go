package button

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/bjlag/go-keeper/internal/cli/style"
)

// Button описывает кнопку в UI.
type Button struct {
	// text содержит текст кнопки.
	text string
	// focus признак имеет ли кнопка фокус.
	focus bool
	// FocusedStyle стиль кнопки, когда она в состоянии фокуса
	FocusedStyle lipgloss.Style
	// BlurredStyle стиль кнопки, когда она не в фокусе.
	BlurredStyle lipgloss.Style
}

// NewButton создает экземпляр кнопки.
func NewButton(text string) Button {
	return Button{
		text:         text,
		FocusedStyle: lipgloss.NewStyle(),
		BlurredStyle: lipgloss.NewStyle(),
	}
}

// Option тип опции кнопки.
type Option func(m *Button)

// WithFocused меняет состояние кнопки на фокус.
func WithFocused() Option {
	return func(m *Button) {
		m.Focus()
	}
}

// CreateDefaultButton создает экземпляр кнопки с заранее примененными стилями.
// Через аргумент opts можно передать дополнительные настройки.
func CreateDefaultButton(text string, opts ...Option) Button {
	b := NewButton(text)
	b.FocusedStyle = style.FocusedStyle
	b.BlurredStyle = style.BlurredStyle

	for _, opt := range opts {
		opt(&b)
	}

	return b
}

// String формирует строковое представление кнопки для вывода в UI.
func (b *Button) String() string {
	if b.focus {
		return fmt.Sprintf("[ %s ]", b.FocusedStyle.Render(b.text))
	}
	return fmt.Sprintf("[ %s ]", b.BlurredStyle.Render(b.text))
}

// Focus меняет состояние на фокус.
func (b *Button) Focus() {
	b.focus = true
}

// Blur снимает фокус с кнопки.
func (b *Button) Blur() {
	b.focus = false
}

// Focused возвращает состояние фокуса кнопки.
func (b *Button) Focused() bool {
	return b.focus
}
