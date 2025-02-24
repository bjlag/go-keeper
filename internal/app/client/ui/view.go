package ui

import (
	"github.com/bjlag/go-keeper/internal/app/client/ui/style"
	"strings"
)

func (m *Model) loginView() string {
	var b strings.Builder

	b.WriteString(style.TitleStyle.Render("АВТОРИЗАЦИЯ"))
	b.WriteRune('\n')

	help := "↑/↓: навигация • q: выход\n"
	b.WriteString(style.HelpStyle.Render(help))
	b.WriteRune('\n')

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if m.inputs[i].Err != nil {
			b.WriteRune('\n')
			b.WriteString(style.ErrorStyle.
				MarginLeft(2). //nolint:mnd
				Render(m.inputs[i].Err.Error()))
		}
		b.WriteRune('\n')
	}

	submitBtn := style.BlurredButton
	if m.focusIndex == len(m.inputs) {
		submitBtn = style.FocusedButton
	}

	b.WriteRune('\n')
	b.WriteString(submitBtn)

	//registerBtn := blurredButton

	return b.String()
}
