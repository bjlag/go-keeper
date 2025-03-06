package element

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
)

func GetValue(elements []interface{}, pos int) string {
	switch e := elements[pos].(type) {
	case textinput.Model:
		return e.Value()
	case textarea.Model:
		return e.Value()
	}

	return ""
}
