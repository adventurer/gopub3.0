package model

import (
	"encoding/json"
	"time"
)

type User struct {
	ID           int    `gorm:"AUTO_INCREMENT"`
	Email        string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	Password     string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	PasswordHash string `gorm:"size:255"`
	Status       int
	Role         int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Result struct {
	Sta  int
	Code int
	Msg  string
	Data interface{}
}

func NewResult(sta int, code int, msg string, data interface{}) []byte {
	result := Result{}
	result.Sta = sta
	result.Code = code
	result.Msg = msg
	result.Data = data
	rs, err := json.Marshal(result)
	if err != nil {
		panic(err.Error())
	}
	return rs
}
