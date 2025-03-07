package list

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/cli/common"
	elist "github.com/bjlag/go-keeper/internal/cli/element/list"
	"github.com/bjlag/go-keeper/internal/cli/style"
	"github.com/bjlag/go-keeper/internal/domain/client"
	"github.com/bjlag/go-keeper/internal/fetcher/item"
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

	selectedCategory client.Category

	usecaseSync *sync.Usecase
	fetcherItem *item.Fetcher
}

func InitModel(usecaseSync *sync.Usecase, fetcherItem *item.Fetcher) *Model {
	f := &Model{
		help:       help.New(),
		header:     "Категории",
		categories: elist.CreateDefaultList("Категории:", defaultWidth, listHeight, elist.CategoryDelegate{}),
		items:      elist.CreateDefaultList("Пароли:", defaultWidth, listHeight, elist.ItemDelegate{}),

		usecaseSync: usecaseSync,
		fetcherItem: fetcherItem,
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
		f.categories.InsertItem(len(f.categories.Items()), elist.Category{Category: client.CategoryPassword, Title: client.CategoryPassword.String()})
		f.categories.InsertItem(len(f.categories.Items()), elist.Category{Category: client.CategoryText, Title: client.CategoryText.String()})
		f.categories.InsertItem(len(f.categories.Items()), elist.Category{Category: client.CategoryBlob, Title: client.CategoryBlob.String()})
		f.categories.InsertItem(len(f.categories.Items()), elist.Category{Category: client.CategoryBankCard, Title: client.CategoryBankCard.String()})

		return f, nil
	case OpenItemListMessage:
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
		for _, p := range items {
			f.items.InsertItem(len(f.categories.Items()), elist.Item{
				GUID:     p.GUID,
				Category: p.Category,
				Title:    p.Title,
				Value:    p.Value,
				Notes:    p.Notes,
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
					return f.Update(OpenItemListMessage{
						Category: c.Category,
					})
				}
			case stateItemList:
				if i, ok := f.items.SelectedItem().(elist.Item); ok {
					f.selectedCategory = i.Category

					return f.main.Update(common.OpenItemMessage{
						BackModel: f,
						BackState: f.state,
						Item:      i,
					})
				}
			}

			return f, nil
		case key.Matches(msg, common.Keys.Back):
			switch f.state {
			case stateCategoryList:
				return f.main.Update(common.BackMessage{})
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

	// выводим прочие ошибки
	if f.err != nil {
		b.WriteRune('\n')
		b.WriteString(style.ErrorBlockStyle.Render(f.err.Error()))
	}

	return b.String()
}
