package model

import (
	"encoding/json"
	"sync"
	"time"
)

var (
	mux           sync.Mutex
	StrategyArray []*Strategy
)

type (
	MonitorConfig struct {
		Id         uint   `json:"id"`
		Project_id uint   `json:"project_id"`
		Token      string `json:"token"`
		Module     string `json:"module"`
		Point      string `json:"point"`
		Strategy   string `json:"strategy"`
		Conf_type  uint   `json:"conf_type"`
		Ctime      string `json:"ctime"`
		Utime      string `json:"utime"`
		Valid      int    `json:"valid"`
	}

	Strategy struct {
		Token      string `json:"token"`
		Module     string `json:"module"`
		Point      string `json:"point"`
		Op         string `json:"op"`
		Opo        string `json:"opo"`
		Field1     string `json:"field1"`
		Field2     string `json:"field2"`
		Field3     string `json:"field3"`
		Field4     string `json:"field4"`
		Level      int    `json:"level"`
		Limit      int    `json:"limit"`
		Range      int    `json:"range"`
		Type       uint   `json:"type"`
		FieldRange int    `json:"field_range"`
		Group_id   int    `json:"group_id"`
	}
)

func (m *MonitorConfig) TableName() string {
	return "mk_monitor_config"
}

func (m *MonitorConfig) Save() (*MonitorConfig, error) {
	if err := db.Save(m).Error; err != nil {
		return nil, err
	}

	return m, nil
}

func (m *MonitorConfig) GetList(param map[string]interface{}) ([]*MonitorConfig, error) {
	var list []*MonitorConfig
	err := db.Where("valid = ? AND project_id = ?", "0", param["id"]).Offset(param["page"]).Limit(param["size"]).Order("id DESC").Find(&list).Error

	return list, err
}

func (m *MonitorConfig) Delete(id uint) bool {
	m.Id = id
	if err := db.Where("valid = ?", "0").First(m).Error; err == nil && m.Id > 0 {
		m.Valid = 1
		m.Utime = time.Now().Format("2006-01-02 15:04:05")
		if err := db.Save(m).Error; err == nil {
			return true
		}
	}

	return false
}

func (m *MonitorConfig) GetStrategyAll() ([]*Strategy, error) {
	var list []*MonitorConfig
	var Array []*Strategy

	err := db.Where("valid = ?", "0").Find(&list).Error
	if err == nil {
		mux.Lock()
		defer mux.Unlock()
		for _, v := range list {
			var data *Strategy
			json.Unmarshal([]byte(v.Strategy), &data)
			data.Token = v.Token
			data.Type = v.Conf_type

			Array = append(Array, data)
		}
	}

	StrategyArray = Array

	return StrategyArray, err
}
