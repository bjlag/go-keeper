package password

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/cli/element/list"
)

type OpenMessage struct {
	BackModel tea.Model
	BackState int
	Item      list.Item
}
