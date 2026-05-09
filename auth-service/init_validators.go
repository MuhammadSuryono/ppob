package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	validate       *validator.Validate
	phoneRegex     = regexp.MustCompile(`^\+62[0-9]{8,12}$`)
	pinRegex       = regexp.MustCompile(`^[0-9]{6}$`)
	sequentialPINs = []string{
		"123456", "654321", "111111", "000000",
		"222222", "333333", "444444", "555555",
		"666666", "777777", "888888", "999999",
	}
)

func init() {
	validate = validator.New()

	_ = validate.RegisterValidation("phone_id", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		if phone == "" {
			return false
		}
		return phoneRegex.MatchString(phone)
	})

	_ = validate.RegisterValidation("pinformat", func(fl validator.FieldLevel) bool {
		pin := fl.Field().String()
		if len(pin) != 6 {
			return false
		}
		if !pinRegex.MatchString(pin) {
			return false
		}
		for _, seq := range sequentialPINs {
			if pin == seq {
				return false
			}
		}
		return true
	})

	_ = validate.RegisterValidation("password_complex", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		if len(password) < 8 {
			return false
		}

		hasUpper, hasLower, hasDigit := false, false, false
		for _, c := range password {
			switch {
			case c >= 'A' && c <= 'Z':
				hasUpper = true
			case c >= 'a' && c <= 'z':
				hasLower = true
			case c >= '0' && c <= '9':
				hasDigit = true
			}
		}

		return hasUpper && hasLower && hasDigit
	})
}