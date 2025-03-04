package register

import tea "github.com/charmbracelet/bubbletea"

type SuccessMessage struct {
	AccessToken  string
	RefreshToken string
}

type OpenMessage struct {
	BackModel tea.Model
}
