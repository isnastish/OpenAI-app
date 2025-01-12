package validator

import (
	"strings"
	"testing"
)

func TestValidateUserPassword(t *testing.T) {
	testData := []struct {
		input    string
		expected error
	}{
		{"", passwordTooShortError},
		{"hello", passwordTooShortError},
		{"this_password_is_supposed_to_be_longer_than_128_charactres_and_thus_it_wont_pass_the_validation_test_and_to_make_that_happen_we_need_to_add_some_extra_characters", passwordTooLongError},
		{"pass8worD", missingSpecialSymbolError},
		{"noNumber_@", missingNumberError},
		{"no_u2per_letter", missingUpperCaseError},
		{"NU_LOWER_LE2TER", missingLowerCaseError},
	}

	for _, data := range testData {
		got := ValidateUserPassword(data.input)
		if strings.Compare(data.expected.Error(), got.Error()) != 0 {
			t.Fatalf("expected: %s, got: %s", data.expected.Error(), got.Error())
		}
	}
}
