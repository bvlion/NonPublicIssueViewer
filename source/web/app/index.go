package app

import (
  "net/http"
  "github.com/labstack/echo"
  "source/utils"
  "source/structs"
)

func IndexView(e echo.Context) error {
  if utils.IsNotLogined(e) {
    return e.Redirect(http.StatusFound, "/login")
  }

  dates, months := utils.CreateDateMonths()

  return e.Render(http.StatusOK, "index.html", structs.IndexData {
    Footer: utils.Yaml().FooterLinks,
    Dates:  dates,
    Months: months,
  })
}

func IndexDetailJson(e echo.Context) error {
  if utils.IsNotLogined(e) {
    return e.JSON(http.StatusOK, map[string]string {
      "error": "login",
    })
  }

  title := e.Param("year") + "/" + e.Param("month") + "/" + e.Param("day")
  body := utils.CreateDetailBody(title)

  return e.JSON(http.StatusOK, map[string]string {
    "title": title,
    "body":  body,
  })
}

func IndexDataJson(e echo.Context) error {
  if utils.IsNotLogined(e) {
    return e.JSON(http.StatusOK, map[string]string {
      "error": "login",
    })
  }

  return e.JSON(http.StatusOK, utils.CreateIssueData(e.Param("minusMonth")))
}