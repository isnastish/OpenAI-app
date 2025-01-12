package validator

import (
	"errors"
	"regexp"
)

// NOTE: Maybe password validation better be done on the client side?

var (
	passwordTooShortError     = errors.New("password is too short")
	passwordTooLongError      = errors.New("password is too long")
	missingUpperCaseError     = errors.New("passwrod should contain at least 1 upped case")
	missingLowerCaseError     = errors.New("password should contain at least 1 lower case")
	missingNumberError        = errors.New("password should contain at least 1 number")
	missingSpecialSymbolError = errors.New("password should contain at least 1 special symbol")
)

func ValidateUserPassword(password string) error {
	// A valid password should contain:
	//
	// at least 8 characters, maximum 128
	pwdLength := len(password)
	if pwdLength < 8 {
		return passwordTooShortError
	} else if pwdLength > 128 {
		return passwordTooLongError
	}

	// at least 1 upper case (A-Z)
	upperRule := regexp.MustCompile("[A-Z]")
	if !upperRule.MatchString(password) {
		return missingUpperCaseError
	}

	// at least 1 lower case (a-z)
	lowerRule := regexp.MustCompile("[a-z]")
	if !lowerRule.MatchString(password) {
		return missingLowerCaseError
	}

	// at least 1 number (0-9)
	numberRule := regexp.MustCompile("[0-9]")
	if !numberRule.MatchString(password) {
		return missingNumberError
	}

	// at least 1 special character (~`!@#$%^&*()_-+={[}]|\:;"'<,>.?/)
	symbolRule := regexp.MustCompile("[^\\d\\w]")
	if !symbolRule.MatchString(password) {
		return missingSpecialSymbolError
	}

	return nil
}
