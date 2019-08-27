package api_v1

import (
	"github.com/kataras/iris"
	"gopub3.0/model"
	"gopub3.0/nat"
	"gopub3.0/service"
)

func DockerContainerDeploy(ctx iris.Context) {
	container := model.ContainerDeploy{}
	err := ctx.ReadForm(&container)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, "参数错误", ""))
		return
	}
	machine := model.Machine{ID: container.Machine}
	model.DB.First(&machine)
	if machine.Name == "" {
		ctx.Write(model.NewResult(0, 0, "没发现机器", ""))
		return
	}
	output, err := service.ContainerNew(machine, container)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, output+err.Error(), ""))
		return
	}
	ctx.Write(model.NewResult(1, 0, "创建成功", container))

}

func DockerMachines(ctx iris.Context) {
	machine := []model.Machine{}
	model.DB.Where("type = ?", 2).Find(&machine)
	ctx.Write(model.NewResult(1, 0, "获取成功", machine))
}

func DockerNetworkList(ctx iris.Context) {
	id, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	machine := model.Machine{}
	model.DB.Where("id = ?", id).First(&machine)
	containers, err := service.NetworkList(machine)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), containers))
		return
	}
	ctx.Write(model.NewResult(1, 0, "获取成功", containers))
}

func DockerContainerList(ctx iris.Context) {
	id, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	machine := model.Machine{}
	model.DB.Where("id = ?", id).First(&machine)
	containers, err := service.ContainerList(machine)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), containers))
		return
	}
	ctx.Write(model.NewResult(1, 0, "获取成功", containers))
}

func DockerStartContainer(ctx iris.Context) {
	id := ctx.PostValue("id")
	machine := model.Machine{}
	model.DB.Where("type = ?", 2).First(&machine)
	output, err := service.StartContainer(machine, id)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	ctx.Write(model.NewResult(1, 0, output, ""))
}

func DockerStopContainer(ctx iris.Context) {
	id := ctx.PostValue("id")
	machine := model.Machine{}
	model.DB.Where("type = ?", 2).First(&machine)
	output, err := service.StopContainer(machine, id)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	ctx.Write(model.NewResult(1, 0, output, ""))
}

func DockerNatTable(ctx iris.Context) {
	machine := model.Machine{}
	model.DB.Where("type = ?", 2).First(&machine)
	service.NatTable(machine)
}

func DockerAddPort(ctx iris.Context) {
	form := model.DockerPort{}
	err := ctx.ReadForm(&form)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
	}
	model.DB.Create(&form)
	nat.DockerPort = append(nat.DockerPort, form)
	ctx.Write(model.NewResult(1, 0, "新增成功", ""))
}

func DockerRemovePort(ctx iris.Context) {
	id, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
	}
	form := model.DockerPort{ID: id}
	model.DB.Delete(&form)
	ctx.Write(model.NewResult(1, 0, "删除成功", ""))
}

func DockerPortList(ctx iris.Context) {
	portList := []model.DockerPort{}
	model.DB.Find(&portList)
	ctx.Write(model.NewResult(1, 0, "成功", portList))
}
