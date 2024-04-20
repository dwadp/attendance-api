package validator

import (
	"errors"
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"reflect"
)

type Validator struct {
	validate *validator.Validate
	trans    ut.Translator
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

	return &Validator{validate: validate, trans: trans}, nil
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
