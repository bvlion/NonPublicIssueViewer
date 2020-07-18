package utils

import (
  "github.com/labstack/echo"
  "github.com/ipfans/echo-session"
)

const SessionName = "logined"

const StartDate = "2020/05/13"
const DateFormat = "2006/01/02"
const MonthFormat = "2006/01"

func IsNotLogined(e echo.Context) bool {
  session := session.Default(e)
  logined := session.Get(SessionName)

  return logined != "true"
}