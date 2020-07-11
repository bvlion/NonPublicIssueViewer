package structs

type LoginParams struct {
  Passphrase string `json:"passphrase"`
}

type IndexData struct {
  Footer map[string] string
  Dates []DateList
  Months []string
}

type DateList struct {
  Title string
  Date string
  Key string
}

type ContentsList struct {
  Date string
  Content string
  Image string
}

type IssuesData struct {
  Breakfasts []ContentsList
  Lunchs []ContentsList
  Dinners []ContentsList
  Others []ContentsList
}