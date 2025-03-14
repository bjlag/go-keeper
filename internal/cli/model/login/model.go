package login

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
	"github.com/bjlag/go-keeper/internal/usecase/client/login"
	"github.com/bjlag/go-keeper/internal/usecase/client/master_key"
)

const (
	posEmail int = iota
	posPassword
	posSubmitBtn
	posRegisterBtn
	posCloseBtn
)

var errPasswordInvalid = common.NewFormError("Неверный email или пароль")

type Model struct {
	main     tea.Model
	help     help.Model
	header   string
	elements []interface{}
	pos      int
	err      error

	fromRegister tea.Model

	usecaseLogin     *login.Usecase
	usecaseMasterKey *master_key.Usecase
}

func InitModel(usecaseLogin *login.Usecase, usecaseMasterKey *master_key.Usecase, fromRegister tea.Model) *Model {
	f := &Model{
		help:   help.New(),
		header: "Авторизация",
		elements: []interface{}{
			posEmail:       tinput.CreateDefaultTextInput("Email"),
			posPassword:    tinput.CreateDefaultTextInput("Пароль"),
			posSubmitBtn:   button.CreateDefaultButton("Вход"),
			posRegisterBtn: button.CreateDefaultButton("Регистрация"),
			posCloseBtn:    button.CreateDefaultButton("Закрыть"),
		},

		fromRegister: fromRegister,

		usecaseLogin:     usecaseLogin,
		usecaseMasterKey: usecaseMasterKey,
	}

	for i := range f.elements {
		if e, ok := f.elements[i].(textinput.Model); ok {
			if i == posEmail {
				e.Focus()
				f.elements[i] = style.SetFocusStyle(e)
				continue
			}
			e.EchoMode = textinput.EchoPassword
			e.EchoCharacter = '•'
			f.elements[i] = e
		}
	}

	return f
}

func (f *Model) SetMainModel(m tea.Model) {
	f.main = m
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
	case message.OpenLoginMsg:
		for i := range f.elements {
			if e, ok := f.elements[i].(textinput.Model); ok {
				e.SetValue("")
				f.elements[i] = e
			}
		}
	case message.SuccessMsg:
		return f.main.Update(msg)
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
		case key.Matches(msg, common.Keys.Enter):
			f.err = nil

			switch {
			case f.pos == posSubmitBtn || f.pos == posEmail || f.pos == posPassword:
				return f.submit()
			case f.pos == posRegisterBtn:
				return f.fromRegister.Update(message.OpenRegisterMsg{
					LoginModel: f,
				})
			case f.pos == posCloseBtn:
				return f, tea.Quit
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
		errValidate.AddError("Неверно заполнен email")
	}

	if password.Value() == "" {
		f.elements[posPassword] = style.SetErrorStyle(password)
		errValidate.AddError("Не заполнен пароль")
	}

	if errValidate.HasErrors() {
		f.err = errValidate
		return f, nil
	}

	err := f.usecaseLogin.Do(context.TODO(), login.Data{
		Email:    email.Value(),
		Password: password.Value(),
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.Unauthenticated {
				f.err = errPasswordInvalid
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

	return f.main.Update(message.SuccessMsg{})
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
