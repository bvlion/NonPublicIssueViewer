package main

import (
  "fmt"
  "html/template"
  "io"
  "net/http"
  "os"
  "source/log"
  "source/utils"
  "source/web/app"

  session "github.com/ipfans/echo-session"
  "github.com/labstack/echo"
)

type Template struct {
  templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
  return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
  log.LogSetUp()

  funcMap := template.FuncMap{
    "safehtml": func(text string) template.HTML { return template.HTML(text) },
  }
  t := &Template{
    templates: template.Must(template.New("").Funcs(funcMap).ParseGlob("web/template/*.html")),
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
        c.JSON(code, map[string]string{"error": err.Error()})
      }
    }
  }

  //セッションを設定
  store := session.NewCookieStore([]byte(utils.Yaml().SessionKey))
  store.MaxAge(86400)
  e.Use(session.Sessions("ESESSION", store))

  e.Static("/css", "./web/static/css")
  e.Static("/js", "./web/static/js")
  e.Static("/images", "./web/static/images")

  e.GET("/", app.IndexView)
  e.GET("/detail/:year/:month/:day", app.IndexDetailJson)
  e.GET("/issues/:minusMonth", app.IndexDataJson)

  e.GET("/login", app.LoginView)
  e.POST("/login", app.LoginPost)
  e.GET("/logout", app.Logout)

  e.GET("/create-issue", func(e echo.Context) error {
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
