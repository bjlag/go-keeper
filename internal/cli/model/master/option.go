package master

import (
	"github.com/bjlag/go-keeper/internal/cli/model/create"
	"github.com/bjlag/go-keeper/internal/cli/model/list"
)

type Option func(*Model)

func WithCreatForm(form *create.Model) Option {
	return func(m *Model) {
		form.SetMainModel(m)
		m.formCreate = form
	}
}

func WithListForm(form *list.Model) Option {
	return func(m *Model) {
		form.SetMainModel(m)
		m.formList = form
	}
}
