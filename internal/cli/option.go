package cli

import (
	"github.com/bjlag/go-keeper/internal/cli/form/login"
	"github.com/bjlag/go-keeper/internal/cli/form/register"
)

type Option func(*MainModel)

func WithLoginForm(form *login.Form) Option {
	return func(m *MainModel) {
		form.SetMainModel(m)
		m.loginForm = form
	}
}

func WithRegisterForm(form *register.Form) Option {
	return func(m *MainModel) {
		form.SetMainModel(m)
		m.registerForm = form
	}
}
