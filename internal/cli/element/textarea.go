package element

import "github.com/charmbracelet/bubbles/textarea"

func CreateDefaultTextArea(placeholder string) textarea.Model {
	m := textarea.New()
	m.Placeholder = placeholder

	return m
}
