package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gitlab.pnlyy.com/monitor_server/model"
	"gitlab.pnlyy.com/monitor_server/utils"
	"io/ioutil"
	"net/http"
	"strings"
)

type Group struct {
	Name     string
	List_id  []string
	response utils.Response
}

func (g *Group) Add(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, err := ioutil.ReadAll(r.Body)
	if err == nil {
		var params Group
		json.Unmarshal(req, &params)
		listId := strings.Join(params.List_id, ",")
		if params.Name != "" && listId != "" {
			if err == nil {
				group := model.Group{
					Name:         params.Name,
					User_id_list: listId,
				}

				if res, err := group.Save(); err == nil {
					g.response.Code = 0
					g.response.Data = res
					g.response.Message = "添加成功"
					g.response.WriteJson(w)
					return
				}
			}
		}
	}

	g.response.Code = 600000
	g.response.Data = ""
	g.response.Message = "添加失败"
	g.response.WriteJson(w)
}

func (g *Group) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	group := model.Group{}
	if data, err := group.GetList(param); err == nil {
		g.response.Code = 0
		g.response.Data = data
		g.response.Message = "获取成功"
		g.response.WriteJson(w)
		return
	}

	g.response.Code = 600001
	g.response.Data = ""
	g.response.Message = "获取失败"
	g.response.WriteJson(w)
}

func (g *Group) Delete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, err := ioutil.ReadAll(r.Body)
	if err == nil {
		user := model.User{}
		json.Unmarshal(req, &user)
		if user.Id > 0 {
			if user.Delete(user.Id) {
				g.response.Code = 0
				g.response.Data = ""
				g.response.Message = "删除成功"
				g.response.WriteJson(w)
				return
			}
		}
	}

	g.response.Code = 600002
	g.response.Data = ""
	g.response.Message = "删除失败"
	g.response.WriteJson(w)
}
