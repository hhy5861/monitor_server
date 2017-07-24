package service

import (
	"github.com/hhy5861/logrus"
	"github.com/julienschmidt/httprouter"
	"gitlab.pnlyy.com/monitor_server/model"
	"gitlab.pnlyy.com/monitor_server/redis"
	"gitlab.pnlyy.com/monitor_server/config"
)

type Service struct {
	Host   string
	Post   int
	router *httprouter.Router
}

func New() *Service {
	svc := &Service{}
	return svc
}

func (svc *Service) Run() {
	lrs := logrus.New(logPath, "monitor_server.log")
	lrs.Init()

	c := config.NewConfig(host, logPath, elasticAddres, port, debug)
	c.ParseConfig(configPath)

	conf := config.Config

	redis.NewRedis(conf.Redis.Host,
		conf.Redis.Port,
		conf.Redis.Db,
		conf.Redis.Password)

	model.NewMysql(conf.Mysql.Host,
		conf.Mysql.Dbname,
		conf.Mysql.User,
		conf.Mysql.Password,
		conf.Mysql.Charset,
		conf.Mysql.Dialect,
		conf.Mysql.Port,
		conf.Servce.Debug)

	monitorConfig := model.MonitorConfig{}
	monitorConfig.GetStrategyAll()

	/*monitor.New(elasticAddres)
	go func() {
		debugger := Debugger{
			DebugTime:    30,
			MonDebugTime: 2,
		}

		debugger.debugger()
	}()*/

	svc.Host = conf.Servce.Host
	svc.Post = conf.Servce.Port
	svc.run()
}
