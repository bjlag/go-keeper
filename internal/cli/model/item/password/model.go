package password

import (
	"errors"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/cli/common"
	"github.com/bjlag/go-keeper/internal/cli/element"
	"github.com/bjlag/go-keeper/internal/cli/style"
	"github.com/bjlag/go-keeper/internal/domain/client"
	"github.com/charmbracelet/bubbles/textarea"

	tarea "github.com/bjlag/go-keeper/internal/cli/element/textarea"
	tinput "github.com/bjlag/go-keeper/internal/cli/element/textinput"
)

const (
	posTitle int = iota
	posLogin
	posPassword
	posNotes
	posEditBtn
	posDeleteBtn
	posBackBtn
)

type Model struct {
	main     tea.Model
	help     help.Model
	header   string
	elements []interface{}
	pos      int
	err      error

	backModel tea.Model
	backState int

	category string
}

func InitModel() *Model {
	f := &Model{
		help:   help.New(),
		header: "Регистрация",

		//usecase: usecase,
	}

	return f
}

func (f *Model) SetMainModel(m tea.Model) {
	f.main = m
}

func (f *Model) Init() tea.Cmd {
	return nil
}

func (f *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		for i := range f.elements {
			switch e := f.elements[i].(type) {
			case textinput.Model:
				e.Width = msg.Width
			}
		}
		return f, nil
	case OpenMessage:
		f.backState = msg.BackState
		f.backModel = msg.BackModel
		f.header = msg.Item.Title
		f.category = msg.Item.Category.String()

		value, ok := msg.Item.Value.(*client.Password)
		if !ok {
			f.err = errors.New("it is not password")
			return f, nil
		}

		f.elements = []interface{}{
			posTitle:     tinput.CreateDefaultTextInput("Название", tinput.WithValue(msg.Item.Title), tinput.WithFocused()),
			posLogin:     tinput.CreateDefaultTextInput("Логин", tinput.WithValue(value.Login)),
			posPassword:  tinput.CreateDefaultTextInput("Пароль", tinput.WithValue(value.Password)),
			posNotes:     tarea.CreateDefaultTextArea("Заметки", tarea.WithValue(msg.Item.Notes)),
			posEditBtn:   element.CreateDefaultButton("Изменить"),
			posDeleteBtn: element.CreateDefaultButton("Удалить"),
			posBackBtn:   element.CreateDefaultButton("Назад"),
		}

		return f, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.Keys.Quit):
			return f, tea.Quit
		case key.Matches(msg, common.Keys.Navigation):
			if key.Matches(msg, common.Keys.Down, common.Keys.Tab) {
				f.pos++
			} else {
				f.pos--
			}

			if f.pos > len(f.elements)-1 {
				f.pos = 0
			} else if f.pos < 0 {
				f.pos = len(f.elements) - 1
			}

			for i := range f.elements {
				switch e := f.elements[i].(type) {
				case textinput.Model:
					if i == f.pos {
						e.Focus()
						f.elements[i] = style.SetFocusStyle(e)
						continue
					}

					e.Blur()
					f.elements[i] = style.SetNoStyle(e)
				case textarea.Model:
					if i == f.pos {
						e.Focus()
						f.elements[i] = e
						continue
					}

					e.Blur()
					f.elements[i] = e
				case element.Button:
					if i == f.pos {
						e.Focus()
						f.elements[i] = e
						continue
					}
					e.Blur()
					f.elements[i] = e
				}
			}

			return f, nil
		case key.Matches(msg, common.Keys.Enter):
			f.err = nil

			switch {
			case f.pos == posBackBtn:
				return f.backModel.Update(common.BackMessage{
					State: f.backState,
				})
			}

			return f, nil
		case key.Matches(msg, common.Keys.Back):
			return f.backModel.Update(common.BackMessage{
				State: f.backState,
			})
		}
	}

	return f, f.updateInputs(msg)
}

func (f *Model) View() string {
	var b strings.Builder

	b.WriteString(style.TitleStyle.Render(f.header))
	b.WriteRune('\n')

	b.WriteString("Категория: ")
	b.WriteString(f.category)
	b.WriteRune('\n')

	for i := range f.elements {
		switch e := f.elements[i].(type) {
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

	for i := range f.elements {
		if e, ok := f.elements[i].(element.Button); ok {
			b.WriteString(e.String())
			b.WriteRune('\n')
		}
	}

	var (
		errValidate *common.ValidateError
		errForm     *common.FormError
	)

	// выводим ошибки валидации
	if f.err != nil && (errors.As(f.err, &errValidate) || errors.As(f.err, &errForm)) {
		b.WriteString(style.ErrorBlockStyle.Render(f.err.Error()))
		b.WriteRune('\n')
	}

	b.WriteRune('\n')
	b.WriteString(f.help.View(common.Keys))

	// выводим прочие ошибки
	if f.err != nil && !(errors.As(f.err, &errValidate) || errors.As(f.err, &errForm)) {
		b.WriteRune('\n')
		b.WriteString(style.ErrorBlockStyle.Render(f.err.Error()))
	}

	return b.String()
}

func (f *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(f.elements))

	for i := range f.elements {
		switch m := f.elements[i].(type) {
		case textinput.Model:
			f.elements[i], cmds[i] = m.Update(msg)
		case textarea.Model:
			f.elements[i], cmds[i] = m.Update(msg)
		}
	}

	return tea.Batch(cmds...)
}
