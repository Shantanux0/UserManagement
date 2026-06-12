package handler

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Validate is the global validator instance.
var Validate *validator.Validate

// InitValidator initializes the validator and registers custom validation functions.
func InitValidator() {
	Validate = validator.New()

	// Register custom validator to ensure DOB is in the past
	_ = Validate.RegisterValidation("past_date", func(fl validator.FieldLevel) bool {
		dateStr := fl.Field().String()
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return false
		}
		// A birthday must be before the current time
		return t.Before(time.Now())
	})
}
