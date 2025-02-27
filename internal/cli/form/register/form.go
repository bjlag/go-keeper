package register

import (
	"errors"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/cli/common"
	"github.com/bjlag/go-keeper/internal/cli/element"
	"github.com/bjlag/go-keeper/internal/cli/message"
	"github.com/bjlag/go-keeper/internal/infrastructure/validator"
)

const (
	posEmail int = iota
	posPassword
	posSubmitBtn
	posBackBtn

	emailCharLimit    = 20
	passwordCharLimit = 20
)

type Form struct {
	main      tea.Model
	prevModel tea.Model
	help      help.Model
	header    string
	elements  []interface{}
	pos       int
	err       error
}

func NewForm(main, prevModel tea.Model) *Form {
	f := &Form{
		main:      main,
		prevModel: prevModel,
		help:      help.New(),
		header:    "Регистрация",
		elements: []interface{}{
			posEmail:     common.CreateDefaultTextInput("Email", emailCharLimit),
			posPassword:  common.CreateDefaultTextInput("Password", passwordCharLimit),
			posSubmitBtn: common.CreateDefaultButton("Регистрация"),
			posBackBtn:   common.CreateDefaultButton("Назад"),
		},
	}

	for i := range f.elements {
		switch e := f.elements[i].(type) {
		case textinput.Model:
			if i == posEmail {
				e.TextStyle = common.FocusedStyle
				e.PromptStyle = common.FocusedStyle
				e.Focus()

				f.elements[i] = e
				continue
			}

			e.EchoMode = textinput.EchoPassword
			e.EchoCharacter = '•'

			f.elements[i] = e
		}
	}

	return f
}

func (f *Form) Init() tea.Cmd {
	return nil
}

func (f *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.Keys.Quit):
			return f, tea.Quit
		case key.Matches(msg, common.Keys.Navigation):
			if key.Matches(msg, common.Keys.Down, common.Keys.Tab) {
				f.pos++
			} else {
				f.pos--
			}

			if f.pos > len(f.elements)-1 {
				f.pos = 0
			} else if f.pos < 0 {
				f.pos = len(f.elements) - 1
			}

			for i := range f.elements {
				switch e := f.elements[i].(type) {
				case textinput.Model:
					if i == f.pos {
						e.Focus()
						f.elements[i] = common.SetFocusStyle(e)
						continue
					}

					e.Blur()
					f.elements[i] = common.SetNoStyle(e)
				case element.Button:
					if i == f.pos {
						e.Focus()
						f.elements[i] = e
						continue
					}
					e.Blur()
					f.elements[i] = e
				}
			}

			return f, nil
		case key.Matches(msg, common.Keys.Enter):
			f.err = nil

			switch {
			case f.pos == posSubmitBtn || f.pos == posEmail || f.pos == posPassword:
				return f.submit()
			case f.pos == posBackBtn:
				return f.prevModel.Update(tea.ClearScreen)
			}

			return f, nil
		}
	}

	return f, f.updateInputs(msg)
}

func (f *Form) View() string {
	var b strings.Builder

	b.WriteString(common.TitleStyle.Render(f.header))
	b.WriteRune('\n')

	for i := range f.elements {
		if e, ok := f.elements[i].(textinput.Model); ok {
			b.WriteString(e.View())
			b.WriteRune('\n')
		}
	}

	b.WriteRune('\n')

	for i := range f.elements {
		if e, ok := f.elements[i].(element.Button); ok {
			b.WriteString(e.String())
			b.WriteRune('\n')
		}
	}

	var errValidate *common.ValidateError

	// выводим ошибки валидации
	if f.err != nil && errors.As(f.err, &errValidate) {
		b.WriteString(common.ErrorBlockStyle.Render(f.err.Error()))
		b.WriteRune('\n')
	}

	b.WriteRune('\n')
	b.WriteString(f.help.View(common.Keys))

	// выводим прочие ошибки
	if f.err != nil && !errors.As(f.err, &errValidate) {
		b.WriteRune('\n')
		b.WriteString(common.ErrorBlockStyle.Render(f.err.Error()))
	}

	return b.String()
}

func (f *Form) submit() (tea.Model, tea.Cmd) {
	errValidate := common.NewValidateError()

	email, ok := f.elements[posEmail].(textinput.Model)
	if !ok {
		f.err = common.ErrInvalidElement
		return f, nil
	}
	password, ok := f.elements[posPassword].(textinput.Model)
	if !ok {
		f.err = common.ErrInvalidElement
		return f, nil
	}

	if !validator.ValidateEmail(email.Value()) {
		f.elements[posEmail] = common.SetErrorStyle(email)
		errValidate.AddError("Неправильный email")
	}

	if !validator.ValidatePassword(password.Value()) {
		f.elements[posPassword] = common.SetErrorStyle(password)
		errValidate.AddError("Недостаточно сложный пароль")
	}

	if errValidate.HasErrors() {
		f.err = errValidate
		return f, nil
	}

	// TODO: выполняем регистрацию

	return f.main.Update(message.RegisterSuccessMessage{
		AccessToken:  "xxxx",
		RefreshToken: "yyyy",
	})
}

func (f *Form) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(f.elements))

	for i := range f.elements {
		if m, ok := f.elements[i].(textinput.Model); ok {
			f.elements[i], cmds[i] = m.Update(msg)
		}
	}

	return tea.Batch(cmds...)
}
