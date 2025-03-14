package master

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/cli/common"
	"github.com/bjlag/go-keeper/internal/cli/element/button"
	"github.com/bjlag/go-keeper/internal/cli/message"
	"github.com/bjlag/go-keeper/internal/cli/model/create"
	listf "github.com/bjlag/go-keeper/internal/cli/model/list"
	"github.com/bjlag/go-keeper/internal/cli/model/login"
	"github.com/bjlag/go-keeper/internal/cli/style"
)

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

	formLogin  *login.Model
	formCreate *create.Model
	formList   *listf.Model
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
			return message.OpenLoginMsg{}
		},
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg == nil {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.Keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, common.Keys.Enter):
			switch m.pos {
			case posViewBtn:
				return m.formList.Update(message.OpenCategoriesMsg{})
			case posCreateBtn:
				return m.formCreate.Update(message.OpenItemMsg{})
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

	case message.BackMsg:
		return m.Update(nil)

	// Forms
	case message.OpenLoginMsg:
		return m.formLogin.Update(msg)
	case message.OpenCategoriesMsg:
		return m.formList.Update(msg)
	case message.OpenItemsMsg:
		return m.formList.Update(msg)

	// Success
	case message.SuccessMsg:
		return m.Update(nil)
	}

	return m.Update(nil)
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
