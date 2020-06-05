package utils

import (
  "github.com/google/go-github/github"
  "golang.org/x/oauth2"
  "sync"
  "time"
  "fmt"
  "context"
  "source/log"
)

var gclient *github.Client
var onceClient sync.Once

const dateFormat = "2006/01/02"

func ReadIssues(token string, user string, project string) []*github.Issue {
  issues, _, err := client(token).Issues.ListByRepo(context.Background(), user, project, nil)
  if err != nil {
    log.ErrorLog(err)
  }
  return issues
}

func CreateNewTodaysIssue(token string, user string, project string) {
  opt := &github.IssueRequest {
	  Title: github.String(time.Now().In(time.FixedZone("Asia/Tokyo", 9 * 60 * 60)).Format(dateFormat)),
	  Body:  github.String(`## やったこと

* 
* 

## 食事

### 朝食


<img src="" width="320"/>

### 昼食


<img src="" width="320"/>

### 夕食


<img src="" width="320"/>

## その他感想的なもの`),
  }
  _, res, err := client(token).Issues.Create(context.Background(), user, project, opt)
  if err != nil {
    log.ErrorLog(err)
  }
  if res != nil {
    log.InfoLog(fmt.Sprintf("%#v", res))
  }
}

func client(token string) *github.Client {
  onceClient.Do(func() {
    log.InfoLog("start loading GitHub setting")
    ts := oauth2.StaticTokenSource(
      &oauth2.Token{AccessToken: token},
    )
    tc := oauth2.NewClient(oauth2.NoContext, ts)
    gclient = github.NewClient(tc)
    log.InfoLog("finished loading GitHub setting")
  })
  return gclient
}