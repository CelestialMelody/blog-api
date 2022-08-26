package validate

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"regexp"
)

var valid = validator.New()

func Struct(s interface{}) error {
	return valid.Struct(s)
}

func Var(field interface{}, tag string) error {
	return valid.Var(field, tag)
}

var phoneNumber validator.Func = func(fl validator.FieldLevel) bool {
	ok, _ := regexp.MatchString(`^1[3-9]\d{9}$`, fl.Field().String())
	return ok
}

func Init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("mobile", phoneNumber)
	}
}
