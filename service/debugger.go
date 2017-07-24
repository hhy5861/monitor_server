package service

import (
	"fmt"
	"github.com/robfig/cron"
	"gitlab.pnlyy.com/monitor_server/model"
	"gitlab.pnlyy.com/monitor_server/monitor"
	"time"
)

type Debugger struct {
	DebugTime    time.Duration
	MonDebugTime time.Duration
}

func (d *Debugger) debugger() {
	specStrategy := fmt.Sprintf("*/%d * * * *", d.DebugTime)
	specCheck := fmt.Sprintf("*/%d * * * *", d.MonDebugTime)

	crond := cron.New()
	crond.AddFunc(specStrategy, func() {
		strategy := model.MonitorConfig{}
		strategy.GetStrategyAll()
	})

	crond.AddFunc(specCheck, func() {
		monitor.Select(model.StrategyArray)
	})

	crond.Start()
}
