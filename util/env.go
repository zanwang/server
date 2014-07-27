package util

import (
  "os"
)

const defaultEnv = "dev"
const envVarName = "GO_ENV"

func Environment() string {
  env := os.Getenv(envVarName)

  if env == "" {
    env = defaultEnv
  }

  return env
}