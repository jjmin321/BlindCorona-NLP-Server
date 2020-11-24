package main

import (
	"BlindCorona-LanguageProcessingServer/controller"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/analyze", controller.Analyze)
	e.Logger.Fatal(e.Start(":3000"))
}
