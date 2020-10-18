package main

import (
	"net/http"
	"time"

	elogrus "github.com/dictor/echologrus"
	echo "github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	elogrus.Attach(e)
	e.GET("/", routePanic)

	go func() {
		time.Sleep(1 * time.Second)
		http.Get("http://127.0.0.1:28000/")
		time.Sleep(1 * time.Second)
		http.Get("http://127.0.0.1:28000/notexist")
	}()

	e.Logger.Fatal(e.Start("127.0.0.1:28000"))
}

func routePanic(c echo.Context) error {
	panic("panic causing test!")
}
