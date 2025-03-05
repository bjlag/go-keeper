package list

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/cli/common"
	"github.com/bjlag/go-keeper/internal/cli/element"
	itemModel "github.com/bjlag/go-keeper/internal/cli/model/item"
	"github.com/bjlag/go-keeper/internal/cli/style"
	"github.com/bjlag/go-keeper/internal/domain/client"
	"github.com/bjlag/go-keeper/internal/usecase/client/item"
	"github.com/bjlag/go-keeper/internal/usecase/client/sync"
)

const (
	stateCategoryList int = iota
	stateItemList
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
	items      list.Model
	err        error

	usecaseSync *sync.Usecase
	usecaseItem *item.Usecase
}

func InitModel(usecaseSync *sync.Usecase, usecaseItem *item.Usecase) *Model {
	f := &Model{
		help:       help.New(),
		header:     "Категории",
		categories: element.CreateDefaultList("Категории:", defaultWidth, listHeight, element.CategoryDelegate{}),
		items:      element.CreateDefaultList("Пароли:", defaultWidth, listHeight, element.ItemDelegate{}),

		usecaseSync: usecaseSync,
		usecaseItem: usecaseItem,
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
		f.items.SetWidth(msg.Width)
		return f, nil
	case common.BackMessage:
		switch msg.State {
		case stateCategoryList:
			return f.Update(OpenCategoryListMessage{})
		case stateItemList:
			return f.Update(OpenItemListMessage{})
		}

	case GetAllDataMessage:
		f.state = stateCategoryList
		f.err = f.usecaseSync.Do(context.TODO())

		return f.Update(OpenCategoryListMessage{})
	case OpenCategoryListMessage:
		f.state = stateCategoryList

		f.categories.SetItems(nil)
		f.categories.InsertItem(len(f.categories.Items()), element.Category{ID: client.CategoryPassword, Name: "Пароли"})
		f.categories.InsertItem(len(f.categories.Items()), element.Category{ID: client.CategoryText, Name: "Тексты"})
		f.categories.InsertItem(len(f.categories.Items()), element.Category{ID: client.CategoryBlob, Name: "Файлы"})
		f.categories.InsertItem(len(f.categories.Items()), element.Category{ID: client.CategoryBankCard, Name: "Банковские карты"})

		return f, nil
	case OpenItemListMessage:
		f.state = stateItemList

		items, err := f.usecaseItem.ItemsByCategory(context.TODO(), msg.Category)
		if err != nil {
			f.err = err
			return f, nil
		}

		f.items.SetItems(nil)
		f.items.Title = f.categories.SelectedItem().(element.Category).Name + ":"

		for _, p := range items {
			f.items.InsertItem(len(f.categories.Items()), element.Item{Name: p.Title})
		}

		return f, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.Keys.Quit):
			return f, tea.Quit
		case key.Matches(msg, common.Keys.Enter):
			switch f.state {
			case stateCategoryList:
				if c, ok := f.categories.SelectedItem().(element.Category); ok {
					return f.Update(OpenItemListMessage{
						Category: c.ID,
					})
				}
			case stateItemList:
				if i, ok := f.items.SelectedItem().(element.Item); ok {
					return f.main.Update(itemModel.OpenMessage{
						BackModel: f,
						BackState: f.state,
						Item:      i,
					})
				}
			}

			return f, nil
		case key.Matches(msg, common.Keys.Back):
			switch f.state {
			case stateItemList:
				return f.Update(OpenCategoryListMessage{})
			}
		}
	}

	var cmd tea.Cmd
	switch f.state {
	case stateCategoryList:
		f.categories, cmd = f.categories.Update(msg)
	case stateItemList:
		f.items, cmd = f.items.Update(msg)
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
	case stateItemList:
		b.WriteString(f.items.View())
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
