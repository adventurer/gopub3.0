package api_v1

import (
	"github.com/kataras/iris"
	"gopub3.0/model"
)

func Login(ctx iris.Context) {
	user := model.User{}
	err := ctx.ReadForm(&user)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), []byte("")))
		return
	}
	findUser, ok := model.ValidateUser(&user)
	if !ok {
		ctx.Write(model.NewResult(0, 0, "账号或密码错误", []byte("")))
		return
	}
	if findUser.Status == 0 {
		ctx.Write(model.NewResult(0, 0, "用户被禁用", []byte("")))
		return
	}

	result := model.NewResult(1, 0, "success", findUser)
	ctx.Write(result)
}

func UserList(ctx iris.Context) {
	user := []model.User{}
	model.DB.Find(&user)
	ctx.Write(model.NewResult(1, 0, "成功", user))
}

func UserAdd(ctx iris.Context) {
	user := model.User{}
	err := ctx.ReadForm(&user)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), []byte("")))
		return
	}

	user.Status = 1
	user.Password = "w123123"
	model.DB.Create(&user)
	if user.ID <= 0 {
		ctx.Write(model.NewResult(0, 0, "创建用户失败", []byte("")))
		return
	}
	ctx.Write(model.NewResult(1, 0, "创建用户成功", ""))
}

func UserEdit(ctx iris.Context) {
	user := model.User{}
	err := ctx.ReadForm(&user)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), []byte("")))
		return
	}
	model.DB.Where("email = ?", user.Email).First(&user)
	if user.ID >= 0 {
		user.Role = user.Role
		user.Name = user.Name
		model.DB.Save(&user)
	}
	ctx.Write(model.NewResult(1, 0, "编辑用户成功", ""))
}

func UserRemove(ctx iris.Context) {
	id, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), []byte("")))
		return
	}
	user := model.User{ID: id}
	model.DB.First(&user)
	if user.Name == "" {
		ctx.Write(model.NewResult(0, 0, "未找到用户", []byte("")))
		return
	}
	model.DB.Delete(&user)
	ctx.Write(model.NewResult(1, 0, "成功", []byte("")))
}

func UserResetPass(ctx iris.Context) {
	pass := ctx.PostValue("pass")
	passwordHash := ctx.GetHeader("token")
	userID := model.ValidatePasswordHash(passwordHash)
	if userID == 0 {
		ctx.Write(model.NewResult(0, 400, "登录超时，请刷新后请重新登录", []byte("")))
		return
	}
	user := model.User{ID: userID}
	model.DB.First(&user)
	if user.Email == "" {
		ctx.Write(model.NewResult(0, 0, "未找到此用户", []byte("")))
	}
	user.Password = pass
	model.DB.Save(&user)
	ctx.Write(model.NewResult(1, 0, "修改成功", []byte("")))
}
