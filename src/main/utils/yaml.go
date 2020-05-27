package utils

import (
  "fmt"
  yml "gopkg.in/yaml.v2"
  "io/ioutil"
  "log"
  "sync"
)

type yaml struct {
  GitHub githubs `yaml:"github"`
}

type githubs struct {
  Token string `yaml:"token"`
  User string `yaml:"user"`
  Project string `yaml:"project"`
}

var instanceYaml *yaml
var onceYaml sync.Once

func Yaml() *yaml {
  onceYaml.Do(func() {
    initializeYaml()
  })
  return instanceYaml
}

func initializeYaml() {
  fmt.Println("start loading yaml file")
  buf, err := ioutil.ReadFile("secret.yaml")
  if err != nil {
    log.Fatalln(err)
  }
  instanceYaml = &yaml{}
  err = yml.Unmarshal(buf, instanceYaml)
  if err != nil {
    log.Fatalln(err)
  }
  fmt.Println("finished loading yaml file")
}