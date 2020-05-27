package utils

import (
  "github.com/google/go-github/github"
  "golang.org/x/oauth2"
  "log"
  "sync"
  "time"
  "context"
)

var gclient *github.Client
var onceClient sync.Once

const dateFormat = "2020/05/23"

func ReadIssues(token string, user string, project string) []*github.Issue {
  issues, _, err := client(token).Issues.ListByRepo(context.Background(), user, project, nil)
  if err != nil {
    log.Fatalln(err)
  }
  return issues
}

func CreateNewTodaysIssue(token string, user string, project string) {
  opt := &github.IssueRequest {
	  Title: github.String(time.Now().Format(dateFormat)),
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
  client(token).Issues.Create(context.Background(), user, project, opt)
}

func client(token string) *github.Client {
  onceClient.Do(func() {
    ts := oauth2.StaticTokenSource(
      &oauth2.Token{AccessToken: token},
    )
    tc := oauth2.NewClient(oauth2.NoContext, ts)
    gclient = github.NewClient(tc)
  })
  return gclient
}