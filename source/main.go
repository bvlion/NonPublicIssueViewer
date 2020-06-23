package main

import (
  "html/template"
  "net/http"
  "io"
  "os"
  "fmt"
  "time"
  "github.com/labstack/echo"
  "github.com/ipfans/echo-session"
  "source/utils"
  "source/log"
  "regexp"
  "strings"
)

const sessionName = "logined"
const startDate = "2020/05/13"
const dateFormat = "2006/01/02"
const monthFormat = "2006/01"

type Template struct {
  templates *template.Template
}

type LoginParams struct {
  Passphrase string `json:"passphrase"`
}

type IndexData struct {
  Footer map[string] string
  Dates []DateList
  Breakfasts []DateList
  Lunchs []DateList
  Dinners []DateList
}

type DateList struct {
  Title string
  Date string
  Key string
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

    // 日付一覧
    t := time.Now().In(time.FixedZone("Asia/Tokyo", 9 * 60 * 60))
    month := ""
    dates := []DateList{}
    for {
      key := ""
      ymd := t.Format(dateFormat)
      ym := t.Format(monthFormat)
      if month != ym {
        month = ym
        key = month
      }
      dates = append(dates, DateList { Title: key, Date: ymd, Key: t.Format("200601") })
      if ymd == startDate {
        break
      }
      t = t.AddDate(0, 0, -1)
    }

    issues := utils.ReadIssues(utils.Yaml().GitHub.Token, utils.Yaml().GitHub.User, utils.Yaml().GitHub.Project, 0)

    breakfasts := []DateList{}
    lunchs := []DateList{}
    dinners := []DateList{}

    for _, s := range issues {
      breakfastMessage := ""
      breakfastImage := ""
      breakfastMessageStart := false
      lunchMessage := ""
      lunchImage := ""
      lunchMessageStart := false
      dinnerMessage := ""
      dinnerImage := ""
      dinnerMessageStart := false
      otherMessage := ""
      otherMessageStart := false

      for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(*s.Body, -1) {
        if v == "### 朝食" {
          breakfastMessageStart = true
        }
        if v == "### 昼食" {
          breakfastMessageStart = false
          lunchMessageStart = true
        }
        if v == "### 夕食" {
          lunchMessageStart = false
          dinnerMessageStart = true
        }
        if v == "## その他感想的なもの" {
          dinnerMessageStart = false
          otherMessageStart = true
        }
        if breakfastMessageStart && v != "### 朝食" {
          if strings.HasPrefix(v, "<img src") {
            breakfastImage += v
            breakfastImage += "\n"
          } else if v != "" {
            breakfastMessage += v
            breakfastMessage += "\n"
          }
        }
        if lunchMessageStart && v != "### 昼食" {
          if strings.HasPrefix(v, "<img src") {
            lunchImage += v
            lunchImage += "\n"
          } else if v != "" {
            lunchMessage += v
            lunchMessage += "\n"
          }
        }
        if dinnerMessageStart && v != "### 夕食" {
          if strings.HasPrefix(v, "<img src") {
            dinnerImage += v
            dinnerImage += "\n"
          } else if v != "" {
            dinnerMessage += v
            dinnerMessage += "\n"
          }
        }
        if otherMessageStart && v != "## その他感想的なもの" {
            otherMessage += v
            otherMessage += "\n"
        }
      }

      fmt.Println(breakfastMessage)
      fmt.Println(lunchMessage)
      fmt.Println(dinnerMessage)
      fmt.Println(breakfastImage)
      fmt.Println(lunchImage)
      fmt.Println(dinnerImage)
    }

    data := IndexData {
      Footer: utils.Yaml().FooterLinks,
      Dates: dates,
      Breakfasts: breakfasts,
      Lunchs: lunchs,
      Dinners: dinners,
    }
    return e.Render(http.StatusOK, "index.html", data)
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
    session.Save()

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
