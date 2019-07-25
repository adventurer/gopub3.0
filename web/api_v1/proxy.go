package api_v1

import (
	"github.com/kataras/iris"
	"gopub3.0/model"
	"gopub3.0/mssh"
)

func ProxyOn(ctx iris.Context) {
	name := ctx.PostValueTrim("name")
	serviceName := ctx.PostValueTrim("service")

	maching := model.Machine{}
	model.DB.Where("name = ?", name).First(&maching)

	service := model.Service{}
	model.DB.Where("name = ?", serviceName).First(&service)

	proxyMachine := mssh.ProxyConn{Name: maching.Name, User: maching.User, Host: maching.Ip, Port: maching.Port, Rsa: maching.Rsa, LocalIp: service.Ip, LocalPort: service.Port}
	mssh.NewProxy(proxyMachine)
	ctx.WriteString(name)
}

func ProxyOff(ctx iris.Context) {
	name := ctx.PostValueTrim("name")
	mssh.StopProxy(name)
	ctx.WriteString(name)
}
