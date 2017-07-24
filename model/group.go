package model

import (
	"github.com/kataras/iris/core/errors"
	"strings"
	"time"
)

type (
	Group struct {
		Id           uint   `json:"id"`
		Name         string `json:"name"`
		User_id_list string `json:"user_id_list"`
		Ctime        string `json:"ctime"`
		Utime        string `json:"utime"`
		Valid        int    `json:"valid"`
	}
)

func (g *Group) TableName() string {
	return "mk_user_group"
}

func (g *Group) Save() (*Group, error) {
	err := db.Where("valid = ? AND name =?", 0, g.Name).First(g).Error
	if err != nil {
		g.Ctime = time.Now().Format("2006-01-02 15:04:05")
		g.Utime = time.Now().Format("2006-01-02 15:04:05")
		if err := db.Save(g).Error; err == nil {
			return g, nil
		}
	}

	return nil, errors.New("名字已经存在")
}

func (g *Group) GetList(param map[string]interface{}) ([]*Group, error) {
	var list []*Group
	err := db.Where("valid = ?", "0").Offset(param["page"]).Limit(param["size"]).Order("id DESC").Find(&list).Error
	for _, v := range list {
		var nameList []string
		list := strings.Split(v.User_id_list, ",")
		for _, val := range list {
			var user User
			err := db.Table("mk_user").Where("id = ? AND valid = ?", val, "0").First(&user).Error
			if err == nil && user.Id > 0 {
				nameList = append(nameList, user.Name)
			}
		}

		v.User_id_list = strings.Join(nameList, ",")
	}

	return list, err
}

func GetGroupIdByEmail(id int) (map[string]string, error) {
	var g Group
	var emailMap = make(map[string]string)
	var err error
	err = db.Where("valid = ? AND id = ?", "0", id).First(&g).Error
	if err == nil {
		user := User{}
		var userList []User

		err = db.Table(user.TableName()).Find(&userList).Error
		if err == nil {
			for _, v := range userList {
				if v.Email != "" {
					emailMap[v.Name] = v.Email
				}
			}
		}
	}

	return emailMap, err
}

/*func (g *Group) Delete(id uint) bool {
	if err := db.Where("valid = ? AND id = ?", "0", id).First(u).Error; err == nil {
		u.Valid = 1
		u.Utime = time.Now().Format("2006-01-02 15:04:05")
		if err := db.Save(u).Error; err == nil {
			return true
		}
	}

	return false
}*/
