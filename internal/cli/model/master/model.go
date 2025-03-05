package master

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/cli/common"
	"github.com/bjlag/go-keeper/internal/cli/model/item/password"
	listf "github.com/bjlag/go-keeper/internal/cli/model/list"
	"github.com/bjlag/go-keeper/internal/cli/model/login"
	"github.com/bjlag/go-keeper/internal/cli/model/register"
	"github.com/bjlag/go-keeper/internal/cli/style"
	"github.com/bjlag/go-keeper/internal/infrastructure/store/client/token"
)

type Model struct {
	help   help.Model
	header string

	formLogin    *login.Model
	formRegister *register.Model
	formList     *listf.Model
	formPassword *password.Model

	storeTokens *token.Store
}

func InitModel(opts ...Option) *Model {
	m := &Model{
		help:   help.New(),
		header: "Go Keeper",
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			//return message.SuccessLoginMessage{}
			return login.OpenMessage{}
		},
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.Keys.Quit):
			return m, tea.Quit
		}

	// Forms
	case login.OpenMessage:
		return m.formLogin.Update(msg)
	case register.OpenMessage:
		return m.formRegister.Update(msg)
	case listf.OpenCategoryListMessage:
		return m.formList.Update(msg)
	case listf.OpenItemListMessage:
		return m.formList.Update(msg)
	case password.OpenMessage:
		return m.formPassword.Update(msg)

	// Success
	case login.SuccessMessage:
		m.storeTokens.SaveTokens(msg.AccessToken, msg.RefreshToken)
		return m.formList.Update(listf.GetAllDataMessage{})
	case register.SuccessMessage:
		m.storeTokens.SaveTokens(msg.AccessToken, msg.RefreshToken)
		return m.formList.Update(listf.GetAllDataMessage{})
	}

	return m, nil
}

func (m *Model) View() string {
	var b strings.Builder

	b.WriteString(style.TitleStyle.Render(m.header))
	b.WriteRune('\n')

	return b.String()
}
