package model

import (
	"crypto/md5"
	"fmt"
	"log"
	"time"
)

func ValidateUser(user *User) (User, bool) {
	findUser := User{}
	DB.Where("email = ?", user.Email).First(&findUser)
	if findUser.Password != user.Password {
		return findUser, false
	}
	hash := genPasswordHash(user)
	findUser.PasswordHash = hash
	DB.Model(&findUser).Update("password_hash", hash)
	return findUser, true
}

func genPasswordHash(user *User) string {
	time, _ := time.Now().MarshalText()
	hash := md5.Sum(time)
	return fmt.Sprintf("%x", hash)
}

func ValidatePasswordHash(hash string) int {
	findUser := User{}
	DB.Where("password_hash = ?", hash).First(&findUser)
	log.Println(findUser)
	if findUser.ID == 0 {
		return 0
	}
	return findUser.ID
}
