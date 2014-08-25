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
		Host   string `yaml:"host"`
		Port   int    `yaml:"port"`
		Secret string `yaml:"secret"`
	} `yaml:"server"`

	Database struct {
		Type string `yaml:"type"`
		Path string `yaml:"path"`
	} `yaml:"database"`

	ReservedDomains []string `yaml:"reserved_domains"`

	Mailgun struct {
		Domain     string `yaml:"domain"`
		PrivateKey string `yaml:"private_key"`
		PublicKey  string `yaml:"public_key"`
	} `yaml:"mailgun"`

	Facebook struct {
		AppID     string `yaml:"app_id"`
		AppSecret string `yaml:"app_secret"`
	} `yaml:"facebook"`

	Twitter struct {
		APIKey    string `yaml:"api_key"`
		APISecret string `yaml:"api_secret"`
	} `yaml:"twitter"`

	EmailActivation bool `yaml:"email_activation"`
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
