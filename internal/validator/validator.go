package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/padurean/purest/internal/auth"
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

	//--> register custom validations
	_ = validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return auth.IsStrongPassword(fl.Field().String()) == nil
	})
	_ = validate.RegisterValidation("role", func(fl validator.FieldLevel) bool {
		return auth.IsValidRole(int(fl.Field().Uint()))
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
	//----> custom translations
	_ = validate.RegisterTranslation(
		"password",
		translator,
		func(ut ut.Translator) error {
			return ut.Add("password", auth.PasswordRequirementsMsg, true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("password", fe.Field())
			return t
		},
	)
	_ = validate.RegisterTranslation(
		"role",
		translator,
		func(ut ut.Translator) error {
			return ut.Add("role", "invalid role - "+auth.ValidRolesMsg, true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("role", fe.Field())
			return t
		},
	)
	//<----
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
