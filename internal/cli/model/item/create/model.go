package create

import (
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
	defaultWidth = 40
	listHeight   = 14
)

type Model struct {
	main       tea.Model
	help       help.Model
	header     string
	categories list.Model
	err        error

	usecaseSync *sync.Usecase
	fetcherItem *item.Fetcher
}

func InitModel() *Model {
	f := &Model{
		help:   help.New(),
		header: "Создание",
		categories: elist.CreateDefaultList("Выберите категорию:", defaultWidth, listHeight, elist.CategoryDelegate{},
			elist.Category{Category: client.CategoryPassword, Title: client.CategoryPassword.String()},
			elist.Category{Category: client.CategoryText, Title: client.CategoryText.String()},
			elist.Category{Category: client.CategoryBlob, Title: client.CategoryBlob.String()},
			elist.Category{Category: client.CategoryBankCard, Title: client.CategoryBankCard.String()},
		),
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
	case Open:
		return f, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.Keys.Enter):
			if c, ok := f.categories.SelectedItem().(elist.Category); ok {
				return f.main.Update(common.OpenItemMessage{
					BackModel: f,
					Category:  c.Category,
				})
			}

			return f, nil
		case key.Matches(msg, common.Keys.Back):
			return f.main.Update(common.BackMessage{})
		}
	}

	var cmd tea.Cmd
	f.categories, cmd = f.categories.Update(msg)

	return f, cmd
}

func (f *Model) View() string {
	var b strings.Builder

	b.WriteString(style.TitleStyle.Render(f.header))
	b.WriteRune('\n')

	b.WriteString(f.categories.View())

	// выводим прочие ошибки
	if f.err != nil {
		b.WriteRune('\n')
		b.WriteString(style.ErrorBlockStyle.Render(f.err.Error()))
	}

	return b.String()
}
