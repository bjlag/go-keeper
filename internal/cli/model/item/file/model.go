package file

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"

	"github.com/bjlag/go-keeper/internal/cli/common"
	"github.com/bjlag/go-keeper/internal/cli/element/button"
	"github.com/bjlag/go-keeper/internal/cli/element/list"
	tarea "github.com/bjlag/go-keeper/internal/cli/element/textarea"
	tinput "github.com/bjlag/go-keeper/internal/cli/element/textinput"
	"github.com/bjlag/go-keeper/internal/cli/style"
	"github.com/bjlag/go-keeper/internal/domain/client"
	"github.com/bjlag/go-keeper/internal/usecase/client/item/create"
	"github.com/bjlag/go-keeper/internal/usecase/client/item/edit"
	"github.com/bjlag/go-keeper/internal/usecase/client/item/remove"
)

const (
	posEditTitle int = iota
	posEditNotes
	posEditEditBtn
	posEditDeleteBtn
	posEditBackBtn
)

const (
	posCreateTitle int = iota
	posCreateNotes
	posCreateSelectFileBtn
	posCreateSaveBtn
	posCreateBackBtn
)

type state int

const (
	stateCreate state = iota
	stateEdit
)

var (
	errUnsupportedCommand = errors.New("unsupported command")
	errInvalidValue       = errors.New("invalid value")
)

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

type Model struct {
	main     tea.Model
	help     help.Model
	header   string
	state    state
	elements []interface{}
	pos      int
	err      error

	filepicker     filepicker.Model
	selectedFile   string
	selectFileMode bool
	fileData       []byte

	backModel tea.Model
	backState int
	item      *list.Item

	guid     uuid.UUID
	category client.Category

	usecaseCreate *create.Usecase
	usecaseEdit   *edit.Usecase
	usecaseDelete *remove.Usecase
}

func InitModel(usecaseCreate *create.Usecase, usecaseSave *edit.Usecase, usecaseDelete *remove.Usecase) *Model {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".txt", ".md", ".jpg", ".jpeg", ".png"}
	fp.CurrentDirectory, _ = os.UserHomeDir()
	fp.Height = 30

	return &Model{
		help:       help.New(),
		header:     "Файл",
		state:      stateCreate,
		filepicker: fp,

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
	case clearErrorMsg:
		m.err = nil
	case tea.WindowSizeMsg:
		for i := range m.elements {
			if e, ok := m.elements[i].(textinput.Model); ok {
				e.Width = msg.Width
			}
		}
		return m, nil
	case OpenMsg:
		m.backState = msg.BackState
		m.backModel = msg.BackModel

		if msg.Item != nil {
			m.state = stateEdit
			m.header = msg.Item.Title
			m.guid = msg.Item.GUID
			m.category = msg.Item.Category
			m.item = msg.Item

			value, ok := msg.Item.Value.(*client.File)
			if !ok {
				m.err = errInvalidValue
				return m, nil
			}

			m.selectedFile = value.Name
			m.fileData = value.Data

			m.elements = []interface{}{
				posEditTitle:     tinput.CreateDefaultTextInput("Название", tinput.WithValue(msg.Item.Title), tinput.WithFocused()),
				posEditNotes:     tarea.CreateDefaultTextArea("Заметки", tarea.WithValue(msg.Item.Notes)),
				posEditEditBtn:   button.CreateDefaultButton("Изменить"),
				posEditDeleteBtn: button.CreateDefaultButton("Удалить"),
				posEditBackBtn:   button.CreateDefaultButton("Назад"),
			}

			return m, nil
		}

		m.state = stateCreate
		m.header = "Новый файл"
		m.elements = []interface{}{
			posCreateTitle:         tinput.CreateDefaultTextInput("Название", tinput.WithFocused()),
			posCreateNotes:         tarea.CreateDefaultTextArea("Заметки"),
			posCreateSelectFileBtn: button.CreateDefaultButton("Выбрать файл"),
			posCreateSaveBtn:       button.CreateDefaultButton("Сохранить"),
			posCreateBackBtn:       button.CreateDefaultButton("Назад"),
		}

		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.Keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, common.Keys.Navigation):
			if !m.selectFileMode {
				m.navigate(msg)
				return m, nil
			}
		case key.Matches(msg, common.Keys.Enter):
			m.err = nil

			if m.selectFileMode {
				var cmd tea.Cmd
				m.filepicker, cmd = m.filepicker.Update(msg)

				if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
					m.selectedFile = path
					m.selectFileMode = false
				}

				if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
					m.err = errors.New(path + " не поддерживается")
					m.selectedFile = ""
					return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
				}

				return m, cmd
			}

			if m.state == stateCreate {
				switch m.pos {
				case posCreateSelectFileBtn:
					m.selectFileMode = true
					return m, m.filepicker.Init()
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
				m.err = m.editAction()
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
			if m.selectFileMode {
				m.selectFileMode = false
				return m.Update(OpenMsg{
					BackModel: m.backModel,
					BackState: m.backState,
					Item:      m.item,
				})
			}

			return m.backModel.Update(common.BackMsg{
				State: m.backState,
			})
		}
	}

	if m.selectFileMode {
		var cmd tea.Cmd
		m.filepicker, cmd = m.filepicker.Update(msg)
		return m, cmd
	}

	return m, m.updateInputs(msg)
}

func (m *Model) View() string {
	var b strings.Builder

	b.WriteString(style.TitleStyle.Render(m.header))
	b.WriteRune('\n')

	if m.selectFileMode {
		switch {
		case m.err != nil:
			b.WriteString(m.err.Error())
		case m.selectedFile == "":
			b.WriteString("Выберите файл:")
		default:
			b.WriteString("Выбранный файл: " + m.selectedFile)
		}

		b.WriteString("\n\n" + m.filepicker.View() + "\n")
		return b.String()
	}

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

	if m.selectedFile != "" {
		b.WriteString("Файл: " + m.selectedFile)
		b.WriteRune('\n')
		b.WriteRune('\n')
	}

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

func (m *Model) navigate(msg tea.KeyMsg) {
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
}
