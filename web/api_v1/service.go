package api_v1

import (
	"github.com/kataras/iris"
	"gopub3.0/model"
	"gopub3.0/mssh"
)

type serviceList struct {
	Name     string
	LocalIp  string
	Port     string
	RemoteIp string
	Status   bool
	Auto     int
	Machine  string
}

func ServiceList(ctx iris.Context) {
	machines := []model.Machine{}
	services := []model.Service{}
	model.DB.Find(&machines)
	model.DB.Find(&services)
	var servicesList []serviceList

	for _, machine := range machines {
		for _, service := range services {
			servicesList = append(servicesList, serviceList{Name: service.Name, Auto: service.Auto, LocalIp: service.Ip, Port: service.Port, RemoteIp: machine.Ip, Status: mssh.GetStatus(machine.Name), Machine: machine.Name})
		}
	}
	ctx.Write(model.NewResult(1, 0, "成功", servicesList))

}

func ServiceAdd(ctx iris.Context) {
	service := model.Service{}
	err := ctx.ReadForm(&service)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	model.DB.Create(&service)
	ctx.Write(model.NewResult(1, 0, "成功", ""))

}

func ServiceRemove(ctx iris.Context) {
	service := model.Service{}
	err := ctx.ReadForm(&service)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	model.DB.Delete(service)
	ctx.Write(model.NewResult(1, 0, "成功", ""))

}

func ServiceOn(ctx iris.Context) {
	name := ctx.PostValue("name")
	service := model.Service{}
	model.DB.Where("name = ?", name).First(&service)
	if service.ID <= 0 {
		ctx.Write(model.NewResult(0, 0, "未找到任务", ""))
		return
	}
	service.Auto = 1
	model.DB.Save(&service)
	ctx.Write(model.NewResult(1, 0, "成功开启", ""))
}

func ServiceOff(ctx iris.Context) {
	name := ctx.PostValue("name")
	service := model.Service{}
	model.DB.Where("name = ?", name).First(&service)
	if service.ID <= 0 {
		ctx.Write(model.NewResult(0, 0, "未找到任务", ""))
		return
	}
	service.Auto = 0
	model.DB.Save(&service)
	ctx.Write(model.NewResult(1, 0, "成功关闭", ""))
}
