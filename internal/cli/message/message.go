package message

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/domain/client"
)

type OpenLoginMsg struct{}

type OpenRegisterMsg struct {
	LoginModel tea.Model
}

type SuccessMsg struct{}

// List

type OpenCategoriesMsg struct{}

type OpenItemsMsg struct {
	Category client.Category
}

type BackMsg struct {
	State int
	Item  *client.Item
}

// Item

type OpenItemMsg struct {
	BackModel tea.Model
	BackState int
	Item      *client.Item
}

type OpenItemMessage struct {
	BackModel tea.Model
	BackState int
	Item      *client.Item
	Category  client.Category
}
