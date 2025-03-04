package list

import (
	"context"
	"strings"

	"github.com/bjlag/go-keeper/internal/cli/common"
	"github.com/bjlag/go-keeper/internal/cli/element"
	"github.com/bjlag/go-keeper/internal/cli/model/item"
	"github.com/bjlag/go-keeper/internal/cli/style"
	"github.com/bjlag/go-keeper/internal/usecase/client/sync"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	stateCategoryList int = iota
	statePasswordList
)

const (
	defaultWidth = 40
	listHeight   = 14
)

type Model struct {
	main       tea.Model
	help       help.Model
	state      int
	header     string
	categories list.Model
	passwords  list.Model
	err        error

	usecase *sync.Usecase
}

func InitModel(usecase *sync.Usecase) *Model {
	f := &Model{
		help:       help.New(),
		header:     "Категории",
		categories: element.CreateDefaultList("Категории:", defaultWidth, listHeight),
		passwords:  element.CreateDefaultList("Пароли:", defaultWidth, listHeight),

		usecase: usecase,
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
		f.categories.SetWidth(msg.Width)
		f.passwords.SetWidth(msg.Width)
		return f, nil
	case common.BackMessage:
		switch msg.State {
		case stateCategoryList:
			return f.Update(OpenCategoryListMessage{})
		case statePasswordList:
			return f.Update(OpenPasswordListMessage{})
		}

	case GetAllDataMessage:
		f.state = stateCategoryList
		f.err = f.usecase.Do(context.TODO())

		return f.Update(OpenCategoryListMessage{})
	case OpenCategoryListMessage:
		f.state = stateCategoryList

		// todo получаем данные из базы

		f.categories.SetItems(nil)
		f.categories.InsertItem(len(f.categories.Items()), element.Item{Name: "Логины"})
		f.categories.InsertItem(len(f.categories.Items()), element.Item{Name: "Тексты"})
		f.categories.InsertItem(len(f.categories.Items()), element.Item{Name: "Файлы"})
		f.categories.InsertItem(len(f.categories.Items()), element.Item{Name: "Банковские карты"})

		return f, nil
	case OpenPasswordListMessage:
		f.state = statePasswordList

		// todo получаем данные из базы

		f.passwords.SetItems(nil)
		f.passwords.Title = f.categories.SelectedItem().(element.Item).Name + ":"
		f.passwords.InsertItem(len(f.categories.Items()), element.Item{Name: "Пароль 1"})
		f.passwords.InsertItem(len(f.categories.Items()), element.Item{Name: "Пароль 2"})
		f.passwords.InsertItem(len(f.categories.Items()), element.Item{Name: "Пароль 3"})
		f.passwords.InsertItem(len(f.categories.Items()), element.Item{Name: "Пароль 4"})

		return f, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.Keys.Quit):
			return f, tea.Quit
		case key.Matches(msg, common.Keys.Enter):
			switch f.state {
			case stateCategoryList:
				if i, ok := f.categories.SelectedItem().(element.Item); ok {
					return f.Update(OpenPasswordListMessage{
						Category: i,
					})
				}
			case statePasswordList:
				if i, ok := f.passwords.SelectedItem().(element.Item); ok {
					return f.main.Update(item.OpenMessage{
						BackModel: f,
						BackState: f.state,
						Item:      i,
					})
				}
			}

			return f, nil
		case key.Matches(msg, common.Keys.Back):
			switch f.state {
			case statePasswordList:
				return f.Update(OpenCategoryListMessage{})
			}
		}
	}

	var cmd tea.Cmd
	switch f.state {
	case stateCategoryList:
		f.categories, cmd = f.categories.Update(msg)
	case statePasswordList:
		f.passwords, cmd = f.passwords.Update(msg)
	}

	return f, cmd
}

func (f *Model) View() string {
	var b strings.Builder

	b.WriteString(style.TitleStyle.Render(f.header))
	b.WriteRune('\n')

	switch f.state {
	case stateCategoryList:
		b.WriteString(f.categories.View())
	case statePasswordList:
		b.WriteString(f.passwords.View())
	}

	//b.WriteRune('\n')
	//b.WriteString(f.help.View(common.Keys))

	// выводим прочие ошибки
	if f.err != nil {
		b.WriteRune('\n')
		b.WriteString(style.ErrorBlockStyle.Render(f.err.Error()))
	}

	return b.String()
}
