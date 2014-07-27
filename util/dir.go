package util

import (
  "path/filepath"
  "os"
)

var baseDir string

func Basedir() string {
  if baseDir != "" {
    return baseDir
  }

  dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

  if err != nil {
    panic(err)
  }

  baseDir = dir

  return baseDir
}

func ViewDir() string {
  return filepath.Join(Basedir(), "views")
}

func PublicDir() string {
  return filepath.Join(Basedir(), "public")
}

func ResolveView(name string) string {
  return filepath.Join(ViewDir(), name + ".html")
}