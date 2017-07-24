package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gitlab.pnlyy.com/monitor_server/model"
	"gitlab.pnlyy.com/monitor_server/utils"
	"io/ioutil"
	"net/http"
	"time"
)

type (
	MonitorConfig struct {
		response utils.Response
	}

	Monitor struct {
		Id         uint      `json:"id"`
		Project_id uint      `json:"project_id"`
		Token      string    `json:"token"`
		Module     string    `json:"module"`
		Point      string    `json:"point"`
		Strategy   *Strategy `json:"strategy"`
		Conf_type  uint      `json:"conf_type"`
		Conf_name  string    `json:"conf_name"`
		Ctime      string    `json:"ctime"`
		Utime      string    `json:"utime"`
		Valid      int       `json:"valid"`
	}

	Strategy struct {
		Module   string `json:"module"`
		Point    string `json:"point"`
		Op       string `json:"op"`
		Opo      string `json:"opo"`
		Field1   string `json:"field1"`
		Field2   string `json:"field2"`
		Field3   string `json:"field3"`
		Field4   string `json:"field4"`
		Level    int    `json:"level"`
		Limit    int    `json:"limit"`
		Range    int    `json:"range"`
		Group_id int    `json:"group_id"`
	}
)

func (m *MonitorConfig) Add(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, err := ioutil.ReadAll(r.Body)
	if err == nil {
		var data Monitor
		json.Unmarshal(req, &data)
		strategy, err := json.Marshal(data.Strategy)
		if data.Project_id > 0 && err == nil {
			nowTime := time.Now().Format("2006-01-02 15:04:05")
			conf := model.MonitorConfig{
				Project_id: data.Project_id,
				Token:      data.Token,
				Module:     data.Module,
				Point:      data.Point,
				Strategy:   string(strategy),
				Conf_type:  data.Conf_type,
				Ctime:      nowTime,
				Utime:      nowTime,
				Valid:      0,
			}

			if res, err := conf.Save(); err == nil {
				m.response.Code = 0
				m.response.Data = res
				m.response.Message = "添加配置成功"
				m.response.WriteJson(w)
				return
			}
		}
	}

	m.response.Code = 400000
	m.response.Data = ""
	m.response.Message = "添加配置失败"
	m.response.WriteJson(w)
}

func (m *MonitorConfig) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	param := make(map[string]interface{})
	param["page"] = 0
	param["size"] = 200

	page := r.URL.Query().Get("page")
	if page != "" {
		param["page"] = page
	}

	size := r.URL.Query().Get("size")
	if size != "" {
		param["size"] = size
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		m.response.Code = 400003
		m.response.Data = ""
		m.response.Message = "项目id不能为空"
		m.response.WriteJson(w)
		return
	}

	param["id"] = id

	conf := model.MonitorConfig{}
	res, err := conf.GetList(param)
	if err == nil {
		var result []*Monitor
		for _, v := range res {
			var data Monitor
			var strategy Strategy
			json.Unmarshal([]byte(v.Strategy), &strategy)

			data.Id = v.Id
			data.Project_id = v.Project_id
			data.Module = v.Module
			data.Point = v.Point
			data.Token = v.Token
			data.Strategy = &strategy
			data.Conf_type = v.Conf_type
			data.Ctime = v.Ctime
			data.Utime = v.Utime
			data.Valid = v.Valid

			switch data.Conf_type {
			case 1:
				data.Conf_name = "单配置"
				break
			case 2:
				data.Conf_name = "双配置"
				break

			}

			result = append(result, &data)
		}

		m.response.Code = 0
		m.response.Data = result
		m.response.Message = "获取成功"
		m.response.WriteJson(w)
		return
	}

	m.response.Code = 400002
	m.response.Data = ""
	m.response.Message = "获取失败"
	m.response.WriteJson(w)
}

func (m *MonitorConfig) Del(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, err := ioutil.ReadAll(r.Body)
	if err == nil {
		config := model.MonitorConfig{}
		json.Unmarshal(req, &config)
		if config.Id > 0 {
			if ok := config.Delete(config.Id); ok {
				m.response.Code = 0
				m.response.Data = ""
				m.response.Message = "删除成功"
				m.response.WriteJson(w)
				return
			}
		}
	}

	m.response.Code = 400001
	m.response.Data = ""
	m.response.Message = "删除失败"
	m.response.WriteJson(w)
}
