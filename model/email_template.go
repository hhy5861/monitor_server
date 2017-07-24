package model

import (
	"time"
)

type (
	EmailTemplate struct {
		Id      uint   `json:"id"`
		Name    string `json:"name"`
		Subject string `json:"subject"`
		Content string `json:"content"`
		To      string `json:"to"`
		Cc      string `json:"cc"`
		Ctime   string `json:"ctime"`
		Utime   string `json:"utime"`
		Valid   int    `json:"valid"`
	}
)

func (e *EmailTemplate) TableName() string {
	return "mk_email_template"
}

func (e *EmailTemplate) Add() (*EmailTemplate, error) {
	var err error
	e.Ctime = time.Now().Format("2006-01-02 15:04:05")
	e.Utime = time.Now().Format("2006-01-02 15:04:05")
	e.Valid = 0

	err = db.Save(e).Error
	return e, err
}

func (e *EmailTemplate) FindByName(name string) (*EmailTemplate, error) {
	err := db.Where("name = ?", name).First(e).Error

	return e, err
}

func (e *EmailTemplate) GetList(param map[string]interface{}) ([]EmailTemplate, error) {
	var list []EmailTemplate
	err := db.Where("valid = ?", "0").Offset(param["page"]).Limit(param["size"]).Order("id DESC").Find(&list).Error

	return list, err
}

func (e *EmailTemplate) Delete(id uint) bool {
	e.Id = id
	if err := db.First(e).Error; err == nil {
		e.Valid = 1
		if err := db.Save(e).Error; err == nil {
			return true
		}
	}

	return false
}

func (e *EmailTemplate) FindById(id uint) (*EmailTemplate, error) {
	e.Id = id
	err := db.First(e).Error
	if err == nil {
		e.To = "huanghaiying@vipgangqin.com"
		e.Subject = "{{.Module}} 模块报警"
	}
	return e, err
}
