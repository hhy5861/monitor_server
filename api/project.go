package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gitlab.pnlyy.com/monitor_server/model"
	"gitlab.pnlyy.com/monitor_server/utils"
	"io/ioutil"
	"net/http"
)

type Project struct {
	Name     string
	response utils.Response
}

func (p *Project) Add(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, err := ioutil.ReadAll(r.Body)
	if err == nil {
		project := model.Project{}
		json.Unmarshal(req, &project)
		if project.Name != "" {
			res, err := project.FindByName(project.Name)
			if err != nil && res.Id == 0 {
				if project.Token == "" {
					project.Token = string(utils.Krand(32, utils.KC_RAND_KIND_LOWER))
				}

				if res, err := project.SaveProject(); err == nil {
					p.response.Code = 0
					p.response.Data = res
					p.response.Message = "添加项目成功"
					p.response.WriteJson(w)
					return
				}
			} else {
				p.response.Code = 200001
				p.response.Data = ""
				p.response.Message = "项目已经存在"
				p.response.WriteJson(w)
				return
			}
		}
	}

	p.response.Code = 200000
	p.response.Data = ""
	p.response.Message = "添加项目失败"
	p.response.WriteJson(w)
}

func (p *Project) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	project := model.Project{}
	res, err := project.GetList(param)
	if err == nil {
		p.response.Code = 0
		p.response.Data = res
		p.response.Message = "获取成功"
		p.response.WriteJson(w)
		return
	}

	p.response.Code = 200002
	p.response.Data = res
	p.response.Message = "获取失败"
	p.response.WriteJson(w)
}

func (p *Project) Del(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		project := model.Project{}
		json.Unmarshal(body, &project)
		if project.Id > 0 {
			if project.Delete(project.Id) {
				p.response.Code = 0
				p.response.Data = ""
				p.response.Message = "删除成功"
				p.response.WriteJson(w)
				return
			}
		}
	}

	p.response.Code = 200003
	p.response.Data = ""
	p.response.Message = "删除失败"
	p.response.WriteJson(w)
}
