package register

import tea "github.com/charmbracelet/bubbletea"

type SuccessMsg struct {
	AccessToken  string
	RefreshToken string
}

type OpenMessage struct {
	BackModel tea.Model
}
