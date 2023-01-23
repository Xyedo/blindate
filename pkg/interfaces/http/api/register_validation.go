package api

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	basicInfoEntity "github.com/xyedo/blindate/pkg/domain/basicinfo/entities"
	userEntity "github.com/xyedo/blindate/pkg/domain/user/entities"
)

func registerValidDObValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("validdob", userEntity.ValidDob)
		if err != nil {
			panic(err)
		}
	} else {
		panic("not ok validator")
	}
}
func registerValidEducationLevelFieldValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("valideducationlevel", basicInfoEntity.ValidEducationLevel)
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
