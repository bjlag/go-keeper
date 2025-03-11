package validator

const minPasswordLength = 8

func ValidatePassword(password string) bool {
	return !(len(password) < minPasswordLength)
}
