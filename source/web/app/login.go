package app

import (
  "net/http"
  "github.com/labstack/echo"
  "source/utils"
  "source/structs"
)

func LoginView(e echo.Context) error {
  if utils.IsNotLogined(e) {
    return e.Render(http.StatusOK, "login.html", "")
  }
  return e.Redirect(http.StatusFound, "/")
}

func LoginPost(e echo.Context) error {
  post := new(structs.LoginParams)
  if err := e.Bind(post); err != nil {
    return e.JSON(http.StatusOK, map[string] string { "error": err.Error() })
  }

  errorVal := ""
  if (post.Passphrase == utils.Yaml().Passphrase) {
    utils.SessionSave(e)
  } else {
    errorVal = "合言葉が正しくありません(T_T)"
  }
  return e.JSON(http.StatusOK, map[string] string { "error": errorVal })
}

func Logout(e echo.Context) error {
  utils.SessionDelete(e)
  return e.Redirect(http.StatusFound, "/login")
}