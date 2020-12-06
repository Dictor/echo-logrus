## echologrus
[![PkgGoDev](https://pkg.go.dev/badge/Dictor/echologrus)](https://pkg.go.dev/Dictor/echologrus)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Middleware echologrus is a [logrus](https://github.com/sirupsen/logrus) logger support for [echo](https://github.com/labstack/echo).
Only support echo **v4**.

### Install

```sh
go get -u github.com/dictor/echologrus
```

### Example
#### Basic
```go
import (
	elogrus "github.com/dictor/echologrus"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()
	elogrus.Attach(e)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":80"))
}
```

#### Using custom formatter
```go
	e := echo.New()
	elogrus.Attach(e).Logger.Formatter = new(prefixed.TextFormatter)
```


