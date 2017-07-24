package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gitlab.pnlyy.com/monitor_server/model"
	"gitlab.pnlyy.com/monitor_server/utils"
	"io/ioutil"
	"net/http"
)

type Email struct {
	response utils.Response
}

func (e *Email) Add(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, err := ioutil.ReadAll(r.Body)
	if err == nil {
		email := model.EmailTemplate{}
		json.Unmarshal(req, &email)
		if email.Name != "" {
			res, err := email.FindByName(email.Name)
			if err != nil && res.Id == 0 {
				if res, err := email.Add(); err == nil {
					e.response.Code = 0
					e.response.Data = res
					e.response.Message = "添加模板成功"
					e.response.WriteJson(w)
					return
				}
			} else {
				e.response.Code = 500001
				e.response.Data = ""
				e.response.Message = "模板名称存在"
				e.response.WriteJson(w)
				return
			}
		}
	}

	e.response.Code = 500000
	e.response.Data = ""
	e.response.Message = "添加模板失败"
	e.response.WriteJson(w)
}

func (e *Email) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	param := make(map[string]interface{})
	param["page"] = 0
	param["size"] = 200

	page := ps.ByName("page")
	if page != "" {
		param["page"] = page
	}

	size := ps.ByName("size")
	if size != "" {
		param["size"] = size
	}

	email := model.EmailTemplate{}
	if data, err := email.GetList(param); err == nil {
		e.response.Code = 0
		e.response.Data = data
		e.response.Message = "获取列表成功"
		e.response.WriteJson(w)
		return
	}

	e.response.Code = 500002
	e.response.Data = ""
	e.response.Message = "获取列表失败"
	e.response.WriteJson(w)
}

func (e *Email) Delete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, err := ioutil.ReadAll(r.Body)
	if err == nil {
		email := model.EmailTemplate{}
		json.Unmarshal(req, &email)
		if email.Id > 0 {
			if email.Delete(email.Id) {
				e.response.Code = 0
				e.response.Data = ""
				e.response.Message = "删除成功"
				e.response.WriteJson(w)
				return
			}
		}
	}

	e.response.Code = 500003
	e.response.Data = ""
	e.response.Message = "删除失败"
	e.response.WriteJson(w)
}
