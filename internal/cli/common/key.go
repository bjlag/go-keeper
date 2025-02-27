package common

import "github.com/charmbracelet/bubbles/key"

type Map struct {
	New        key.Binding
	Edit       key.Binding
	Delete     key.Binding
	Up         key.Binding
	Down       key.Binding
	Right      key.Binding
	Left       key.Binding
	Enter      key.Binding
	Help       key.Binding
	Quit       key.Binding
	Back       key.Binding
	Tab        key.Binding
	Navigation key.Binding
}

func (k Map) ShortHelp() []key.Binding {
	return []key.Binding{k.Navigation, k.Enter, k.Quit}
}

func (k Map) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Navigation, k.Enter},
		{k.Help, k.Quit},
	}
}

var Keys = Map{
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "добавить"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "редактировать"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "удалить"),
	),
	Navigation: key.NewBinding(
		key.WithKeys("up", "down", "tab"),
		key.WithHelp("↑/↓/tab", "navigation"),
	),
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "наверх"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "вниз"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "табуляция"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "направо"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/l", "налево"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "выполнить"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "выход"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "назад"),
	),
}
