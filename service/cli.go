package service

import (
	"github.com/urfave/cli"
)

var (
	host          string
	port          int
	debug         bool
	logPath       string
	configPath    string
	elasticAddres string
)

func NewCli() *cli.App {
	app := cli.NewApp()
	app.Name = "vip ractice monitor server"
	app.Usage = "vip ractice monitor server"
	app.Author = "Mike Huang"
	app.Email = "huanghaiying@vipractice.cn"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config-path",
			Value:       string("/etc/monitor_servce"),
			Destination: &configPath,
			Usage:       "config file path",
		},

		cli.StringFlag{
			Name:        "log-path",
			Value:       string("/data/logs/monitor_server"),
			Destination: &logPath,
			Usage:       "logs path distribution",
		},

		cli.StringFlag{
			Name:        "elastic-addres",
			Value:       string("http://192.168.40.222:9200"),
			Destination: &elasticAddres,
			Usage:       "config file path",
		},

		cli.StringFlag{
			Name:        "host",
			Value:       string("0.0.0.0"),
			Destination: &host,
			Usage:       "api http listen host",
		},

		cli.IntFlag{
			Name:        "port",
			Value:       int(8080),
			Destination: &port,
			Usage:       "api http listen port",
		},

		cli.BoolFlag{
			Name:        "debug",
			Destination: &debug,
			Usage:       "server debug logs info",
		},
	}

	return app
}
