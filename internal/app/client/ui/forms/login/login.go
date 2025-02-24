package login

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	key2 "github.com/bjlag/go-keeper/internal/app/client/ui/key"
	"github.com/bjlag/go-keeper/internal/app/client/ui/style"
	"github.com/bjlag/go-keeper/internal/infrastructure/validator"
)

type Form struct {
	help     help.Model
	header   string
	email    textinput.Model
	password textinput.Model
}

func NewLoginForm() *Form {
	f := &Form{
		help:     help.New(),
		header:   "Авторизация",
		email:    textinput.New(),
		password: textinput.New(),
	}

	f.help.ShowAll = true

	f.email.Cursor.Style = style.CursorStyle
	f.email.PlaceholderStyle = style.BlurredStyle
	f.email.CharLimit = 20
	f.email.Placeholder = "Email"
	f.email.TextStyle = style.FocusedStyle
	f.email.PromptStyle = style.FocusedStyle
	f.email.Focus()

	f.password.Cursor.Style = style.CursorStyle
	f.password.PlaceholderStyle = style.BlurredStyle
	f.password.CharLimit = 60
	f.password.Placeholder = "Password"
	f.password.EchoMode = textinput.EchoPassword
	f.password.EchoCharacter = '•'

	return f
}

func (f *Form) Init() tea.Cmd {
	return nil
}

func (f *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key2.Keys.Quit):
			return f, tea.Quit
		case key.Matches(msg, key2.Keys.Up, key2.Keys.Down):
			if f.email.Focused() {
				f.email.Blur()
				f.password.Focus()

				f.email.PromptStyle = style.NoStyle
				f.email.TextStyle = style.NoStyle

				f.password.PromptStyle = style.FocusedStyle
				f.password.TextStyle = style.FocusedStyle

				return f, textarea.Blink
			}

			if f.password.Focused() {
				f.password.Blur()
				f.email.Focus()

				f.password.PromptStyle = style.NoStyle
				f.password.TextStyle = style.NoStyle

				f.email.PromptStyle = style.FocusedStyle
				f.email.TextStyle = style.FocusedStyle
			}

			return f, textarea.Blink
		case key.Matches(msg, key2.Keys.Enter):

			if !validator.ValidateEmail(f.email.Value()) {
				f.email.PromptStyle = style.ErrorStyle
				f.email.TextStyle = style.ErrorStyle
			}

			if f.password.Value() == "" {
				f.password.PromptStyle = style.ErrorStyle
				f.password.TextStyle = style.ErrorStyle
			}

			return f, textarea.Blink
		}
	}

	return f, f.updateInputs(msg)
}

func (f *Form) View() string {
	var b strings.Builder

	b.WriteString(style.TitleStyle.Render(f.header))
	b.WriteRune('\n')

	b.WriteString(f.email.View())
	b.WriteRune('\n')

	if f.email.Err != nil {
		b.WriteString(style.ErrorStyle.
			MarginLeft(2). //nolint:mnd
			Render(f.email.Err.Error()))
		b.WriteRune('\n')
	}

	b.WriteString(f.password.View())
	b.WriteRune('\n')

	if f.password.Err != nil {
		b.WriteString(style.ErrorStyle.
			MarginLeft(2). //nolint:mnd
			Render(f.password.Err.Error()))
		b.WriteRune('\n')
	}

	b.WriteRune('\n')
	b.WriteString(f.help.View(key2.Keys))

	return b.String()
}

func (f *Form) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 2)

	f.email, cmds[0] = f.email.Update(msg)
	f.password, cmds[1] = f.password.Update(msg)

	return tea.Batch(cmds...)
}
