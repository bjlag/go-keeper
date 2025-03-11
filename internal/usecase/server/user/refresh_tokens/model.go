package refresh_tokens

type Data struct {
	RefreshToken string
}

type Result struct {
	AccessToken  string
	RefreshToken string
}
