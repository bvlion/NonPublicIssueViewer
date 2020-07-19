package app

import (
  "net/http"
  "source/utils"
  "github.com/labstack/echo"
)

func CreateIssue(e echo.Context) error {
  if e.Request().Header.Get("X-Appengine-Cron") != "true" {
    return e.Render(http.StatusNotFound, "404.html", "")
  }
  utils.CreateNewTodaysIssue(utils.Yaml().GitHub.Token, utils.Yaml().GitHub.User, utils.Yaml().GitHub.Project)
  return e.String(http.StatusOK, "")
}