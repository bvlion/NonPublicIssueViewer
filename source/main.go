package main

import (
  "html/template"
  "net/http"
  "io"
  "os"
  "fmt"
  "github.com/labstack/echo"
  "github.com/ipfans/echo-session"
  "source/utils"
  "source/log"
)

const sessionName = "logined"

type Template struct {
  templates *template.Template
}

type LoginParams struct {
  Passphrase string `json:"passphrase"`
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
  return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
  log.LogSetUp()

  t := &Template{
    templates: template.Must(template.ParseGlob("views/*.html")),
  }

  e := echo.New()

  http.Handle("/", e)

  e.Renderer = t

  e.HTTPErrorHandler = func(err error, c echo.Context) {
    if he, ok := err.(*echo.HTTPError); ok {
      code := he.Code
      log.ErrorLog(err.Error())
      if code == http.StatusNotFound {
        c.Render(code, "404.html", "")
      } else {
        c.JSON(code, map[string] string { "error": err.Error() })
      }
    }
  }

  //セッションを設定
  store := session.NewCookieStore([]byte(utils.Yaml().SessionKey))
  store.MaxAge(86400)
  e.Use(session.Sessions("ESESSION", store))

  e.Static("/css", "./public/css")
  e.Static("/js", "./public/js")
  e.Static("/images", "./public/images")

  e.GET("/", func (e echo.Context) error {
    session := session.Default(e)
    logined := session.Get(sessionName)

    if logined != "true" {
      return e.Redirect(http.StatusFound, "/login")
    }

    log.DebugLog(utils.ReadIssues(utils.Yaml().GitHub.Token, utils.Yaml().GitHub.User, utils.Yaml().GitHub.Project))
    return e.Render(http.StatusOK, "index.html", "")
  })

  e.GET("/login", func (e echo.Context) error {
    session := session.Default(e)
    logined := session.Get(sessionName)

    if logined == "true" {
      return e.Redirect(http.StatusFound, "/")
    }
    return e.Render(http.StatusOK, "login.html", "")
  })

  e.POST("/login", func (e echo.Context) error {
    post := new(LoginParams)
    if err := e.Bind(post); err != nil {
      return e.JSON(http.StatusOK, map[string] string { "error": err.Error() })
    }

    errorVal := ""
    if (post.Passphrase == utils.Yaml().Passphrase) {
      session := session.Default(e)
      session.Set(sessionName, "true")
      session.Save()
    } else {
      errorVal = "合言葉が正しくありません(T_T)"
    }
    return e.JSON(http.StatusOK, map[string] string { "error": errorVal })
  })

  e.GET("/logout", func (e echo.Context) error {
    session := session.Default(e)
    session.Clear()

    return e.Redirect(http.StatusFound, "/login")
  })

  e.GET("/create-issue", func (e echo.Context) error {
    if e.Request().Header.Get("X-Appengine-Cron") != "true" {
      return e.Render(http.StatusNotFound, "404.html", "")
    }
    utils.CreateNewTodaysIssue(utils.Yaml().GitHub.Token, utils.Yaml().GitHub.User, utils.Yaml().GitHub.Project)
    return e.String(http.StatusOK, "")
  })

  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
    log.DebugLog("Defaulting to port", port)
  }

  log.DebugLog("Listening on port", port)
  log.FatalLog(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
