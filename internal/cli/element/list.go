package element

import (
	"fmt"
	"io"
	"strings"

	"github.com/bjlag/go-keeper/internal/cli/style"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func CreateDefaultList(title string, with, height int, items ...list.Item) list.Model {
	l := list.New(items, ItemDelegate{}, with, height)

	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	//l.SetShowTitle(false)
	l.Styles.Title = style.ListTitleStyle
	l.Styles.PaginationStyle = style.ListPaginationStyle
	l.Styles.HelpStyle = style.ListHelpStyle

	return l
}

type Item struct {
	ID   string
	Name string
}

func (i Item) FilterValue() string { return "" }

type ItemDelegate struct{}

func (d ItemDelegate) Height() int {
	return 1
}

func (d ItemDelegate) Spacing() int {
	return 0
}

func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Name)

	fn := style.ListItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return style.SelectedListItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprint(w, fn(str))
}
