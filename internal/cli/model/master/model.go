package master

import (
	"errors"
	bank_card2 "github.com/bjlag/go-keeper/internal/cli/message/item/bank_card"
	create2 "github.com/bjlag/go-keeper/internal/cli/message/item/create"
	file2 "github.com/bjlag/go-keeper/internal/cli/message/item/file"
	password2 "github.com/bjlag/go-keeper/internal/cli/message/item/password"
	text2 "github.com/bjlag/go-keeper/internal/cli/message/item/text"
	"github.com/bjlag/go-keeper/internal/cli/message/list"
	login2 "github.com/bjlag/go-keeper/internal/cli/message/login"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/cli/common"
	"github.com/bjlag/go-keeper/internal/cli/element/button"
	"github.com/bjlag/go-keeper/internal/cli/model/item/bank_card"
	"github.com/bjlag/go-keeper/internal/cli/model/item/create"
	"github.com/bjlag/go-keeper/internal/cli/model/item/file"
	"github.com/bjlag/go-keeper/internal/cli/model/item/password"
	"github.com/bjlag/go-keeper/internal/cli/model/item/text"
	listf "github.com/bjlag/go-keeper/internal/cli/model/list"
	"github.com/bjlag/go-keeper/internal/cli/model/login"
	"github.com/bjlag/go-keeper/internal/cli/model/register"
	"github.com/bjlag/go-keeper/internal/cli/style"
	"github.com/bjlag/go-keeper/internal/domain/client"
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
	formBankCard *bank_card.Model
	formFile     *file.Model
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
			return login2.OpenMsg{}
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
				return m.formList.Update(list.GetDataMsg{})
			case posCreateBtn:
				return m.formCreate.Update(create2.Open{})
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
				if e, ok := m.elements[i].(button.Button); ok {
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
		}

	case OpenMsg:
		return m, nil
	case common.BackMsg:
		return m.Update(OpenMsg{})

	// Forms
	case login2.OpenMsg:
		return m.formLogin.Update(msg)
	case register.OpenMessage:
		return m.formRegister.Update(msg)
	case list.OpenCategoriesMsg:
		return m.formList.Update(msg)
	case list.OpenItemsMsg:
		return m.formList.Update(msg)
	case common.OpenItemMessage:
		switch msg.Category {
		case client.CategoryPassword:
			return m.formPassword.Update(password2.OpenMsg{
				BackModel: msg.BackModel,
				BackState: msg.BackState,
				Item:      msg.Item,
			})
		case client.CategoryText:
			return m.formText.Update(text2.OpenMsg{
				BackModel: msg.BackModel,
				BackState: msg.BackState,
				Item:      msg.Item,
			})
		case client.CategoryBankCard:
			return m.formBankCard.Update(bank_card2.OpenMsg{
				BackModel: msg.BackModel,
				BackState: msg.BackState,
				Item:      msg.Item,
			})
		case client.CategoryFile:
			return m.formFile.Update(file2.OpenMsg{
				BackModel: msg.BackModel,
				BackState: msg.BackState,
				Item:      msg.Item,
			})
		default:
			m.err = errUnsupportedCategory
		}

	// Success
	case login2.SuccessMsg:
		return m.Update(OpenMsg{})
	case register.SuccessMsg:
		return m.Update(OpenMsg{})
	}

	return m.Update(OpenMsg{})
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
