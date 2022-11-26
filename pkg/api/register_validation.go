package api

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/xyedo/blindate/pkg/domain/validation"
)

func registerValidDObValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("validdob", validation.ValidDob)
		if err != nil {
			panic(err)
		}
	} else {
		panic("not ok validator")
	}
}
func registerTagName() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "" {
				return fld.Name
			}

			if name == "-" {
				return ""
			}
			return name
		})
	} else {
		panic("not ok validator")
	}

}
