package utils

import (
  "time"
  "source/structs"
  "regexp"
  "strconv"
  "strings"
)

const DateFormat = "2006/01/02"
const MonthFormat = "2006/01"

const startDate = "2020/05/13"

func CreateDateMonths() ([]structs.DateList, []string) {
  t := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
  month := ""
  dates := []structs.DateList{}
  months := []string{}

  for {
    key := ""
    ymd := t.Format(DateFormat)
    ym := t.Format(MonthFormat)

    if month != ym {
      month = ym
      key = month
      months = append(months, ym)
    }

    dates = append(dates, structs.DateList{Title: key, Date: ymd, Key: t.Format("200601")})
    
    if ymd == startDate {
      break
    }
    
    t = t.AddDate(0, 0, -1)
  }

  return dates, months
}

func CreateDetailBody(title string) string {
  issues := ReadOneIssue(
    Yaml().GitHub.Token,
    Yaml().GitHub.User,
    Yaml().GitHub.Project,
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

  return body
}

func CreateIssueData(param string) structs.IssuesData {
  minusMonth, _ := strconv.Atoi(param)
  issues := ReadIssues(Yaml().GitHub.Token, Yaml().GitHub.User, Yaml().GitHub.Project, minusMonth)

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
          lunchImage += imgReplace.ReplaceAllString(v, "img class='lazyload' data-src")
          lunchImage += "\n"
        } else if v != "" {
          lunchMessage += v
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

    date, _ := time.Parse(DateFormat, *s.Title)
    dateString := *s.Title + "（" + wdays[date.Weekday()] + "）"

    breakfasts = append(breakfasts, structs.ContentsList{Date: dateString, Content: breakfastMessage, Image: breakfastImage})
    lunchs = append(lunchs, structs.ContentsList{Date: dateString, Content: lunchMessage, Image: lunchImage})
    dinners = append(dinners, structs.ContentsList{Date: dateString, Content: dinnerMessage, Image: dinnerImage})
    others = append(others, structs.ContentsList{Date: dateString, Content: otherMessage, Image: ""})
  }

  return structs.IssuesData {
    Breakfasts: breakfasts,
    Lunchs:   lunchs,
    Dinners:  dinners,
    Others:   others,
  }
}
