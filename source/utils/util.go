package utils

import (
  "github.com/labstack/echo"
  "github.com/ipfans/echo-session"
)

const SessionName = "logined"

func IsNotLogined(e echo.Context) bool {
  session := session.Default(e)
  logined := session.Get(SessionName)

  return logined != "true"
}