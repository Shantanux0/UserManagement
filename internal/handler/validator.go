package handler

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()

	_ = Validate.RegisterValidation("past_date", func(fl validator.FieldLevel) bool {
		dateStr := fl.Field().String()
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return false
		}
		return t.Before(time.Now())
	})
}
