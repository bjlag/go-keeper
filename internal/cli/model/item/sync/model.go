package sync

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/cli/common"
	"github.com/bjlag/go-keeper/internal/cli/element/button"
	"github.com/bjlag/go-keeper/internal/cli/message/item/bank_card"
	"github.com/bjlag/go-keeper/internal/cli/message/item/file"
	"github.com/bjlag/go-keeper/internal/cli/message/item/password"
	"github.com/bjlag/go-keeper/internal/cli/message/item/sync"
	"github.com/bjlag/go-keeper/internal/cli/message/item/text"
	"github.com/bjlag/go-keeper/internal/cli/style"
	"github.com/bjlag/go-keeper/internal/domain/client"
	itemSync "github.com/bjlag/go-keeper/internal/usecase/client/item/sync"
)

const (
	posSyncBtn int = iota
	posCancelBtn
)

type Model struct {
	main     tea.Model
	help     help.Model
	header   string
	elements []button.Button
	pos      int
	err      error

	item      client.Item
	prevModel tea.Model

	usecaseSync *itemSync.Usecase
}

func InitModel(usecaseSync *itemSync.Usecase) *Model {
	return &Model{
		help:        help.New(),
		usecaseSync: usecaseSync,
	}
}

func (m *Model) SetMainModel(model tea.Model) {
	m.main = model
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case sync.OpenMsg:
		m.prevModel = msg.BackModel

		if msg.Item == nil {
			return m.prevModel.Update(msg)
		}

		m.item = *msg.Item
		m.header = "Синхронизация"

		switch m.item.Category {
		case client.CategoryPassword:
			m.header += fmt.Sprintf(" пароля: %s", m.item.Title)
		case client.CategoryText:
			m.header += fmt.Sprintf(" текста: %s", m.item.Title)
		case client.CategoryFile:
			m.header += fmt.Sprintf(" файла: %s", m.item.Title)
		case client.CategoryBankCard:
			m.header += fmt.Sprintf(" банковской карты: %s", m.item.Title)
		}

		m.elements = []button.Button{
			posSyncBtn:   button.CreateDefaultButton("Синхронизировать"),
			posCancelBtn: button.CreateDefaultButton("Отмена"),
		}

		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.Keys.Quit):
			return m, tea.Quit
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

			for i, e := range m.elements {
				if i == m.pos {
					e.Focus()
					m.elements[i] = e
					continue
				}
				e.Blur()
				m.elements[i] = e
			}

			return m, nil
		case key.Matches(msg, common.Keys.Enter):
			m.err = nil

			switch m.pos {
			case posSyncBtn:
				err := m.syncAction()
				if err != nil {
					m.err = err
					return m, nil
				}

				switch m.item.Category {
				case client.CategoryPassword:
					return m.prevModel.Update(password.OpenMsg{
						Item: &m.item,
					})
				case client.CategoryText:
					return m.prevModel.Update(text.OpenMsg{
						Item: &m.item,
					})
				case client.CategoryFile:
					return m.prevModel.Update(file.OpenMsg{
						Item: &m.item,
					})
				case client.CategoryBankCard:
					return m.prevModel.Update(bank_card.OpenMsg{
						Item: &m.item,
					})
				}
			case posCancelBtn:
				switch m.item.Category {
				case client.CategoryPassword:
					return m.prevModel.Update(password.OpenMsg{
						Item: &m.item,
					})
				case client.CategoryText:
					return m.prevModel.Update(text.OpenMsg{
						Item: &m.item,
					})
				case client.CategoryFile:
					return m.prevModel.Update(file.OpenMsg{
						Item: &m.item,
					})
				case client.CategoryBankCard:
					return m.prevModel.Update(bank_card.OpenMsg{
						Item: &m.item,
					})
				}
			}

			return m, nil
		case key.Matches(msg, common.Keys.Back):
			return m.prevModel.Update(common.BackMsg{})
		}
	}

	return m, nil
}

func (m *Model) View() string {
	var b strings.Builder

	b.WriteString(style.TitleStyle.Render(m.header))
	b.WriteRune('\n')

	b.WriteString("Версия на сервере отличается от версии на клиенте.\n")
	b.WriteString("Синхронизируйте сначала данные с сервером, после меняйте запись.\n")
	b.WriteRune('\n')

	b.WriteString("Запись на сервере:\n")
	b.WriteRune('\n')

	b.WriteString("Название")
	b.WriteRune('\n')
	b.WriteString(m.item.Title)
	b.WriteRune('\n')
	b.WriteRune('\n')

	switch v := m.item.Value.(type) {
	case *client.Password:
		b.WriteString("Логин")
		b.WriteRune('\n')
		b.WriteString(v.Login)
		b.WriteRune('\n')
		b.WriteRune('\n')

		b.WriteString("Пароль")
		b.WriteRune('\n')
		b.WriteString(v.Password)
		b.WriteRune('\n')
		b.WriteRune('\n')
	case *client.File:
		b.WriteString("Название файла")
		b.WriteRune('\n')
		b.WriteString(v.Name)
		b.WriteRune('\n')
		b.WriteRune('\n')
	case *client.BankCard:
		b.WriteString("Номер карты")
		b.WriteRune('\n')
		b.WriteString(v.Number)
		b.WriteRune('\n')
		b.WriteRune('\n')

		b.WriteString("Истекает")
		b.WriteRune('\n')
		b.WriteString(v.Expiry)
		b.WriteRune('\n')
		b.WriteRune('\n')

		b.WriteString("CVV")
		b.WriteRune('\n')
		b.WriteString(v.CVV)
		b.WriteRune('\n')
		b.WriteRune('\n')
	}

	b.WriteString("Заметки")
	b.WriteRune('\n')
	b.WriteString(m.item.Notes)
	b.WriteRune('\n')
	b.WriteRune('\n')

	b.WriteRune('\n')

	for _, e := range m.elements {
		b.WriteString(e.String())
		b.WriteRune('\n')
	}

	// выводим прочие ошибки
	if m.err != nil {
		b.WriteRune('\n')
		b.WriteString(style.ErrorBlockStyle.Render(m.err.Error()))
	}

	return b.String()
}
