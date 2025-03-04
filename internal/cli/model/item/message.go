package item

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/cli/element"
)

type OpenMessage struct {
	BackModel tea.Model
	BackState int
	Item      element.Item
}
