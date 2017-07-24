package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gitlab.pnlyy.com/monitor_server/model"
	"gitlab.pnlyy.com/monitor_server/utils"
	"io/ioutil"
	"net/http"
)

type User struct {
	response utils.Response
}

func (u *User) Add(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, err := ioutil.ReadAll(r.Body)
	if err == nil {
		user := model.User{}
		json.Unmarshal(req, &user)
		if user.Name != "" {
			if res, err := user.Save(); err == nil {
				u.response.Code = 0
				u.response.Data = res
				u.response.Message = "添加用户成功"
				u.response.WriteJson(w)
				return
			}
		}
	}

	u.response.Code = 300000
	u.response.Data = ""
	u.response.Message = "添加用户失败"
	u.response.WriteJson(w)
}

func (u *User) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	user := model.User{}
	if data, err := user.GetList(param); err == nil {
		u.response.Code = 0
		u.response.Data = data
		u.response.Message = "获取成功"
		u.response.WriteJson(w)
		return
	}

	u.response.Code = 300001
	u.response.Data = ""
	u.response.Message = "获取失败"
	u.response.WriteJson(w)
}

func (u *User) Delete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, err := ioutil.ReadAll(r.Body)
	if err == nil {
		user := model.User{}
		json.Unmarshal(req, &user)
		if user.Id > 0 {
			if user.Delete(user.Id) {
				u.response.Code = 0
				u.response.Data = ""
				u.response.Message = "删除成功"
				u.response.WriteJson(w)
				return
			}
		}
	}

	u.response.Code = 300002
	u.response.Data = ""
	u.response.Message = "删除失败"
	u.response.WriteJson(w)
}
