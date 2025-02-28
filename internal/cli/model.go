package cli

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/cli/common"
	"github.com/bjlag/go-keeper/internal/cli/element"
	"github.com/bjlag/go-keeper/internal/cli/form/login"
	"github.com/bjlag/go-keeper/internal/cli/form/register"
	"github.com/bjlag/go-keeper/internal/cli/message"
)

const defaultWidth = 20
const listHeight = 14

type MainModel struct {
	help   help.Model
	header string
	list   list.Model

	loginForm    *login.Form
	registerForm *register.Form

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
			return message.OpenLoginFormMessage{}
		},
	)
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.Keys.Quit):
			return m, tea.Quit
		}
	case message.OpenLoginFormMessage:
		return m.loginForm.Update(tea.ClearScreen())
	case message.OpenRegisterFormMessage:
		return m.registerForm.Update(tea.ClearScreen())
	case message.LoginSuccessMessage:
		m.accessToken = msg.AccessToken
		m.refreshToken = msg.RefreshToken

		m.list = element.CreateDefaultList("Категории:", defaultWidth, listHeight,
			element.Item("Логины"),
			element.Item("Тексты"),
			element.Item("Файлы"),
			element.Item("Банковские карты"),
		)
	case message.RegisterSuccessMessage:
		m.accessToken = msg.AccessToken
		m.refreshToken = msg.RefreshToken

		m.list = element.CreateDefaultList("Категории:", defaultWidth, listHeight,
			element.Item("Логины"),
			element.Item("Тексты"),
			element.Item("Файлы"),
			element.Item("Банковские карты"),
		)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m *MainModel) View() string {
	var b strings.Builder

	b.WriteString(element.TitleStyle.Render(m.header))
	b.WriteRune('\n')

	b.WriteString(m.list.View())

	return b.String()
}
