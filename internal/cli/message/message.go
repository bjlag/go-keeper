package message

type LoginSuccessMessage struct {
	AccessToken  string
	RefreshToken string
}

type RegisterSuccessMessage struct {
	AccessToken  string
	RefreshToken string
}

type OpenLoginFormMessage struct{}

type OpenRegisterFormMessage struct{}
