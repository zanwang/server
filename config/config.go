package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-martini/martini"
	"gopkg.in/yaml.v1"
)

// Config stores global configuration
type Config struct {
	Server struct {
		Host   string
		Port   int
		Secret string
	}

	Database struct {
		Type string
		Path string
	}
}

const configDir = "config"

var acceptExtnameForConfig = []string{"yml", "yaml"}
var baseDir, _ = os.Getwd()

// Load loads configuration file
func Load() *Config {
	env := martini.Env
	obj := Config{}

	for _, ext := range acceptExtnameForConfig {
		path := filepath.Join(baseDir, configDir, env+"."+ext)

		if exists(path) {
			data, err := ioutil.ReadFile(path)

			if err != nil {
				panic(err)
			}

			if err = yaml.Unmarshal(data, &obj); err != nil {
				panic(err)
			}

			break
		}
	}

	return &obj
}

func exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}
