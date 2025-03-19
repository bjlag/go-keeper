// Package list описывает модель для работы UI списка категорий и элементов.
package list

import (
	"context"
	"errors"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/cli/common"
	elist "github.com/bjlag/go-keeper/internal/cli/element/list"
	"github.com/bjlag/go-keeper/internal/cli/message"
	"github.com/bjlag/go-keeper/internal/cli/style"
	"github.com/bjlag/go-keeper/internal/domain/client"
	"github.com/bjlag/go-keeper/internal/fetcher/item"
	"github.com/bjlag/go-keeper/internal/usecase/client/sync"
)

var errUnsupportedCategory = errors.New("unsupported category")

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

	selectedCategory client.Category

	formPassword tea.Model
	formText     tea.Model
	formBankCard tea.Model
	formFile     tea.Model

	usecaseSync *sync.Usecase
	fetcherItem *item.Fetcher
}

func InitModel(
	usecaseSync *sync.Usecase,
	fetcherItem *item.Fetcher,
	formPassword tea.Model,
	formText tea.Model,
	formBankCard tea.Model,
	formFile tea.Model,
) *Model {
	f := &Model{
		help:       help.New(),
		header:     "Категории",
		categories: elist.CreateDefaultList("Категории:", defaultWidth, listHeight, elist.CategoryDelegate{}),
		items:      elist.CreateDefaultList("Пароли:", defaultWidth, listHeight, elist.ItemDelegate{}),

		usecaseSync: usecaseSync,
		fetcherItem: fetcherItem,

		formPassword: formPassword,
		formText:     formText,
		formBankCard: formBankCard,
		formFile:     formFile,
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
	case message.BackMsg:
		switch msg.State {
		case stateCategoryList:
			return f.Update(message.OpenCategoriesMsg{})
		case stateItemList:
			return f.Update(message.OpenItemsMsg{})
		}
	case message.OpenCategoriesMsg:
		f.state = stateCategoryList

		f.err = f.usecaseSync.Do(context.TODO())

		f.categories.SetItems(nil)
		f.categories.InsertItem(len(f.categories.Items()), elist.Category{Category: client.CategoryPassword, Title: client.CategoryPassword.String()})
		f.categories.InsertItem(len(f.categories.Items()), elist.Category{Category: client.CategoryText, Title: client.CategoryText.String()})
		f.categories.InsertItem(len(f.categories.Items()), elist.Category{Category: client.CategoryFile, Title: client.CategoryFile.String()})
		f.categories.InsertItem(len(f.categories.Items()), elist.Category{Category: client.CategoryBankCard, Title: client.CategoryBankCard.String()})

		return f, nil
	case message.OpenItemsMsg:
		f.state = stateItemList

		if c, ok := f.categories.SelectedItem().(elist.Category); ok {
			f.items.Title = c.Title + ":"
		}

		items, err := f.fetcherItem.ItemsByCategory(context.TODO(), f.selectedCategory)
		if err != nil {
			f.err = err
			return f, nil
		}

		f.items.SetItems(nil)
		for _, model := range items {
			f.items.InsertItem(len(f.categories.Items()), elist.Item{
				Model: model,
			})
		}

		return f, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.Keys.Quit):
			return f, tea.Quit
		case key.Matches(msg, common.Keys.Enter):
			switch f.state {
			case stateCategoryList:
				if c, ok := f.categories.SelectedItem().(elist.Category); ok {
					f.selectedCategory = c.Category
					return f.Update(message.OpenItemsMsg{
						Category: c.Category,
					})
				}
			case stateItemList:
				if i, ok := f.items.SelectedItem().(elist.Item); ok {
					f.selectedCategory = i.Model.Category

					switch i.Model.Category {
					case client.CategoryPassword:
						return f.formPassword.Update(message.OpenItemMsg{
							BackModel: f,
							BackState: f.state,
							Item:      &i.Model,
						})
					case client.CategoryText:
						return f.formText.Update(message.OpenItemMsg{
							BackModel: f,
							BackState: f.state,
							Item:      &i.Model,
						})
					case client.CategoryBankCard:
						return f.formBankCard.Update(message.OpenItemMsg{
							BackModel: f,
							BackState: f.state,
							Item:      &i.Model,
						})
					case client.CategoryFile:
						return f.formFile.Update(message.OpenItemMsg{
							BackModel: f,
							BackState: f.state,
							Item:      &i.Model,
						})
					default:
						f.err = errUnsupportedCategory
					}
				}
			}

			return f, nil
		case key.Matches(msg, common.Keys.Back):
			switch f.state {
			case stateCategoryList:
				return f.main.Update(message.BackMsg{})
			case stateItemList:
				return f.Update(message.OpenCategoriesMsg{})
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

	// выводим прочие ошибки
	if f.err != nil {
		b.WriteRune('\n')
		b.WriteString(style.ErrorBlockStyle.Render(f.err.Error()))
	}

	return b.String()
}
