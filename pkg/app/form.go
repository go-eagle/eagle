package app

import (
	"strings"

	"github.com/gin-gonic/gin"
	// nolint: typecheck
	ut "github.com/go-playground/universal-translator"
	// nolint: typecheck
	val "github.com/go-playground/validator/v10"
)

// ValidError .
type ValidError struct {
	Key     string
	Message string
}

// ValidErrors .
type ValidErrors []*ValidError

// Error return error msg
func (v *ValidError) Error() string {
	return v.Message
}

// Error return error string
func (v ValidErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

// Errors return some error
func (v ValidErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

// BindAndValid valid params
func BindAndValid(c *gin.Context, v interface{}) (bool, ValidErrors) {
	var errs ValidErrors
	err := c.ShouldBind(v)
	if err != nil {
		v := c.Value("trans")
		trans, _ := v.(ut.Translator)
		verrs, ok := err.(val.ValidationErrors)
		if !ok {
			return false, errs
		}

		for key, value := range verrs.Translate(trans) {
			errs = append(errs, &ValidError{
				Key:     key,
				Message: value,
			})
		}

		return false, errs
	}

	return true, nil
}
