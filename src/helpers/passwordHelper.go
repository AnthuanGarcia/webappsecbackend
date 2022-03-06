package helpers

import (
	"log"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var patterns = []string{`[a-z]`, `[A-Z]`, `\d`, `\W`}

// hashPassword -  used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Panic(err)
	}

	return string(bytes)

}

// Validate Password - validates a user password with the following patterns:
// An Uppercase ([A-Z]), a lowercase ([a-z]), a symbol (\W) and a digit (\d)
func ValidatePassword(password string) bool {

	for _, pattern := range patterns {

		valid, _ := regexp.MatchString(pattern, password)

		if !valid {
			return false
		}

	}

	return true

}
