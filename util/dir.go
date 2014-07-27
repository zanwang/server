package util

import (
  "path/filepath"
  "os"
)

func Basedir() string {
  dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

  if err != nil {
    panic(err)
  }

  return dir
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