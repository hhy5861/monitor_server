package model

import (
	"errors"
	"gitlab.pnlyy.com/monitor_server/utils"
)

type (
	User struct {
		Id       uint   `json:"id"`
		Name     string `json:"name"`
		Mobile   string `json:"mobile"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Token    string `json:"token"`
		Ctime    string `json:"ctime"`
		Utime    string `json:"utime"`
		Valid    int    `json:"valid"`
	}
)

func (u *User) TableName() string {
	return "mk_user"
}

func (u *User) Save() (*User, error) {
	var err error
	if u.Mobile != "" && u.Email != "" {
		err = db.Where("valid = ? AND (name =? OR mobile = ? OR email = ?)", 0, u.Name, u.Mobile, u.Email).First(&u).Error
	} else if u.Mobile != "" {
		err = db.Where("valid = ? AND (name =? OR mobile = ?)", 0, u.Name, u.Mobile).First(&u).Error
	} else if u.Email != "" {
		err = db.Where("valid = ? AND (name =? OR email = ?)", 0, u.Name, u.Email).First(&u).Error
	} else {
		err = db.Where("valid = ? AND name =?", 0, u.Name).First(&u).Error
	}

	if err != nil {
		u.Ctime = utils.TimesTamp()
		u.Utime = utils.TimesTamp()
		u.Password = utils.GetMD5Hash(u.Password)
		if err := db.Save(u).Error; err == nil {
			return u, nil
		}
	}

	return nil, errors.New("save error")
}

func (u *User) GetUserByName(name string) (*User, error) {
	err := db.Where("valid = ? AND name = ?", "0", name).First(&u).Error

	return u, err
}

func (u *User) GetList(param map[string]interface{}) ([]User, error) {
	var list []User
	err := db.Where("valid = ?", "0").Offset(param["page"]).Limit(param["size"]).Order("id DESC").Find(&list).Error

	return list, err
}

func (u *User) Delete(id uint) bool {
	if err := db.Where("valid = ? AND id = ?", "0", id).First(u).Error; err == nil {
		u.Valid = 1
		u.Utime = utils.TimesTamp()
		if err := db.Save(u).Error; err == nil {
			return true
		}
	}

	return false
}
