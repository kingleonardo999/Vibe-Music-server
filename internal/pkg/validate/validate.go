package validate

import (
	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

var (
	usernameRe         = regexp2.MustCompile(`^[a-zA-Z0-9_-]{4,16}$`, regexp2.None)
	passwordRe         = regexp2.MustCompile(`^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z\W]{8,18}$`, regexp2.None)
	verificationCodeRe = regexp2.MustCompile(`^[0-9a-zA-Z]{6}$`, regexp2.None)
	phoneRe            = regexp2.MustCompile(`^1[3-9]\d{9}$`, regexp2.None)
)

func username(fl validator.FieldLevel) bool {
	ok, _ := usernameRe.MatchString(fl.Field().String())
	return ok
}

func password(fl validator.FieldLevel) bool {
	ok, _ := passwordRe.MatchString(fl.Field().String())
	return ok
}

func verificationCode(fl validator.FieldLevel) bool {
	ok, _ := verificationCodeRe.MatchString(fl.Field().String())
	return ok
}

func phone(fl validator.FieldLevel) bool {
	ok, _ := phoneRe.MatchString(fl.Field().String())
	return ok
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		Validate = v // ← 关键：复用 Gin 的实例
	} else {
		panic("Gin validator is not *validator.Validate")
	}
	_ = Validate.RegisterValidation("username", username)
	_ = Validate.RegisterValidation("password", password)
	_ = Validate.RegisterValidation("verificationCode", verificationCode)
	_ = Validate.RegisterValidation("phone", phone)
}
