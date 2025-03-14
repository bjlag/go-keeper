// Package register описывает модель для работы UI регистрации пользователя.
package register

import (
	"context"
	"errors"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bjlag/go-keeper/internal/cli/common"
	"github.com/bjlag/go-keeper/internal/cli/element/button"
	tinput "github.com/bjlag/go-keeper/internal/cli/element/textinput"
	"github.com/bjlag/go-keeper/internal/cli/message"
	"github.com/bjlag/go-keeper/internal/cli/style"
	"github.com/bjlag/go-keeper/internal/infrastructure/validator"
	"github.com/bjlag/go-keeper/internal/usecase/client/master_key"
	"github.com/bjlag/go-keeper/internal/usecase/client/register"
)

const (
	posEmail int = iota
	posPassword
	posSubmitBtn
	posBackBtn
)

var errUserAlreadyRegistered = common.NewFormError("Пользователь уже зарегистрирован")

type Model struct {
	help     help.Model
	header   string
	elements []interface{}
	pos      int
	err      error

	loginModel tea.Model

	usecaseRegister  *register.Usecase
	usecaseMasterKey *master_key.Usecase
}

func InitModel(usecaseRegister *register.Usecase, usecaseMasterKey *master_key.Usecase) *Model {
	f := &Model{
		help:   help.New(),
		header: "Регистрация",
		elements: []interface{}{
			posEmail:     tinput.CreateDefaultTextInput("Email"),
			posPassword:  tinput.CreateDefaultTextInput("Пароль"),
			posSubmitBtn: button.CreateDefaultButton("Регистрация"),
			posBackBtn:   button.CreateDefaultButton("Назад"),
		},

		usecaseRegister:  usecaseRegister,
		usecaseMasterKey: usecaseMasterKey,
	}

	for i := range f.elements {
		if e, ok := f.elements[i].(textinput.Model); ok {
			if i == posEmail {
				e.TextStyle = style.FocusedStyle
				e.PromptStyle = style.FocusedStyle
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

func (f *Model) Init() tea.Cmd {
	return nil
}

func (f *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		for i := range f.elements {
			if e, ok := f.elements[i].(textinput.Model); ok {
				e.Width = msg.Width
			}
		}
		return f, nil
	case message.OpenRegisterMsg:
		f.loginModel = msg.LoginModel
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
						f.elements[i] = style.SetFocusStyle(e)
						continue
					}

					e.Blur()
					f.elements[i] = style.SetNoStyle(e)
				case button.Button:
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
		case key.Matches(msg, common.Keys.Back):
			return f.loginModel.Update(message.BackMsg{})
		case key.Matches(msg, common.Keys.Enter):
			f.err = nil

			switch {
			case f.pos == posSubmitBtn || f.pos == posEmail || f.pos == posPassword:
				return f.submit()
			case f.pos == posBackBtn:
				return f.loginModel.Update(message.OpenLoginMsg{})
			}

			return f, nil
		}
	}

	return f, f.updateInputs(msg)
}

func (f *Model) View() string {
	var b strings.Builder

	b.WriteString(style.TitleStyle.Render(f.header))
	b.WriteRune('\n')

	for i := range f.elements {
		if e, ok := f.elements[i].(textinput.Model); ok {
			b.WriteString(e.View())
			b.WriteRune('\n')
		}
	}

	b.WriteRune('\n')

	for i := range f.elements {
		if e, ok := f.elements[i].(button.Button); ok {
			b.WriteString(e.String())
			b.WriteRune('\n')
		}
	}

	var (
		errValidate *common.ValidateError
		errForm     *common.FormError
	)

	// выводим ошибки валидации
	if f.err != nil && (errors.As(f.err, &errValidate) || errors.As(f.err, &errForm)) {
		b.WriteString(style.ErrorBlockStyle.Render(f.err.Error()))
		b.WriteRune('\n')
	}

	b.WriteRune('\n')
	b.WriteString(f.help.View(common.Keys))

	// выводим прочие ошибки
	if f.err != nil && !(errors.As(f.err, &errValidate) || errors.As(f.err, &errForm)) {
		b.WriteRune('\n')
		b.WriteString(style.ErrorBlockStyle.Render(f.err.Error()))
	}

	return b.String()
}

func (f *Model) submit() (tea.Model, tea.Cmd) {
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
		f.elements[posEmail] = style.SetErrorStyle(email)
		errValidate.AddError("Неправильный email")
	}

	if !validator.ValidatePassword(password.Value()) {
		f.elements[posPassword] = style.SetErrorStyle(password)
		errValidate.AddError("Недостаточно сложный пароль")
	}

	if errValidate.HasErrors() {
		f.err = errValidate
		return f, nil
	}

	err := f.usecaseRegister.Do(context.TODO(), register.Data{
		Email:    email.Value(),
		Password: password.Value(),
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.AlreadyExists {
				f.err = errUserAlreadyRegistered
				return f, nil
			}
			f.err = common.NewFormError(s.Message())
			return f, nil
		}

		f.err = err
		return f, nil
	}

	err = f.usecaseMasterKey.Do(context.TODO(), master_key.Data{
		Password: password.Value(),
	})
	if err != nil {
		f.err = err
		return f, nil
	}

	return f.loginModel.Update(message.SuccessMsg{})
}

func (f *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(f.elements))

	for i := range f.elements {
		if m, ok := f.elements[i].(textinput.Model); ok {
			f.elements[i], cmds[i] = m.Update(msg)
		}
	}

	return tea.Batch(cmds...)
}
