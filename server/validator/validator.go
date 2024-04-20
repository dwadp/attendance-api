package validator

import (
	"errors"
	"fmt"
	"github.com/dwadp/attendance-api/store/db"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"reflect"
	"regexp"
	"time"
)

type Validator struct {
	validate           *validator.Validate
	trans              ut.Translator
	timeRegexValidator *regexp.Regexp
}

func New() (*Validator, error) {
	validate := validator.New()

	english := en.New()
	uni := ut.New(english, english)
	trans, ok := uni.GetTranslator("en")

	if !ok {
		return nil, errors.New("translator not found")
	}

	if err := enTranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		return nil, fmt.Errorf("could not register default translations: %w", err)
	}

	v := Validator{validate: validate, trans: trans}
	if err := v.registerRegexes(); err != nil {
		return nil, err
	}

	if err := v.registerTimeValidator(); err != nil {
		return nil, err
	}

	return &v, nil
}

func (v *Validator) registerRegexes() error {
	timeRegex, err := regexp.Compile("^([0-1]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$")
	if err != nil {
		return err
	}

	v.timeRegexValidator = timeRegex

	return nil
}

func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

func (v *Validator) SerializeErrors(err error, kind any) map[string]string {
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	if !ok {
		return nil
	}

	result := map[string]string{}
	t := reflect.TypeOf(kind)

	for _, e := range errs {
		field, ok := t.FieldByName(e.Field())
		msg := e.Translate(v.trans)

		if !ok {
			result[e.Field()] = msg
		}

		tag := field.Tag.Get("json")
		result[tag] = msg
	}

	return result
}

func (v *Validator) registerTimeValidator() (err error) {
	err = v.validate.RegisterValidation("time", func(fl validator.FieldLevel) bool {
		value := fl.Field().Interface()

		switch a := value.(type) {
		case db.Time:
			return v.timeRegexValidator.MatchString(a.String())
		case time.Time:
			return v.timeRegexValidator.MatchString(a.Format("15:04:05"))
		case string:
			return v.timeRegexValidator.MatchString(a)
		}

		return false
	})

	if err != nil {
		return err
	}

	err = v.validate.RegisterTranslation(
		"time",
		v.trans,
		func(ut ut.Translator) error {
			return ut.Add("time", "{0} must be a valid time format (hh:mm)", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("time", fe.Field())
			return t
		},
	)

	if err != nil {
		return err
	}

	return nil
}
