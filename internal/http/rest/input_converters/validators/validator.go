package validators

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/KompiTech/itsm-ticket-management-service/internal/http/rest/presenters"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

func NewPayloadValidator() PayloadValidator {
	return &payloadValidator{}
}

type payloadValidator struct {}


func (v payloadValidator) Validate(payload interface{}) error {
	validate := validator.New()

	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	// use the names which have been specified for JSON representations of structs, rather than normal Go field names
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return fmt.Sprintf("'%s'", name)
	})

	err := validate.Struct(payload)
	if err != nil {
		errs := translateError(err, trans)
		var errStrings []string
		for _, transErr := range errs {
			errStrings = append(errStrings, transErr.Error())
		}
		translErr := strings.Join(errStrings, ", ")
		return presenters.NewErrorf(http.StatusBadRequest, translErr)
	}

	return nil
}

func translateError(err error, trans ut.Translator) (errs []error) {
	if err == nil {
		return nil
	}

	validatorErrs := err.(validator.ValidationErrors)

	validatorErrs.Translate(trans)

	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(trans))
		errs = append(errs, translatedErr)
	}

	return errs
}
