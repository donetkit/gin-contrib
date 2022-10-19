package gin_bind

import (
	"github.com/donetkit/contrib-gin/pkg/gin/validation"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Bind(s interface{}, c *gin.Context) (interface{}, error) {
	b := binding.Default(c.Request.Method, c.ContentType())
	if err := c.ShouldBindWith(s, b); err != nil {
		return nil, err
	}
	return s, nil
}

// BindAndValid binds and validates data
func BindAndValid(c *gin.Context, s interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	if err := c.ShouldBindWith(s, b); err != nil {
		//return nil, err
	}
	valid := validation.Validation{}
	check, err := valid.Valid(s)
	if err != nil {
		return err
	}
	if !check {
		return valid.Errors[0]
	}
	return nil
}

// ShouldBindAndValid binds and validates data
func ShouldBindAndValid(c *gin.Context, s interface{}) error {
	if err := c.ShouldBind(s); err != nil {
		//return nil, err
	}
	valid := validation.Validation{}
	check, err := valid.Valid(s)
	if err != nil {
		return err
	}
	if !check {
		return valid.Errors[0]
	}
	return nil
}

func ShouldBindAndValidV10(c *gin.Context, s interface{}) error {
	if err := c.ShouldBindWith(&s, binding.Query); err == nil {
		return err
	} else {
		_, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}
	}
	return nil
}
