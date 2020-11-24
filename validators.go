package main

import (
	"errors"
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
)

// type ValidationAPI interface {
// 	// v *validator.Validate
// 	// RegisterValidations  func()
// 	// RegisterTranslations func()
// 	registerCustomErrors() error
// }

type ReminderValidator struct {
	v *validator.Validate
}

func NewReminderValidator() (*ReminderValidator, error) {
	rv := &ReminderValidator{v: validator.New()}
	err := rv.registerCustomErrors()
	if err != nil {
		return nil, err
	}

	return rv, nil
}

// func (rv *ReminderValidator) registerValidations() error {
// 	// _ = rv.v.RegisterValidation("passwd", func(fl validator.FieldLevel) bool {
// 	// 	return len(fl.Field().String()) > 6
// 	// })
// }

func (rv *ReminderValidator) registerCustomErrors() error {
	translator := en.New()
	uni := ut.New(translator, translator)

	trans, found := uni.GetTranslator("en")
	if !found {
		return errors.New("ReminderValidator: translator not found")
	}

	err := rv.v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})
	if err != nil {
		return errors.New(fmt.Sprintf("ReminderValidator: required %v", err))
	}

	err = rv.v.RegisterTranslation("latitude", trans, func(ut ut.Translator) error {
		return ut.Add("latitude", "{0} must be a valid latitude", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("latitude", fe.Field())
		return t
	})

	if err != nil {
		return errors.New(fmt.Sprintf("ReminderValidator: required %v", err))
	}

	err = rv.v.RegisterTranslation("longitude", trans, func(ut ut.Translator) error {
		return ut.Add("longitude", "{0} must be a valid latitude", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("longitude", fe.Field())
		return t
	})

	if err != nil {
		return errors.New(fmt.Sprintf("ReminderValidator: required %v", err))
	}

	return nil
}
