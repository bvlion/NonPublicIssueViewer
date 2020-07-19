package utils

import (
  "github.com/labstack/echo"
  "github.com/ipfans/echo-session"
)

const sessionName = "logined"

func IsNotLogined(e echo.Context) bool {
  session := session.Default(e)
  logined := session.Get(sessionName)

  return logined != "true"
}

func SessionSave(e echo.Context) {
  session := session.Default(e)
  session.Set(sessionName, "true")
  session.Save()
}

func SessionDelete(e echo.Context) {
  session := session.Default(e)
  session.Clear()
  session.Save()
}