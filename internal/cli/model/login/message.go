package login

type OpenMessage struct{}

type SuccessMessage struct {
	AccessToken  string
	RefreshToken string
}
