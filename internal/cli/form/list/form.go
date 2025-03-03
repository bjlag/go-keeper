package list

import (
	"context"
	"github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/cli/common"
	"github.com/bjlag/go-keeper/internal/cli/element"
	"github.com/bjlag/go-keeper/internal/cli/message"
)

type state int

const (
	stateShowCategoryList state = iota
	stateShowPasswordList
)

const (
	defaultWidth = 40
	listHeight   = 14
)

type Form struct {
	main       tea.Model
	help       help.Model
	state      state
	header     string
	categories list.Model
	passwords  list.Model
	err        error

	rpcClient *client.RPCClient
	//usecase *login.Usecase
}

func NewForm(rpcClient *client.RPCClient) *Form {
	f := &Form{
		help:       help.New(),
		header:     "Категории",
		categories: element.CreateDefaultList("Категории:", defaultWidth, listHeight),
		passwords:  element.CreateDefaultList("Пароли:", defaultWidth, listHeight),

		rpcClient: rpcClient,
		//usecase: usecase,
	}

	return f
}

func (f *Form) SetMainModel(m tea.Model) {
	f.main = m
}

func (f *Form) Init() tea.Cmd {
	return nil
}

func (f *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.categories.SetWidth(msg.Width)
		f.passwords.SetWidth(msg.Width)
		return f, nil
	case message.OpenCategoryListFormMessage:
		f.state = stateShowCategoryList

		// todo получаем данные из базы
		// todo получить все данные
		in := &client.GetAllDataIn{
			Limit:  10,
			Offset: 0,
		}
		out, err := f.rpcClient.GetAllData(context.TODO(), in)
		if err != nil {
			if s, ok := status.FromError(err); ok {
				if s.Code() == codes.PermissionDenied {
					return f.main.Update(message.OpenLoginFormMessage{})
				}
				return f, nil
			}

			f.err = err
			return f, nil
		}

		_ = out

		f.categories.SetItems(nil)
		f.categories.InsertItem(len(f.categories.Items()), element.Item{Name: "Логины"})
		f.categories.InsertItem(len(f.categories.Items()), element.Item{Name: "Тексты"})
		f.categories.InsertItem(len(f.categories.Items()), element.Item{Name: "Файлы"})
		f.categories.InsertItem(len(f.categories.Items()), element.Item{Name: "Банковские карты"})

		return f, nil
	case message.OpenPasswordListFormMessage:
		f.state = stateShowPasswordList

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
			case stateShowCategoryList:
				if i, ok := f.categories.SelectedItem().(element.Item); ok {
					return f.Update(message.OpenPasswordListFormMessage{
						Category: i,
					})
				}
			case stateShowPasswordList:
				if i, ok := f.passwords.SelectedItem().(element.Item); ok {
					return f.main.Update(message.OpenPasswordFormMessage{
						Item: i,
					})
				}
			}

			return f, nil
		case key.Matches(msg, common.Keys.Back):
			switch f.state {
			case stateShowPasswordList:
				return f.Update(message.OpenCategoryListFormMessage{})
			}
		}
	}

	var cmd tea.Cmd
	switch f.state {
	case stateShowCategoryList:
		f.categories, cmd = f.categories.Update(msg)
	case stateShowPasswordList:
		f.passwords, cmd = f.passwords.Update(msg)
	}

	return f, cmd
}

func (f *Form) View() string {
	var b strings.Builder

	b.WriteString(element.TitleStyle.Render(f.header))
	b.WriteRune('\n')

	switch f.state {
	case stateShowCategoryList:
		b.WriteString(f.categories.View())
	case stateShowPasswordList:
		b.WriteString(f.passwords.View())
	}

	//b.WriteRune('\n')
	//b.WriteString(f.help.View(common.Keys))

	// выводим прочие ошибки
	if f.err != nil {
		b.WriteRune('\n')
		b.WriteString(element.ErrorBlockStyle.Render(f.err.Error()))
	}

	return b.String()
}
