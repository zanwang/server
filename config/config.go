package config

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v1"
)

// Config stores global configuration
type config struct {
	Server struct {
		Host   string `yaml:"host"`
		Port   int    `yaml:"port"`
		Secret string `yaml:"secret"`
		Logger bool   `yaml:"logger"`
	} `yaml:"server"`

	Database struct {
		Type string `yaml:"type"`
		Path string `yaml:"path"`
	} `yaml:"database"`

	Mailgun struct {
		Domain     string `yaml:"domain"`
		PrivateKey string `yaml:"private_key"`
		PublicKey  string `yaml:"public_key"`
	} `yaml:"mailgun"`

	Facebook struct {
		AppId     string `yaml:"app_id"`
		AppSecret string `yaml:"app_secret"`
	} `yaml:"facebook"`

	EmailActivation bool `yaml:"email_activation"`
}

const (
	configDir   = "config"
	Development = "dev"
	Production  = "prod"
	Test        = "test"
)

var (
	configExtname = []string{"yml", "yaml"}
	BaseDir, _    = os.Getwd()
	Config        config
	Env           string
)

// Load loads configuration file
func init() {
	if Env = os.Getenv("GO_ENV"); Env == "" {
		Env = Development
	}

	switch Env {
	case Production:
		gin.SetMode(gin.ReleaseMode)
	case Test:
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	if strings.HasPrefix(Env, "test_") {
		gin.SetMode(gin.TestMode)
	}

	for _, ext := range configExtname {
		path := path.Join(BaseDir, configDir, Env+"."+ext)

		if !exists(path) {
			continue
		}

		data, err := ioutil.ReadFile(path)

		if err != nil {
			panic(err)
		}

		if err = yaml.Unmarshal(data, &Config); err != nil {
			panic(err)
		}

		break
	}
}

func exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}
