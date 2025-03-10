package text

import (
	"context"
	"errors"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"

	"github.com/bjlag/go-keeper/internal/cli/common"
	"github.com/bjlag/go-keeper/internal/cli/element"
	"github.com/bjlag/go-keeper/internal/cli/element/button"
	tarea "github.com/bjlag/go-keeper/internal/cli/element/textarea"
	tinput "github.com/bjlag/go-keeper/internal/cli/element/textinput"
	"github.com/bjlag/go-keeper/internal/cli/style"
	"github.com/bjlag/go-keeper/internal/domain/client"
	"github.com/bjlag/go-keeper/internal/usecase/client/item/create"
	"github.com/bjlag/go-keeper/internal/usecase/client/item/delete"
	"github.com/bjlag/go-keeper/internal/usecase/client/item/save"
	"github.com/charmbracelet/bubbles/textarea"
)

const (
	posTitle int = iota
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

	guid     uuid.UUID
	category client.Category

	usecaseCreate *create.Usecase
	usecaseSave   *save.Usecase
	usecaseDelete *delete.Usecase
}

func InitModel(usecaseCreate *create.Usecase, usecaseSave *save.Usecase, usecaseDelete *delete.Usecase) *Model {
	return &Model{
		help:          help.New(),
		header:        "Регистрация",
		usecaseCreate: usecaseCreate,
		usecaseSave:   usecaseSave,
		usecaseDelete: usecaseDelete,
	}
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
		f.guid = msg.Item.GUID
		f.category = msg.Item.Category

		f.elements = []interface{}{
			posTitle:     tinput.CreateDefaultTextInput("Название", tinput.WithValue(msg.Item.Title), tinput.WithFocused()),
			posNotes:     tarea.CreateDefaultTextArea("Текст", tarea.WithValue(msg.Item.Notes)),
			posEditBtn:   button.CreateDefaultButton("Изменить"),
			posDeleteBtn: button.CreateDefaultButton("Удалить"),
			posBackBtn:   button.CreateDefaultButton("Назад"),
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
				case button.Button:
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

			switch f.pos {
			case posEditBtn:
				f.err = f.edit()
				return f, nil
			case posDeleteBtn:
				f.err = f.delete()
				return f, nil
			case posBackBtn:
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
	b.WriteString(f.category.String())
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
		if e, ok := f.elements[i].(button.Button); ok {
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

func (f *Model) edit() error {
	i := client.Item{
		GUID:     f.guid,
		Category: f.category,
		Title:    element.GetValue(f.elements, posTitle),
		Notes:    element.GetValue(f.elements, posNotes),
	}

	return f.usecaseSave.Do(context.TODO(), i)
}

func (f *Model) delete() error {
	return f.usecaseDelete.Do(context.TODO(), f.guid)
}
