package register

type Data struct {
	Email    string
	Password string
}

type Result struct {
	AccessToken  string
	RefreshToken string
}
