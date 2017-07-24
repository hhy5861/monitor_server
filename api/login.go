package api

import (
	"encoding/json"
	"fmt"
	"github.com/hhy5861/token"
	"github.com/julienschmidt/httprouter"
	"gitlab.pnlyy.com/monitor_server/model"
	"gitlab.pnlyy.com/monitor_server/redis"
	"gitlab.pnlyy.com/monitor_server/utils"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	EXPIRETIME = 3600 * 24 * 365
)

type (
	Login struct {
		UserName   string
		Password   string
		RememberMe bool
		response   utils.Response
	}

	CheckLogin struct {
		Token string
	}
)

func (login *Login) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		login.response.Code = 100000
		login.response.Data = err
		login.response.Message = "读取body出错了"
		login.response.WriteJson(w)
		return
	}

	var body Login
	err = json.Unmarshal(req, &body)
	if err != nil {
		login.response.Code = 100001
		login.response.Data = err
		login.response.Message = "解析json出错了"
		login.response.WriteJson(w)
		return
	}

	r.Body.Close()

	if body.UserName != "" {
		var users model.User
		if userInfo, err := users.GetUserByName(body.UserName); err == nil {
			password := utils.GetMD5Hash(body.Password)
			if strings.EqualFold(password, userInfo.Password) {
				login.response.Code = 0

				genNum := token.GenerateToken(4)
				tk := fmt.Sprintf("%s|%d", genNum, userInfo.Id)
				proToken := token.Encrypt([]byte(tk))
				var expire int
				if body.RememberMe {
					expire = EXPIRETIME
				} else {
					expire = 3600
				}

				userInfo.Token = proToken
				userJson, err := json.Marshal(userInfo)
				if err != nil {
					login.response.Code = 100007
					login.response.Data = err
					login.response.Message = "生成json失败"
					login.response.WriteJson(w)
					return
				}

				ok, err := redis.GetConnect().Do("SET", proToken, string(userJson), "ex", expire)
				defer redis.GetConnect().Close()
				if err != nil || ok != "OK" {
					login.response.Code = 100004
					login.response.Data = err
					login.response.Message = "设置登陆状态失败"
					login.response.WriteJson(w)
					return
				}

				userInfo.Password = ""
				login.response.Data = userInfo
				login.response.Message = "登陆成功"
				login.response.WriteJson(w)
				return
			} else {
				login.response.Code = 100005
				login.response.Data = err
				login.response.Message = "密码不正确"
				login.response.WriteJson(w)
				return
			}
		} else {
			login.response.Code = 100006
			login.response.Data = err
			login.response.Message = "用户不存在"
			login.response.WriteJson(w)
			return
		}
	}

	login.response.Code = 100005
	login.response.Data = err
	login.response.Message = "参数错误"
	login.response.WriteJson(w)
}

func (login *Login) CheckLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		login.response.Code = 100000
		login.response.Data = err
		login.response.Message = "读取body出错了"
		login.response.WriteJson(w)
		return
	}

	var check CheckLogin
	json.Unmarshal(req, &check)
	data, err := redis.GetConnect().Do("GET", check.Token)
	if err == nil && data != nil {
		var user model.User
		json.Unmarshal(data.([]byte), &user)
		if strings.EqualFold(user.Token, check.Token) {
			login.response.Code = 0
			login.response.Data = user
			login.response.Message = "用户正常登陆"
			login.response.WriteJson(w)
			return
		}
	}

	login.response.Code = 100005
	login.response.Data = ""
	login.response.Message = "用户异常登陆"
	login.response.WriteJson(w)
}
