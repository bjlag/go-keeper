package master

import (
	"github.com/bjlag/go-keeper/internal/cli/model/item/create"
	"github.com/bjlag/go-keeper/internal/cli/model/item/password"
	"github.com/bjlag/go-keeper/internal/cli/model/item/text"
	"github.com/bjlag/go-keeper/internal/cli/model/list"
	"github.com/bjlag/go-keeper/internal/cli/model/login"
	"github.com/bjlag/go-keeper/internal/cli/model/register"
)

type Option func(*Model)

func WithLoginForm(form *login.Model) Option {
	return func(m *Model) {
		form.SetMainModel(m)
		m.formLogin = form
	}
}

func WithRegisterForm(form *register.Model) Option {
	return func(m *Model) {
		form.SetMainModel(m)
		m.formRegister = form
	}
}

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

func WithPasswordItemForm(form *password.Model) Option {
	return func(m *Model) {
		form.SetMainModel(m)
		m.formPassword = form
	}
}

func WithTextItemForm(form *text.Model) Option {
	return func(m *Model) {
		form.SetMainModel(m)
		m.formText = form
	}
}
