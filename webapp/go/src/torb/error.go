package main

import (
	"github.com/labstack/echo"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	c.Logger().Error(err)
}
