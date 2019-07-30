package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open("mysql", "root:112215334@tcp(192.168.1.202:3306)/gopub_v3?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}
	DB.LogMode(false)

	DB.AutoMigrate(&User{}, &Machine{}, Service{}, Project{}, Task{}, DeployStep{}, Cron{})
	initAdmin()
}

func initAdmin() {
	user := User{Email: "admin@qq.com"}
	DB.Where("email = ?", "admin@qq.com").First(&user)
	if user.ID <= 0 {
		DB.Create(&User{Email: "admin@qq.com", Password: "w123123", Name: "管理员", Status: 1, Role: 999})
	}
}
