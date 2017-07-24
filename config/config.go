package config

import (
	"io/ioutil"
	"os"

	"github.com/hhy5861/logrus"
	"gitlab.pnlyy.com/monitor_server/model"
	"gopkg.in/yaml.v2"
)

type (
	Conf struct {
		Email  *Email
		Redis  *Redis
		Mysql  *model.Mysql
		Servce *Servce
	}

	Redis struct {
		Host     string `yaml:"host"`
		Port     int64  `yaml:"port"`
		Password string `yaml:"password"`
		Db       int    `yaml:"db"`
	}

	Email struct {
		From     string `yaml:"from"`
		Smtp     string `yaml:"smtp"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	}

	Servce struct {
		Host           string `yaml:"host"`
		Port           int    `yaml:"port"`
		Debug          bool   `yaml:"debug"`
		Log_path       string `yaml:"log_path"`
		Elastic_addres string `yaml:"elastic_addres"`
	}
)

var (
	Config *Conf
)

func NewConfig(host, logPath, elasticAddres string, port int, debug bool) *Conf {
	conf := &Conf{
		Servce: &Servce{
			Host:           host,
			Port:           port,
			Debug:          debug,
			Log_path:       logPath,
			Elastic_addres: elasticAddres,
		},
	}

	return conf
}

func (c *Conf) ParseConfig(configFile string) {
	configFile = configFile + "/config.yaml"
	if fileExists(configFile) {
		body, err := ioutil.ReadFile(configFile)
		if err != nil {
			var ps logrus.Params
			logrus.Fatal(ps, err, "config file read err")
		}

		err = yaml.Unmarshal(body, c)
		if err != nil {
			var ps logrus.Params
			logrus.Fatal(ps, err, "config unmarshal err")
		}

		Config = c
	} else {
		var ps logrus.Params
		logrus.Fatal(ps, nil, "config file does not exist")
	}
}

func fileExists(configFile string) bool {
	if _, err := os.Stat(configFile); !os.IsExist(err) {
		return true
	}

	return false
}
