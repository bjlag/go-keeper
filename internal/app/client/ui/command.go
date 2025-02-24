package ui

import (
	"errors"
	"github.com/bjlag/go-keeper/internal/app/client/ui/style"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/infrastructure/validator"
)

const (
	EmailInput = iota
	PasswordInput
)

type logonCmd struct {
}

func submitLogin(m *Model) tea.Cmd {
	return func() tea.Msg {
		isValid := true
		for i := range m.inputs {
			switch i {
			case EmailInput:
				if !validator.ValidateEmail(m.inputs[i].Value()) {
					m.inputs[i].PromptStyle = style.ErrorStyle
					m.inputs[i].TextStyle = style.ErrorStyle
					m.inputs[i].Err = errors.New("Неверный емейл")
					isValid = false
				}
			case PasswordInput:
				if m.inputs[i].Value() == "" {
					m.inputs[i].PromptStyle = style.ErrorStyle
					m.inputs[i].TextStyle = style.ErrorStyle
					m.inputs[i].Err = errors.New("Не указан пароль")
					isValid = false
				}
			}
		}

		if !isValid {
			return nil
		}

		return logonCmd{}
	}
}

func (m logonCmd) Login() {
	//fmt.Println("login")
}
