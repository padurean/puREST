package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/rs/zerolog/log"
)

// use a single instance of Validate, it caches struct info
var (
	validate   *validator.Validate
	translator ut.Translator
)

func init() {
	validate = validator.New()

	//--> register JSON field names instead of field names for validation errors
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	//<--

	//--> register validator translation
	en := en.New()
	uni := ut.New(en, en)
	var found bool
	translator, found = uni.GetTranslator("en")
	if !found {
		log.Error().Msg("error geting english translator for validation errors")
		return
	}
	en_translations.RegisterDefaultTranslations(validate, translator)
	//<--
}

// Validate ...
func Validate(s interface{}) error {
	err := validate.Struct(s)
	if err != nil && translator != nil {
		// translate all error at once
		errsMap := err.(validator.ValidationErrors).Translate(translator)
		var errs []string
		for _, v := range errsMap {
			errs = append(errs, v)
		}
		return fmt.Errorf("%v", strings.Join(errs, ", "))
	}

	return nil
}
