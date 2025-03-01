package message

import "github.com/bjlag/go-keeper/internal/cli/element"

type SuccessLoginMessage struct {
	AccessToken  string
	RefreshToken string
}

type SuccessRegisterMessage struct {
	AccessToken  string
	RefreshToken string
}

type OpenLoginFormMessage struct{}

type OpenRegisterFormMessage struct{}

type OpenCategoryListFormMessage struct{}

type OpenPasswordFormMessage struct {
	Item element.Item
}

type OpenPasswordListFormMessage struct {
	Category element.Item
}
