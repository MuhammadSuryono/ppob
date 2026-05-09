package validation

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

var (
	phoneRegex    = regexp.MustCompile(`^\+62[0-9]{8,12}$`)
	pinRegex      = regexp.MustCompile(`^[0-9]{6}$`)
	sequentialPIN = []string{"123456", "654321", "111111", "000000", "222222", "333333", "444444", "555555", "666666", "777777", "888888", "999999"}
	customerNoMin = 4
	customerNoMax = 25
)

func init() {
	_ = validate.RegisterValidation("phone_id", validatePhone)
	_ = validate.RegisterValidation("pinformat", validatePIN)
	_ = validate.RegisterValidation("customer_no", validateCustomerNo)
	_ = validate.RegisterValidation("password_complex", validatePassword)
}

func ValidateUserCreate(req interface{}) error {
	return validate.Struct(req)
}

func ValidatePhone(phone string) bool {
	if phone == "" {
		return false
	}
	return phoneRegex.MatchString(phone)
}

func ValidatePIN(pin string) bool {
	if len(pin) != 6 {
		return false
	}
	if !pinRegex.MatchString(pin) {
		return false
	}
	for _, seq := range sequentialPIN {
		if pin == seq {
			return false
		}
	}
	return true
}

func ValidateCustomerNo(no string) bool {
	cleaned := strings.ReplaceAll(no, " ", "")
	if len(cleaned) < customerNoMin || len(cleaned) > customerNoMax {
		return false
	}
	return true
}

func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, c := range password {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasDigit = true
		case strings.Contains("!@#$%^&*()_+-=[]{}|;:,.<>?", string(c)):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return false
	}
	return phoneRegex.MatchString(phone)
}

func validatePIN(fl validator.FieldLevel) bool {
	pin := fl.Field().String()
	if len(pin) != 6 {
		return false
	}
	if !pinRegex.MatchString(pin) {
		return false
	}
	for _, seq := range sequentialPIN {
		if pin == seq {
			return false
		}
	}
	return true
}

func validateCustomerNo(fl validator.FieldLevel) bool {
	no := fl.Field().String()
	cleaned := strings.ReplaceAll(no, " ", "")
	if len(cleaned) < customerNoMin || len(cleaned) > customerNoMax {
		return false
	}
	return true
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, c := range password {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasDigit = true
		case strings.Contains("!@#$%^&*()_+-=[]{}|;:,.<>?", string(c)):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

func GetValidationErrors(err error) []string {
	if err == nil {
		return nil
	}

	var errors []string
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, err.Field()+": "+err.Tag())
	}
	return errors
}