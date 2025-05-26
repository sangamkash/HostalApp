package ValidatorSystem

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"sync"
)

type ValidatorManager struct {
	validate *validator.Validate
}

var getValidatorManager = sync.OnceValue(func() *ValidatorManager {
	validate := validator.New()

	// Custom phone validation
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	_ = validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		return phoneRegex.MatchString(fl.Field().String())
	})

	// Optional: Register strong password validation
	_ = validate.RegisterValidation("strong_password", func(fl validator.FieldLevel) bool {
		return isPasswordStrong(fl.Field().String())
	})

	return &ValidatorManager{
		validate: validate,
	}
})

// Strong password check
func isPasswordStrong(password string) bool {
	var (
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString
		hasSpecial = regexp.MustCompile(`[!@#$%^&*()\-_=+{};:,<.>]`).MatchString
	)

	return hasUpper(password) &&
		hasLower(password) &&
		hasNumber(password) &&
		hasSpecial(password) &&
		len(password) >= 8
}

// Public accessor (thread-safe lazy initialization)
func GetValidator() *ValidatorManager {
	return getValidatorManager()
}

// Validate a struct with registered rules
func (m *ValidatorManager) IsValid(data interface{}) error {
	return m.validate.Struct(data)
}
