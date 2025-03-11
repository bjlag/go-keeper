package common

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/cli/element/list"
	"github.com/bjlag/go-keeper/internal/domain/client"
)

type BackMessage struct {
	State int
}

type OpenItemMessage struct {
	BackModel tea.Model
	BackState int
	Item      *list.Item
	Category  client.Category
}
