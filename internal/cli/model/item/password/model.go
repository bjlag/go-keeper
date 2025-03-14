package password

import (
	"errors"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"

	"github.com/bjlag/go-keeper/internal/cli/common"
	"github.com/bjlag/go-keeper/internal/cli/element/button"
	tarea "github.com/bjlag/go-keeper/internal/cli/element/textarea"
	tinput "github.com/bjlag/go-keeper/internal/cli/element/textinput"
	message "github.com/bjlag/go-keeper/internal/cli/message/item/password"
	"github.com/bjlag/go-keeper/internal/cli/message/item/sync"
	modelSync "github.com/bjlag/go-keeper/internal/cli/model/item/sync"
	"github.com/bjlag/go-keeper/internal/cli/style"
	"github.com/bjlag/go-keeper/internal/domain/client"
	"github.com/bjlag/go-keeper/internal/usecase/client/item/create"
	"github.com/bjlag/go-keeper/internal/usecase/client/item/edit"
	"github.com/bjlag/go-keeper/internal/usecase/client/item/remove"
)

const (
	posEditTitle int = iota
	posEditLogin
	posEditPassword
	posEditNotes
	posEditEditBtn
	posEditDeleteBtn
	posEditBackBtn
)

const (
	posCreateTitle int = iota
	posCreateLogin
	posCreatePassword
	posCreateNotes
	posCreateSaveBtn
	posCreateBackBtn
)

type state int

const (
	stateCreate state = iota
	stateEdit
)

var (
	errUnsupportedCommand   = errors.New("unsupported command")
	errInvalidValuePassword = errors.New("invalid value password")
)

type Model struct {
	main     tea.Model
	help     help.Model
	header   string
	state    state
	elements []interface{}
	pos      int
	err      error

	backModel tea.Model
	backState int

	guid     uuid.UUID
	item     *client.Item
	category client.Category

	formSync *modelSync.Model

	usecaseCreate *create.Usecase
	usecaseEdit   *edit.Usecase
	usecaseDelete *remove.Usecase
}

