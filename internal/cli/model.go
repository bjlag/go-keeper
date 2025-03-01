package cli

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/cli/common"
	"github.com/bjlag/go-keeper/internal/cli/element"
	listf "github.com/bjlag/go-keeper/internal/cli/form/list"
	"github.com/bjlag/go-keeper/internal/cli/form/login"
	"github.com/bjlag/go-keeper/internal/cli/form/password"
	"github.com/bjlag/go-keeper/internal/cli/form/register"
	"github.com/bjlag/go-keeper/internal/cli/message"
)

type MainModel struct {
	help   help.Model
	header string

	formLogin    *login.Form
	formRegister *register.Form
	formList     *listf.Form
	formPassword *password.Form

	accessToken  string
	refreshToken string
}

func InitModel(opts ...Option) *MainModel {
	m := &MainModel{
		help:   help.New(),
		header: "Go Keeper",
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *MainModel) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			return message.SuccessLoginMessage{}
			//return message.OpenLoginFormMessage{}
		},
	)
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.Keys.Quit):
			return m, tea.Quit
		}

	// Forms
	case message.OpenLoginFormMessage:
		return m.formLogin.Update(tea.ClearScreen())
	case message.OpenRegisterFormMessage:
		return m.formRegister.Update(tea.ClearScreen())
	case message.OpenCategoryListFormMessage:
		return m.formList.Update(msg)
	case message.OpenPasswordListFormMessage:
		return m.formList.Update(msg)
	case message.OpenPasswordFormMessage:
		return m.formPassword.Update(msg)

	// Success
	case message.SuccessLoginMessage:
		m.accessToken = msg.AccessToken
		m.refreshToken = msg.RefreshToken

		return m.Update(message.OpenCategoryListFormMessage{})
	case message.SuccessRegisterMessage:
		m.accessToken = msg.AccessToken
		m.refreshToken = msg.RefreshToken

		return m.Update(message.OpenCategoryListFormMessage{})
	}

	return m, nil
}

func (m *MainModel) View() string {
	var b strings.Builder

	b.WriteString(element.TitleStyle.Render(m.header))
	b.WriteRune('\n')

	return b.String()
}
