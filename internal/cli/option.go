package cli

import (
	"github.com/bjlag/go-keeper/internal/cli/form/list"
	"github.com/bjlag/go-keeper/internal/cli/form/login"
	"github.com/bjlag/go-keeper/internal/cli/form/password"
	"github.com/bjlag/go-keeper/internal/cli/form/register"
)

type Option func(*MainModel)

func WithLoginForm(form *login.Form) Option {
	return func(m *MainModel) {
		form.SetMainModel(m)
		m.formLogin = form
	}
}

func WithRegisterForm(form *register.Form) Option {
	return func(m *MainModel) {
		form.SetMainModel(m)
		m.formRegister = form
	}
}

func WithListFormForm(form *list.Form) Option {
	return func(m *MainModel) {
		form.SetMainModel(m)
		m.formList = form
	}
}

func WithShowPasswordForm(form *password.Form) Option {
	return func(m *MainModel) {
		form.SetMainModel(m)
		m.formPassword = form
	}
}
