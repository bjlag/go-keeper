package master

import (
	"github.com/bjlag/go-keeper/internal/cli/model/item"
	"github.com/bjlag/go-keeper/internal/cli/model/list"
	"github.com/bjlag/go-keeper/internal/cli/model/login"
	"github.com/bjlag/go-keeper/internal/cli/model/register"
	"github.com/bjlag/go-keeper/internal/infrastructure/store/client/token"
)

type Option func(*Model)

func WithStoreTokens(store *token.Store) Option {
	return func(m *Model) {
		m.storeTokens = store
	}
}

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

func WithListFormForm(form *list.Model) Option {
	return func(m *Model) {
		form.SetMainModel(m)
		m.formList = form
	}
}

func WithShowPasswordForm(form *item.Model) Option {
	return func(m *Model) {
		form.SetMainModel(m)
		m.formPassword = form
	}
}
