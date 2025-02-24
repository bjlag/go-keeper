package client

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/app/client/ui"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Run(ctx context.Context) error {
	_, err := tea.NewProgram(ui.InitModel(), tea.WithAltScreen()).Run()
	return err
}
