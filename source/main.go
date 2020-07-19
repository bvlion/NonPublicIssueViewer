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

  e := echo.New()
  http.Handle("/", e)

  e.Renderer = &Template{
    templates: template.Must(template.ParseGlob("web/template/*.html")),
  }
  e.HTTPErrorHandler = app.ProjectHandler

  // Set Session
  store := session.NewCookieStore([]byte(utils.Yaml().SessionKey))
  store.MaxAge(86400)
  e.Use(session.Sessions("ESESSION", store))

  // Set Static
  e.Static("/css", "./web/static/css")
  e.Static("/js", "./web/static/js")
  e.Static("/images", "./web/static/images")

  // Set Index
  e.GET("/", app.IndexView)
  e.GET("/detail/:year/:month/:day", app.IndexDetailJson)
  e.GET("/issues/:minusMonth", app.IndexDataJson)

  // Set Login
  e.GET("/login", app.LoginView)
  e.POST("/login", app.LoginPost)
  e.GET("/logout", app.Logout)

  // Set GitHub
  e.GET("/create-issue", app.CreateIssue)

  // Set Port
  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
    log.DebugLog("Defaulting to port", port)
  }

  log.DebugLog("Listening on port", port)
  log.FatalLog(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
