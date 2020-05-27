package main

import (
  "html/template"
  "net/http"
  "io"
  "fmt"
  "github.com/labstack/echo"
  "log"
  "os"
  "main/utils"
)

type Template struct {
  templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
  return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

  t := &Template{
    templates: template.Must(template.ParseGlob("views/*.html")),
  }

  e := echo.New()

  http.Handle("/", e)

  e.Renderer = t

  e.HTTPErrorHandler = func(err error, c echo.Context) {
    if he, ok := err.(*echo.HTTPError); ok {
      code := he.Code
      fmt.Println(err)
      c.JSON(code, err.Error())
    }
  }

  e.Static("/css", "./public/css")
  e.Static("/js", "./public/js")
  e.Static("/images", "./public/images")

  e.GET("/", func (e echo.Context) error {
    fmt.Println(utils.ReadIssues(utils.Yaml().GitHub.Token, utils.Yaml().GitHub.User, utils.Yaml().GitHub.Project))
    return e.Render(http.StatusOK, "index.html", "")
  })
  e.GET("/login", func (e echo.Context) error {
    return e.Render(http.StatusOK, "login.html", "")
  })

  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
    log.Printf("Defaulting to port %s", port)
  }

  log.Printf("Listening on port %s", port)
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
