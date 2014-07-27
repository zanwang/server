package main

import (
  "io/ioutil"
  "path/filepath"
  "log"
  "os"
  "gopkg.in/yaml.v1"
  "github.com/tommy351/maji.moe/util"
)

type config struct {
  Path string
  Data map[string]interface{}
}

const configDir = "config"

var acceptExtnameForConfig = []string{"yml", "yaml"}
var Config = config{}

func init() {
  env := util.Environment()
  basedir := util.Basedir()

  for _, ext := range acceptExtnameForConfig {
    path := filepath.Join(basedir, configDir, env + "." + ext)

    if exists(path) {
      data, err := ioutil.ReadFile(path)

      if err != nil {
        panic(err)
        return
      }

      err = yaml.Unmarshal(data, &Config.Data)

      if err != nil {
        panic(err)
        return
      }

      Config.Path = path
      break
    }
  }

  log.Printf("Environment: %s", env)
  log.Printf("Config path: %s", Config.Path)
}

func exists(path string) bool {
  if _, err := os.Stat(path); err == nil {
    return true
  } else {
    return false
  }
}

func (c *config) Get(key string) interface{} {
  return c.Data[key]
}

func (c *config) Set(key string, value interface{}) {
  c.Data[key] = value
}