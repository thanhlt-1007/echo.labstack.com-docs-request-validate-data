package main

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type User struct {
	Name  string `json:"name" form:"name" validate:"required"`
	Email string `json:"email" form:"email" validate:"required,email"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func (custmValidator *CustomValidator) Validate(i interface{}) error {
	err := custmValidator.validator.Struct(i)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			err.Error(),
		)
	}

	return nil
}

func postCreateUserHandler(context echo.Context) error {
	var user User
	err := context.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadGateway,
			err.Error(),
		)
	}

	err = context.Validate(user)
	if err != nil {
		return err
	}

	return context.JSON(
		http.StatusOK,
		user,
	)
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &CustomValidator{validator: validator.New()}
	e.POST("/users", postCreateUserHandler)
	e.Logger.Fatal(e.Start(":1323"))
}
