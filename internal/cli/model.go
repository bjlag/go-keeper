package cli

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/cli/common"
	"github.com/bjlag/go-keeper/internal/cli/form/login"
	"github.com/bjlag/go-keeper/internal/cli/message"
)

type MainModel struct {
	help   help.Model
	header string

	accessToken  string
	refreshToken string
}

func InitModel() *MainModel {
	return &MainModel{
		help:   help.New(),
		header: "Go Keeper",
	}
}

func (m *MainModel) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			return login.NewForm(m)
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
	case *login.Form:
		return msg.Update(msg)
	case message.LoginSuccessMessage:
		m.accessToken = msg.AccessToken
		m.refreshToken = msg.RefreshToken
	case message.RegisterSuccessMessage:
		m.accessToken = msg.AccessToken
		m.refreshToken = msg.RefreshToken
	}

	return m, nil
}

func (m *MainModel) View() string {
	var b strings.Builder

	b.WriteString(common.TitleStyle.Render(m.header))
	b.WriteRune('\n')

	b.WriteString(common.TextStyle.Render("Access token:", m.accessToken, "\n\n"))
	b.WriteString(common.TextStyle.Render("Refresh token:", m.accessToken, "\n\n"))

	return b.String()
}
