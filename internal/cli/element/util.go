package element

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
)

// GetValue вспомогательная функция, которая достает из переданного elements элемент
// в позиции pos и возвращает значение модели.
func GetValue(elements []interface{}, pos int) string {
	switch e := elements[pos].(type) {
	case textinput.Model:
		return e.Value()
	case textarea.Model:
		return e.Value()
	}

	return ""
}
