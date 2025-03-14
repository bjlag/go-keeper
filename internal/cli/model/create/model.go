package create

import (
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
)

var errUnsupportedCategory = errors.New("unsupported category")

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

	formPassword tea.Model
	formText     tea.Model
	formBankCard tea.Model
	formFile     tea.Model
}

func InitModel(
	formPassword tea.Model,
	formText tea.Model,
	formBankCard tea.Model,
	formFile tea.Model,
) *Model {
	return &Model{
		help:   help.New(),
		header: "Создание",
		categories: elist.CreateDefaultList("Выберите категорию:", defaultWidth, listHeight, elist.CategoryDelegate{},
			elist.Category{Category: client.CategoryPassword, Title: client.CategoryPassword.String()},
			elist.Category{Category: client.CategoryText, Title: client.CategoryText.String()},
			elist.Category{Category: client.CategoryFile, Title: client.CategoryFile.String()},
			elist.Category{Category: client.CategoryBankCard, Title: client.CategoryBankCard.String()},
		),

		formPassword: formPassword,
		formText:     formText,
		formBankCard: formBankCard,
		formFile:     formFile,
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
	case message.OpenItemMsg:
		return f, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.Keys.Enter):
			if c, ok := f.categories.SelectedItem().(elist.Category); ok {
				itemMsg := message.OpenItemMsg{
					BackModel: f,
				}

				switch c.Category {
				case client.CategoryPassword:
					return f.formPassword.Update(itemMsg)
				case client.CategoryText:
					return f.formText.Update(itemMsg)
				case client.CategoryFile:
					return f.formFile.Update(itemMsg)
				case client.CategoryBankCard:
					return f.formBankCard.Update(itemMsg)
				default:
					f.err = errUnsupportedCategory
					return f, nil
				}
			}

			return f, nil
		case key.Matches(msg, common.Keys.Back):
			return f.main.Update(message.BackMsg{})
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
