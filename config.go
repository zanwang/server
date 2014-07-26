package main

import (
  "io/ioutil"
  "path/filepath"
  "log"
  "os"
  "gopkg.in/yaml.v1"
)

type config struct {
  Path string
  Data map[string]interface{}
}

const configDir = "config"
const defaultEnv = "dev"
const envVarName = "MARTINI_ENV"

var acceptExtnameForConfig = []string{"yml", "yaml"}
var Config = config{}

func init() {
  env := getenv()
  basedir := Basedir()

  for _, ext := range acceptExtnameForConfig {
    path := filepath.Join(basedir, configDir, env + "." + ext)

    if exists(path) {
      data, err := ioutil.ReadFile(path)
      check(err)

      err = yaml.Unmarshal(data, &Config.Data)
      check(err)

      Config.Path = path
      break
    }
  }

  log.Printf("Environment: %s", env)
  log.Printf("Config path: %s", Config.Path)
}

func getenv() string {
  env := os.Getenv(envVarName)

  if env == "" {
    env = defaultEnv
    os.Setenv(envVarName, defaultEnv)
  }

  return env
}

func Basedir() string {
  dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
  check(err)

  return dir
}

func exists(path string) bool {
  if _, err := os.Stat(path); err == nil {
    return true
  } else {
    return false
  }
}

func check(e error) {
  if e != nil {
    panic(e)
  }
}

func (c *config) Get(key string) interface{} {
  return c.Data[key]
}

func (c *config) Set(key string, value interface{}) {
  c.Data[key] = value
}