package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
)

type ReminderValidator struct {
	v     *validator.Validate
	trans ut.Translator
}

func NewReminderValidator() (*ReminderValidator, error) {
	rv := &ReminderValidator{v: validator.New()}
	translator := en.New()
	uni := ut.New(translator, translator)

	trans, found := uni.GetTranslator("en")
	if !found {
		return nil, errors.New("ReminderValidator: translator not found")
	}

	rv.trans = trans

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

	err := rv.v.RegisterTranslation("required", rv.trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})
	if err != nil {
		return errors.New(fmt.Sprintf("ReminderValidator: required %v", err))
	}

	err = rv.v.RegisterTranslation("latitude", rv.trans, func(ut ut.Translator) error {
		return ut.Add("latitude", "{0} must be a valid latitude", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("latitude", fe.Field())
		return t
	})

	if err != nil {
		return errors.New(fmt.Sprintf("ReminderValidator: required %v", err))
	}

	err = rv.v.RegisterTranslation("longitude", rv.trans, func(ut ut.Translator) error {
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

func (rv *ReminderValidator) Validate(d interface{}) interface{} {
	err := rv.v.Struct(d)
	if err == nil {
		return err
	}

	er := make(map[string]string)

	for _, e := range err.(validator.ValidationErrors) {
		er[strings.ToLower(e.Field())] = strings.ToLower(e.Translate(rv.trans))
	}

	return er
}
