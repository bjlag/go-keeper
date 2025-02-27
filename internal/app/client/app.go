package client

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjlag/go-keeper/internal/cli"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Run(ctx context.Context) error {
	_, err := tea.NewProgram(cli.InitModel(), tea.WithAltScreen(), tea.WithContext(ctx)).Run()
	return err
}
