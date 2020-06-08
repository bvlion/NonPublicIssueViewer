package utils

import (
  yml "gopkg.in/yaml.v2"
  "io/ioutil"
  "sync"
  "source/log"
)

type yaml struct {
  GitHub githubs `yaml:"github"`
  Passphrase string `yaml:"login_pass"`
  SessionKey string `yaml:"session_key"`
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
  log.InfoLog("start loading yaml file")
  buf, err := ioutil.ReadFile("secret.yaml")
  if err != nil {
    log.ErrorLog(err)
  }
  instanceYaml = &yaml{}
  err = yml.Unmarshal(buf, instanceYaml)
  if err != nil {
    log.ErrorLog(err)
  }
  log.InfoLog("finished loading yaml file")
}