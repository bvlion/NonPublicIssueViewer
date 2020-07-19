package app

import (
  "net/http"
  "source/log"
  "github.com/labstack/echo"
)

func ProjectHandler(err error, c echo.Context) {
  if he, ok := err.(*echo.HTTPError); ok {
    code := he.Code
    log.ErrorLog(err.Error())
    if code == http.StatusNotFound {
      c.Render(code, "404.html", "")
    } else {
      c.JSON(code, map[string]string{"error": err.Error()})
    }
  }
}