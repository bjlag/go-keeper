package master

import (
	"errors"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/cli/common"
	"github.com/bjlag/go-keeper/internal/cli/element/button"
	"github.com/bjlag/go-keeper/internal/cli/model/item/create"
	"github.com/bjlag/go-keeper/internal/cli/model/item/password"
	"github.com/bjlag/go-keeper/internal/cli/model/item/text"
	listf "github.com/bjlag/go-keeper/internal/cli/model/list"
	"github.com/bjlag/go-keeper/internal/cli/model/login"
	"github.com/bjlag/go-keeper/internal/cli/model/register"
	"github.com/bjlag/go-keeper/internal/cli/style"
	"github.com/bjlag/go-keeper/internal/domain/client"
	"github.com/bjlag/go-keeper/internal/infrastructure/store/client/token"
)

var errUnsupportedCategory = errors.New("unsupported category")

const (
	posViewBtn int = iota
	posCreateBtn
	posCloseBtn
)

type Model struct {
	help     help.Model
	header   string
	elements []interface{}
	pos      int
	err      error

	formLogin    *login.Model
	formRegister *register.Model
	formCreate   *create.Model
	formList     *listf.Model
	formPassword *password.Model
	formText     *text.Model

	storeTokens *token.Store
}

func InitModel(opts ...Option) *Model {
	m := &Model{
		help:   help.New(),
		header: "Go Keeper",
		elements: []interface{}{
			posViewBtn:   button.CreateDefaultButton("Просмотр", button.WithFocused()),
			posCreateBtn: button.CreateDefaultButton("Создать"),
			posCloseBtn:  button.CreateDefaultButton("Выйти"),
		},
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
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
		case key.Matches(msg, common.Keys.Enter):
			switch m.pos {
			case posViewBtn:
				return m.formList.Update(listf.GetAllDataMessage{})
			case posCreateBtn:
				return m.formCreate.Update(create.Open{})
			case posCloseBtn:
				return m, tea.Quit
			}
		case key.Matches(msg, common.Keys.Navigation):
			if key.Matches(msg, common.Keys.Down, common.Keys.Tab) {
				m.pos++
			} else {
				m.pos--
			}

			if m.pos > len(m.elements)-1 {
				m.pos = 0
			} else if m.pos < 0 {
				m.pos = len(m.elements) - 1
			}

			for i := range m.elements {
				switch e := m.elements[i].(type) {
				case button.Button:
					if i == m.pos {
						e.Focus()
						m.elements[i] = e
						continue
					}
					e.Blur()
					m.elements[i] = e
				}
			}

			return m, nil
		case key.Matches(msg, common.Keys.Quit):
			return m, tea.Quit
		}

	case Open:
		return m, nil
	case common.BackMessage:
		return m.Update(Open{})

	// Forms
	case login.OpenMessage:
		return m.formLogin.Update(msg)
	case register.OpenMessage:
		return m.formRegister.Update(msg)
	case listf.OpenCategoryListMessage:
		return m.formList.Update(msg)
	case listf.OpenItemListMessage:
		return m.formList.Update(msg)
	case common.OpenItemMessage:
		switch msg.Category {
		case client.CategoryPassword:
			return m.formPassword.Update(password.OpenMsg{
				BackModel: msg.BackModel,
				BackState: msg.BackState,
				Item:      msg.Item,
			})
		case client.CategoryText:
			return m.formText.Update(text.OpenMsg{
				BackModel: msg.BackModel,
				BackState: msg.BackState,
				Item:      msg.Item,
			})
		default:
			m.err = errUnsupportedCategory
		}

	// Success
	case login.SuccessMessage:
		m.storeTokens.SaveTokens(msg.AccessToken, msg.RefreshToken)
	case register.SuccessMessage:
		m.storeTokens.SaveTokens(msg.AccessToken, msg.RefreshToken)
	}

	return m.Update(Open{})
}

func (m *Model) View() string {
	var b strings.Builder

	b.WriteString(style.TitleStyle.Render(m.header))
	b.WriteRune('\n')

	for i := range m.elements {
		if e, ok := m.elements[i].(button.Button); ok {
			b.WriteString(e.String())
			b.WriteRune('\n')
		}
	}

	// выводим прочие ошибки
	if m.err != nil {
		b.WriteRune('\n')
		b.WriteString(style.ErrorBlockStyle.Render(m.err.Error()))
	}

	return b.String()
}
