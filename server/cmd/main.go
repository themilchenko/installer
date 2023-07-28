package main

import (
	"github.com/labstack/echo/v4"

	"server/internal/app"
)

func main() {
	e := echo.New()
	s := app.New(e)

	if err := s.Start(); err != nil {
		s.Echo.Logger.Error(err)
	}
}
