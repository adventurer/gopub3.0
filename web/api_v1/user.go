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

	result := model.NewResult(1, 0, "success", findUser)
	ctx.Write(result)
}