func InitModel(usecaseCreate *create.Usecase, usecaseSave *edit.Usecase, usecaseDelete *remove.Usecase, formSync *modelSync.Model) *Model {
	return &Model{
		help:   help.New(),
		header: "Пароль",
		state:  stateCreate,

		formSync: formSync,

		usecaseCreate: usecaseCreate,
		usecaseEdit:   usecaseSave,
		usecaseDelete: usecaseDelete,
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
	case tea.WindowSizeMsg:
		for i := range m.elements {
			if e, ok := m.elements[i].(textinput.Model); ok {
				e.Width = msg.Width
			}
		}
		return m, nil
	case message.OpenMsg:
		m.backState = msg.BackState
		m.backModel = msg.BackModel

		if msg.Item != nil {
			m.state = stateEdit
			m.header = msg.Item.Title
			m.item = msg.Item
			m.guid = msg.Item.GUID
			m.category = msg.Item.Category

			value, ok := msg.Item.Value.(*client.Password)
			if !ok {
				m.err = errInvalidValuePassword
				return m, nil
			}

			m.elements = []interface{}{
				posEditTitle:     tinput.CreateDefaultTextInput("Название", tinput.WithValue(msg.Item.Title), tinput.WithFocused()),
				posEditLogin:     tinput.CreateDefaultTextInput("Логин", tinput.WithValue(value.Login)),
				posEditPassword:  tinput.CreateDefaultTextInput("Пароль", tinput.WithValue(value.Password)),
				posEditNotes:     tarea.CreateDefaultTextArea("Заметки", tarea.WithValue(msg.Item.Notes)),
				posEditEditBtn:   button.CreateDefaultButton("Изменить"),
				posEditDeleteBtn: button.CreateDefaultButton("Удалить"),
				posEditBackBtn:   button.CreateDefaultButton("Назад"),
			}

			return m, nil
		}

		m.state = stateCreate
		m.header = "Новый пароль"
		m.elements = []interface{}{
			posCreateTitle:    tinput.CreateDefaultTextInput("Название", tinput.WithFocused()),
			posCreateLogin:    tinput.CreateDefaultTextInput("Логин"),
			posCreatePassword: tinput.CreateDefaultTextInput("Пароль"),
			posCreateNotes:    tarea.CreateDefaultTextArea("Заметки"),
			posCreateSaveBtn:  button.CreateDefaultButton("Сохранить"),
			posCreateBackBtn:  button.CreateDefaultButton("Назад"),
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

			for i := range m.elements {
				switch e := m.elements[i].(type) {
				case textinput.Model:
					if i == m.pos {
						e.Focus()
						m.elements[i] = style.SetFocusStyle(e)
						continue
					}

					e.Blur()
					m.elements[i] = style.SetNoStyle(e)
				case textarea.Model:
					if i == m.pos {
						e.Focus()
						m.elements[i] = e
						continue
					}

					e.Blur()
					m.elements[i] = e
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
		case key.Matches(msg, common.Keys.Enter):
			m.err = nil

			if m.state == stateCreate {
				switch m.pos {
				case posCreateSaveBtn:
					m.err = m.createAction()
					return m, nil
				case posCreateBackBtn:
					return m.backModel.Update(common.BackMsg{
						State: m.backState,
					})
				default:
					m.err = errUnsupportedCommand
				}

				return m, nil
			}

			switch m.pos {
			case posEditEditBtn:
				err := m.editAction()
				if err != nil && errors.Is(err, edit.ErrConflict) {
					return m.formSync.Update(sync.OpenMsg{
						BackModel: m,
						Item:      m.item,
					})
				}

				m.err = err
				return m, nil
			case posEditDeleteBtn:
				m.err = m.deleteAction()
				return m, nil
			case posEditBackBtn:
				return m.backModel.Update(common.BackMsg{
					State: m.backState,
				})
			default:
				m.err = errUnsupportedCommand
			}

			return m, nil
		case key.Matches(msg, common.Keys.Back):
			return m.backModel.Update(common.BackMsg{
				State: m.backState,
			})
		}
	}

	return m, m.updateInputs(msg)
}

func (m *Model) View() string {
	var b strings.Builder

	b.WriteString(style.TitleStyle.Render(m.header))
	b.WriteRune('\n')

	b.WriteString("Категория: ")
	b.WriteString(m.category.String())
	b.WriteRune('\n')

	for i := range m.elements {
		switch e := m.elements[i].(type) {
		case textinput.Model:
			b.WriteString(e.Placeholder)
			b.WriteRune('\n')
			b.WriteString(e.View())
			b.WriteRune('\n')
			b.WriteRune('\n')
		case textarea.Model:
			b.WriteString(e.Placeholder)
			b.WriteRune('\n')
			b.WriteString(e.View())
			b.WriteRune('\n')
			b.WriteRune('\n')
		}
	}

	b.WriteRune('\n')

	for i := range m.elements {
		if e, ok := m.elements[i].(button.Button); ok {
			b.WriteString(e.String())
			b.WriteRune('\n')
		}
	}

	var (
		errValidate *common.ValidateError
		errForm     *common.FormError
	)

	// выводим ошибки валидации
	if m.err != nil && (errors.As(m.err, &errValidate) || errors.As(m.err, &errForm)) {
		b.WriteString(style.ErrorBlockStyle.Render(m.err.Error()))
		b.WriteRune('\n')
	}

	b.WriteRune('\n')
	b.WriteString(m.help.View(common.Keys))

	// выводим прочие ошибки
	if m.err != nil && !(errors.As(m.err, &errValidate) || errors.As(m.err, &errForm)) {
		b.WriteRune('\n')
		b.WriteString(style.ErrorBlockStyle.Render(m.err.Error()))
	}

	return b.String()
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.elements))

	for i := range m.elements {
		switch e := m.elements[i].(type) {
		case textinput.Model:
			m.elements[i], cmds[i] = e.Update(msg)
		case textarea.Model:
			m.elements[i], cmds[i] = e.Update(msg)
		}
	}

	return tea.Batch(cmds...)
}
