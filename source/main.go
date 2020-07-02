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
  "strconv"
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
}

type DateList struct {
  Title string
  Date string
  Key string
}

type IssuesData struct {
  Breakfasts []DateList
  Lunchs []DateList
  Dinners []DateList
  Others []DateList
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
    templates: template.Must(template.New("").Funcs(funcMap).ParseGlob("views/*.html")),
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

    data := IndexData {
      Footer: utils.Yaml().FooterLinks,
      Dates: dates,
    }

    return e.Render(http.StatusOK, "index.html", data)
  })

  e.GET("/detail/:year/:month/:day", func (e echo.Context) error {
    title := e.Param("year") + "/" + e.Param("month") + "/" + e.Param("day")
    issues := utils.ReadOneIssue(
      utils.Yaml().GitHub.Token,
      utils.Yaml().GitHub.User,
      utils.Yaml().GitHub.Project,
      title,
    )

    body := ""
    imgReplace := regexp.MustCompile("img.*src")

    for _, s := range issues {
      for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(*s.Body, -1) {
        if strings.HasPrefix(v, "<img") {
          body += imgReplace.ReplaceAllString(v, "img class='lazyload' data-src")
          body += "<br>"
        } else if strings.HasPrefix(v, "##") {
          body += "\n#"
          body += v
        } else if strings.HasPrefix(v, "* ") || strings.HasPrefix(v, "- ") {
          body += "\n"
          body += v
        } else if v != "" {
          body += v
          body += "<br>"
        }
        body += "\n"
      }
    }

    return e.JSON(http.StatusOK, map [string] string {
      "title": title,
      "body": body,
    })
  })

  e.GET("/issues/:minusMonth", func (e echo.Context) error {
    minusMonth, _ := strconv.Atoi(e.Param("minusMonth"))
    issues := utils.ReadIssues(utils.Yaml().GitHub.Token, utils.Yaml().GitHub.User, utils.Yaml().GitHub.Project, minusMonth)

    breakfasts := []DateList{}
    lunchs := []DateList{}
    dinners := []DateList{}
    others := []DateList{}

    wdays := [...] string{ "日", "月", "火", "水", "木", "金", "土" }

    imgReplace := regexp.MustCompile("img.*src")

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
          if strings.HasPrefix(v, "<img") {
            breakfastImage += imgReplace.ReplaceAllString(v, "img class='lazyload' data-src")
            breakfastImage += "\n"
          } else if v != "" {
            breakfastMessage += v
            breakfastMessage += "\n"
          }
        }
        if lunchMessageStart && v != "### 昼食" {
          if strings.HasPrefix(v, "<img") {
            lunchImage += v
            lunchImage += "\n"
          } else if v != "" {
            lunchMessage += imgReplace.ReplaceAllString(v, "img class='lazyload' data-src")
            lunchMessage += "\n"
          }
        }
        if dinnerMessageStart && v != "### 夕食" {
          if strings.HasPrefix(v, "<img") {
            dinnerImage += imgReplace.ReplaceAllString(v, "img class='lazyload' data-src")
            dinnerImage += "\n"
          } else if v != "" {
            dinnerMessage += v
            dinnerMessage += "\n"
          }
        }
        if otherMessageStart && v != "## その他感想的なもの" {
            otherMessage += imgReplace.ReplaceAllString(v, "img class='lazyload' data-src")
            otherMessage += "<br>"
        }
      }

      if breakfastMessage == "" {
        breakfastMessage = "未記入"
      }
      if lunchMessage == "" {
        lunchMessage = "未記入"
      }
      if dinnerMessage == "" {
        dinnerMessage = "未記入"
      }

      date, _ := time.Parse(dateFormat, *s.Title)
      dateString := *s.Title + "（" + wdays[date.Weekday()] + "）"

      breakfasts = append(breakfasts, DateList { Title: breakfastMessage, Date: dateString, Key: breakfastImage })
      lunchs = append(lunchs, DateList { Title: lunchMessage, Date: dateString, Key: lunchImage })
      dinners = append(dinners, DateList { Title: dinnerMessage, Date: dateString, Key: dinnerImage })
      others = append(others, DateList { Title: otherMessage, Date: dateString, Key: "" })
    }

    data := IssuesData {
      Breakfasts: breakfasts,
      Lunchs: lunchs,
      Dinners: dinners,
      Others: others,
    }
    
    return e.JSON(http.StatusOK, data)
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
