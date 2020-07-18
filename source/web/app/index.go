package app

import (
  "net/http"
  "github.com/labstack/echo"
  "source/utils"
  "source/structs"
  "time"
  "regexp"
  "strconv"
  "strings"
)

func IndexView(e echo.Context) error {

  if utils.IsNotLogined(e) {
    return e.Redirect(http.StatusFound, "/login")
  }

  // 日付一覧
  t := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
  month := ""
  dates := []structs.DateList{}
  months := []string{}
  for {
    key := ""
    ymd := t.Format(utils.DateFormat)
    ym := t.Format(utils.MonthFormat)
    if month != ym {
    month = ym
    key = month
    months = append(months, ym)
    }
    dates = append(dates, structs.DateList{Title: key, Date: ymd, Key: t.Format("200601")})
    if ymd == utils.StartDate {
    break
    }
    t = t.AddDate(0, 0, -1)
  }

  data := structs.IndexData{
    Footer: utils.Yaml().FooterLinks,
    Dates:  dates,
    Months: months,
  }

  return e.Render(http.StatusOK, "index.html", data)
  }

func IndexDetailJson(e echo.Context) error {
  if utils.IsNotLogined(e) {
    return e.JSON(http.StatusOK, map[string]string{
    "error": "login",
    })
  }

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

  return e.JSON(http.StatusOK, map[string]string{
    "title": title,
    "body":  body,
  })
  }

func IndexDataJson(e echo.Context) error {
  if utils.IsNotLogined(e) {
    return e.JSON(http.StatusOK, map[string]string{
    "error": "login",
    })
  }

  minusMonth, _ := strconv.Atoi(e.Param("minusMonth"))
  issues := utils.ReadIssues(utils.Yaml().GitHub.Token, utils.Yaml().GitHub.User, utils.Yaml().GitHub.Project, minusMonth)

  breakfasts := []structs.ContentsList{}
  lunchs := []structs.ContentsList{}
  dinners := []structs.ContentsList{}
  others := []structs.ContentsList{}

  wdays := [...]string{"日", "月", "火", "水", "木", "金", "土"}

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
      if strings.HasPrefix(v, "<img") {
      otherMessage += imgReplace.ReplaceAllString(v, "img class='lazyload' data-src")
      otherMessage += "<br>"
      } else if strings.HasPrefix(v, "##") {
      otherMessage += "\n#"
      otherMessage += v
      } else if strings.HasPrefix(v, "* ") || strings.HasPrefix(v, "- ") {
      otherMessage += "\n"
      otherMessage += v
      } else if v != "" {
      otherMessage += v
      otherMessage += "<br>"
      }
      otherMessage += "\n"
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

    date, _ := time.Parse(utils.DateFormat, *s.Title)
    dateString := *s.Title + "（" + wdays[date.Weekday()] + "）"

    breakfasts = append(breakfasts, structs.ContentsList{Date: dateString, Content: breakfastMessage, Image: breakfastImage})
    lunchs = append(lunchs, structs.ContentsList{Date: dateString, Content: lunchMessage, Image: lunchImage})
    dinners = append(dinners, structs.ContentsList{Date: dateString, Content: dinnerMessage, Image: dinnerImage})
    others = append(others, structs.ContentsList{Date: dateString, Content: otherMessage, Image: ""})
  }

  data := structs.IssuesData{
    Breakfasts: breakfasts,
    Lunchs:   lunchs,
    Dinners:  dinners,
    Others:   others,
  }

  return e.JSON(http.StatusOK, data)
  }