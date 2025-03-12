package common

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/domain/client"
)

type BackMsg struct {
	State int
}

type OpenItemMessage struct {
	BackModel tea.Model
	BackState int
	Item      *client.Item
	Category  client.Category
}
