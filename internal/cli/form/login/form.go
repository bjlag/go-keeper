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
	"github.com/bjlag/go-keeper/internal/cli/element"
	"github.com/bjlag/go-keeper/internal/cli/message"
	"github.com/bjlag/go-keeper/internal/infrastructure/validator"
	"github.com/bjlag/go-keeper/internal/usecase/client/login"
)

const (
	posEmail int = iota
	posPassword
	posSubmitBtn
	posRegisterBtn
	posCloseBtn

	emailCharLimit    = 20
	passwordCharLimit = 20
)

var errPasswordInvalid = common.NewFormError("Неверный email или пароль")

type Form struct {
	main     tea.Model
	help     help.Model
	header   string
	elements []interface{}
	pos      int
	err      error

	usecase *login.Usecase
}

func NewForm(usecase *login.Usecase) *Form {
	f := &Form{
		help:   help.New(),
		header: "Авторизация",
		elements: []interface{}{
			posEmail:       element.CreateDefaultTextInput("Email", emailCharLimit),
			posPassword:    element.CreateDefaultTextInput("Password", passwordCharLimit),
			posSubmitBtn:   element.CreateDefaultButton("Вход"),
			posRegisterBtn: element.CreateDefaultButton("Регистрация"),
			posCloseBtn:    element.CreateDefaultButton("Закрыть"),
		},

		usecase: usecase,
	}

	for i := range f.elements {
		if e, ok := f.elements[i].(textinput.Model); ok {
			if i == posEmail {
				e.TextStyle = element.FocusedStyle
				e.PromptStyle = element.FocusedStyle
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

func (f *Form) SetMainModel(m tea.Model) {
	f.main = m
}

func (f *Form) Init() tea.Cmd {
	return nil
}

func (f *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		for i := range f.elements {
			switch e := f.elements[i].(type) {
			case textinput.Model:
				e.Width = msg.Width
			}
		}
		return f, nil
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
						f.elements[i] = element.SetFocusStyle(e)
						continue
					}

					e.Blur()
					f.elements[i] = element.SetNoStyle(e)
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
			case f.pos == posRegisterBtn:
				return f.main.Update(message.OpenRegisterFormMessage{})
			case f.pos == posCloseBtn:
				return f, tea.Quit
			}

			return f, nil
		}
	}

	return f, f.updateInputs(msg)
}

func (f *Form) View() string {
	var b strings.Builder

	b.WriteString(element.TitleStyle.Render(f.header))
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

	var (
		errValidate *common.ValidateError
		errForm     *common.FormError
	)

	// выводим ошибки валидации
	if f.err != nil && (errors.As(f.err, &errValidate) || errors.As(f.err, &errForm)) {
		b.WriteString(element.ErrorBlockStyle.Render(f.err.Error()))
		b.WriteRune('\n')
	}

	b.WriteRune('\n')
	b.WriteString(f.help.View(common.Keys))

	// выводим прочие ошибки
	if f.err != nil && !(errors.As(f.err, &errValidate) || errors.As(f.err, &errForm)) {
		b.WriteRune('\n')
		b.WriteString(element.ErrorBlockStyle.Render(f.err.Error()))
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
		f.elements[posEmail] = element.SetErrorStyle(email)
		errValidate.AddError("Неверно заполнен email")
	}

	if password.Value() == "" {
		f.elements[posPassword] = element.SetErrorStyle(password)
		errValidate.AddError("Не заполнен пароль")
	}

	if errValidate.HasErrors() {
		f.err = errValidate
		return f, nil
	}

	result, err := f.usecase.Do(context.TODO(), login.Data{
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

	return f.main.Update(message.LoginSuccessMessage{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
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
