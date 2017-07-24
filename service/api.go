package service

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gitlab.pnlyy.com/monitor_server/api"
	"log"
	"net/http"
	"strings"
)

var (
	TOKEN = "xxxx"
)

func (svc *Service) run() {
	svc.router = httprouter.New()

	addres := fmt.Sprintf("%s:%d", svc.Host, svc.Post)
	log.Println("run server:", addres)
	http.ListenAndServe(addres, svc.routerConfig())
}

func (svc *Service) routerConfig() *httprouter.Router {
	login := api.Login{}
	svc.router.POST("/login/login", login.Login)
	svc.router.POST("/user/check-login", svc.auth(login.CheckLogin))

	project := api.Project{}
	svc.router.POST("/project/add", svc.auth(project.Add))
	svc.router.POST("/project/delete", svc.auth(project.Del))
	svc.router.GET("/project/list", svc.auth(project.List))

	config := api.MonitorConfig{}
	svc.router.POST("/config/add", svc.auth(config.Add))
	svc.router.POST("/config/delete", svc.auth(config.Del))
	svc.router.GET("/config/list", svc.auth(config.List))

	email := api.Email{}
	svc.router.POST("/email/add", svc.auth(email.Add))
	svc.router.GET("/email/list", svc.auth(email.List))

	user := api.User{}
	svc.router.POST("/user/add", svc.auth(user.Add))
	svc.router.POST("/user/delete", svc.auth(user.Delete))
	svc.router.GET("/user/list", svc.auth(user.List))

	group := api.Group{}
	svc.router.POST("/group/add", svc.auth(group.Add))
	svc.router.GET("/group/list", svc.auth(group.List))

	return svc.router
}

func (svc *Service) auth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token := r.Header.Get("token")
		hasAuth := svc.basicAuth(token)
		if hasAuth {
			h(w, r, ps)
		} else {
			w.Header().Set("Content-Type", "application/json;charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			res := map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"data":    "",
				"message": http.StatusText(http.StatusUnauthorized),
			}

			body, _ := json.Marshal(res)
			fmt.Fprint(w, string(body))
		}
	}
}

func (svc *Service) basicAuth(token string) bool {
	if strings.EqualFold(token, TOKEN) {
		return true
	}

	return false
}
