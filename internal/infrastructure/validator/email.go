package validator

import (
	"regexp"
	"sync"
)

var (
	once              sync.Once
	regexEmailPattern *regexp.Regexp
)

func ValidateEmail(email string) bool {
	once.Do(func() {
		regexEmailPattern = regexp.MustCompile("^[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$")
	})

	return regexEmailPattern.MatchString(email)
}
