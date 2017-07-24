package model

import (
	"time"
)

type (
	Project struct {
		Id    uint   `json:"id"`
		Name  string `json:"name"`
		Token string `json:"token"`
		Ctime string `json:"ctime"`
		Utime string `json:"utime"`
		Valid int    `json:"valid"`
	}
)

func (p *Project) TableName() string {
	p.Ctime = time.Now().Format("2006-01-02 15:04:05")
	p.Utime = time.Now().Format("2006-01-02 15:04:05")
	p.Valid = 0
	return "mk_project"
}

func (p *Project) SaveProject() (*Project, error) {
	if err := db.Save(p).Error; err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Project) FindByName(name string) (*Project, error) {
	err := db.Where("name = ? AND valid = ?", name, "0").First(p).Error
	return p, err
}

func (p *Project) GetList(param map[string]interface{}) ([]Project, error) {
	var list []Project
	err := db.Where("valid = ?", "0").Offset(param["page"]).Limit(param["size"]).Order("id DESC").Find(&list).Error

	return list, err
}

func (p *Project) Delete(id uint) bool {
	p.Id = id
	if err := db.Where("valid = ?", "0").First(p).Error; err == nil && p.Id > 0 {
		p.Valid = 1
		p.Utime = time.Now().Format("2006-01-02 15:04:05")
		if err := db.Save(p).Error; err == nil {
			return true
		}
	}

	return false
}
