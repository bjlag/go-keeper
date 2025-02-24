package ui

import (
	"github.com/bjlag/go-keeper/internal/app/client/ui/forms/login"
	key2 "github.com/bjlag/go-keeper/internal/app/client/ui/key"
	"github.com/bjlag/go-keeper/internal/app/client/ui/style"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const initCountInputs = 2

type state int

const (
	loginState state = iota
	listCategories
	listPasswords
	viewPassword
)

type Model struct {
	help       help.Model
	state      state
	inputs     []textinput.Model
	focusIndex int
}

func InitModel() *Model {
	m := &Model{
		help:   help.New(),
		state:  loginState,
		inputs: make([]textinput.Model, initCountInputs),
	}

	for i := range m.inputs {
		t := textinput.New()
		t.Cursor.Style = style.CursorStyle
		t.PlaceholderStyle = style.BlurredStyle

		switch i {
		case 0:
			t.CharLimit = 20
			t.Focus()
			t.Placeholder = "Email"
			t.TextStyle = style.FocusedStyle
			t.PromptStyle = style.FocusedStyle
		case 1:
			t.CharLimit = 60
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		func() tea.Msg {
			return login.NewLoginForm()
		},
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key2.Keys.Quit):
			return m, tea.Quit
		}
	case *login.Form:
		return msg.Update(msg)
	}

	switch m.state {
	case loginState:
		//return m.updateLogin(msg)
	}

	return m, nil
}

func (m *Model) View() string {
	if m.state == loginState {

	}

	return ""
}
