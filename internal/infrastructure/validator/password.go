package validator

const minPasswordLength = 8

// ValidatePassword валидация пароля.
func ValidatePassword(password string) bool {
	return !(len(password) < minPasswordLength)
}
