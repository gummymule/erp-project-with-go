package utils

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom validators with Gin's validator
		_ = v.RegisterValidation("sku", validateSKU)
		_ = v.RegisterValidation("phone", validatePhone)
		validate = v
	} else {
		// Fallback if Gin's validator is not available
		validate = validator.New()
		_ = validate.RegisterValidation("sku", validateSKU)
		_ = validate.RegisterValidation("phone", validatePhone)
	}
}

func validateSKU(fl validator.FieldLevel) bool {
	sku := fl.Field().String()

	// SKU should be alphanumeric with optional hyphens
	matched, _ := regexp.MatchString(`^[A-Za-z0-9\-]+$`, sku)
	return matched && len(sku) >= 3 && len(sku) <= 50
}

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	// Phone is optional, but if provided, validate it
	if phone == "" {
		return true
	}

	// Simple phone validation - adjust as needed
	matched, _ := regexp.MatchString(`^[\d\s\-\+\(\)]+$`, phone)
	return matched && len(phone) >= 10 && len(phone) <= 20
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// Validation error formatter
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func FormatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			errorMsg := getErrorMessage(fieldErr)
			errors = append(errors, ValidationError{
				Field:   strings.ToLower(fieldErr.Field()),
				Message: errorMsg,
			})
		}
	}

	return errors
}

func getErrorMessage(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	case "gt":
		return "Value must be greater than " + fieldErr.Param()
	case "gte":
		return "Value must be greater than or equal to " + fieldErr.Param()
	case "sku":
		return "SKU must be alphanumeric with hyphens, 3-50 characters"
	case "phone":
		return "Invalid phone number format"
	default:
		return "Invalid value"
	}
}
