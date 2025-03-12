package bank_card

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/domain/client"
)

type OpenMsg struct {
	BackModel tea.Model
	BackState int
	Item      *client.Item
}
