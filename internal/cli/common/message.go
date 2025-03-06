package common

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/cli/element"
)

type BackMessage struct {
	State int
}

type OpenItemMessage struct {
	BackModel tea.Model
	BackState int
	Item      element.Item
}
