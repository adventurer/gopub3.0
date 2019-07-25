package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/gopub_v3?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}
	DB.LogMode(false)

	DB.AutoMigrate(&User{}, &Machine{}, Service{}, Project{}, Task{}, DeployStep{})
}
