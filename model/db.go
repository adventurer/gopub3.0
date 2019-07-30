package model

import (
	"fmt"

	"gopub3.0/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func init() {
	var err error
	dbURI := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/gopub_v3?charset=utf8&parseTime=True&loc=Local",
		config.Variable.Database.User,
		config.Variable.Database.Password,
		config.Variable.Database.Host,
		config.Variable.Database.Port,
	)

	DB, err = gorm.Open("mysql", dbURI)
	if err != nil {
		panic(err.Error())
	}
	DB.LogMode(true)

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
