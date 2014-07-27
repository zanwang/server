package util

import (
  "os"
)

const (
  defaultEnv = "dev"
  envVarName = "GO_ENV"
)

var env string

func Environment() string {
  env = os.Getenv(envVarName)

  if env == "" {
    env = defaultEnv
    os.Setenv(envVarName, defaultEnv)
  }

  return env
}