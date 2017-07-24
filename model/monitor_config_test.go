package model

import "testing"

func TestNewMysql(t *testing.T) {
	host, dbname, user, password, charset, dialect := "127.0.0.1", "monitor_servce", "mike", "tlslpc", "utf8", "mysql"
	post := 3306
	NewMysql(host, dbname, user, password, charset, dialect, post)
}

func TestMonitorConfig_GetStrategyAll(t *testing.T) {
	config := MonitorConfig{}
	data, err := config.GetStrategyAll()
	if err == nil {
		for _, v := range data {
			t.Logf("%#v", v)
		}
	}
}
