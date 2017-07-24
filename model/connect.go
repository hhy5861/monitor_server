package model

import (
	"github.com/hhy5861/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"fmt"
)

type (
	Mysql struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
		Charset  string `yaml:"charset"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Dialect  string `yaml:"dialect"`
		Debug    bool   `yaml:"debug"`
	}
)

var (
	db *gorm.DB
)

func NewMysql(host, dbname, user, password, charset, dialect string, post int, debug bool) {
	m := &Mysql{
		Host:     host,
		Dbname:   dbname,
		User:     user,
		Password: password,
		Charset:  charset,
		Port:     post,
		Dialect:  dialect,
		Debug:    debug,
	}

	m.connectDb()
}

func (m *Mysql) connectDb() {
	var err error
	addres := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s&parseTime=False&loc=Local", m.User, m.Password, m.Host, m.Port, m.Dbname, m.Charset)
	db, err = gorm.Open("mysql", addres)
	if err != nil {
		ps := logrus.Params{
			"dialect": m.Dialect,
			"connect": addres,
		}

		logrus.Fatal(ps, err, "mysql connect error")
	}

	db.LogMode(m.Debug)
	db.DB().SetMaxOpenConns(10)
}
