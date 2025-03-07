package text

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/cli/element/list"
)

type OpenMsg struct {
	BackModel tea.Model
	BackState int
	Item      *list.Item
}
