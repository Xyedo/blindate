package api

import (
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
	}
}
