package cli

import (
	"github.com/bjlag/go-keeper/internal/cli/form/list"
	"github.com/bjlag/go-keeper/internal/cli/form/login"
	"github.com/bjlag/go-keeper/internal/cli/form/password"
	"github.com/bjlag/go-keeper/internal/cli/form/register"
	"github.com/bjlag/go-keeper/internal/infrastructure/store/client/token"
)

type Option func(*MainModel)

func WithStoreTokens(store *token.Store) Option {
	return func(m *MainModel) {
		m.storeTokens = store
	}
}

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
