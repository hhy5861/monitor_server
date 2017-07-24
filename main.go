package main

import (
	"os"

	"github.com/hhy5861/logrus"
	"github.com/urfave/cli"
	"gitlab.pnlyy.com/monitor_server/service"
)

func main() {
	app := service.NewCli()
	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "run Commands",
			Action: func(c *cli.Context) {
				svc := service.New()
				svc.Run()
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		var ps logrus.Params
		logrus.Fatal(ps, err, "cli run error")
	}
}
