package api_v1

import (
	"mime/multipart"
	"os"
	"strings"

	"github.com/kataras/iris"
	"gopub3.0/model"
	"gopub3.0/nat"
	"gopub3.0/service"
)

type Files struct {
	Name string
}

func DockerFileRemove(ctx iris.Context) {
	file := ctx.PostValue("File")
	err := os.Remove("./uploads/" + file)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	ctx.Write(model.NewResult(1, 0, "删除成功", ""))
}

func DockerFileUp(ctx iris.Context) {
	ctx.UploadFormFiles("./uploads", beforeSave)
}

func beforeSave(ctx iris.Context, file *multipart.FileHeader) {
	ip := ctx.RemoteAddr()
	// make sure you format the ip in a way
	// that can be used for a file name (simple case):
	ip = strings.Replace(ip, ".", "_", -1)
	ip = strings.Replace(ip, ":", "_", -1)

	// you can use the time.Now, to prefix or suffix the files
	// based on the current time as well, as an exercise.
	// i.e unixTime :=	time.Now().Unix()
	// prefix the Filename with the $IP-
	// no need for more actions, internal uploader will use this
	// name to save the file into the "./uploads" folder.
	// file.Filename = ip + "-" + file.Filename
}

func DockerFiles(ctx iris.Context) {
	output, err := service.DockerFiles()
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	filesArr := strings.Split(strings.TrimSpace(output), "\n")
	if len(filesArr) <= 0 {
		ctx.Write(model.NewResult(0, 0, "目录为空", ""))
		return
	}
	files := []Files{}
	for _, file := range filesArr {
		files = append(files, Files{Name: file})
	}
	ctx.Write(model.NewResult(1, 0, "获取成功", files))

}

func DockerFileDeploy(ctx iris.Context) {
	machine, err := ctx.PostValueInt("Machine")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	file := ctx.PostValue("File")

	Machine := model.Machine{ID: machine}
	model.DB.First(&Machine)
	if Machine.Name == "" {
		ctx.Write(model.NewResult(0, 0, "没发现机器", ""))
		return
	}
	service.ContainerNewFromFile(Machine, file)
}

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

func DockerNetworkNew(ctx iris.Context) {
	name := ctx.PostValue("Name")
	if name == "" {
		ctx.Write(model.NewResult(0, 0, "网络名称不能为空", ""))
		return
	}
	subnet := ctx.PostValue("SubNet")
	if subnet == "" {
		ctx.Write(model.NewResult(0, 0, "子网不能为空", ""))
		return
	}
	machine := ctx.PostValue("Machine")
	if machine == "" {
		ctx.Write(model.NewResult(0, 0, "未选择主机", ""))
		return
	}
	Machine := model.Machine{}
	model.DB.Where("id = ?", machine).First(&Machine)

	ok, err := service.NewNetwork(Machine, name, subnet)
	if !ok {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	ctx.Write(model.NewResult(1, 0, "创建成功", ""))
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
	machineName := ctx.PostValue("machine")
	form := model.DockerPort{ID: id}
	model.DB.Find(&form)
	if form.MachineName == "" {
		ctx.Write(model.NewResult(0, 0, "未找到数据库记录", ""))
		return
	}

	machine := model.Machine{}
	model.DB.Where("name = ?", machineName).First(&machine)
	if machine.ID <= 0 {
		ctx.Write(model.NewResult(0, 0, "未找到主机", ""))
		return
	}
	isDel := nat.RemoveRule(machine, form)
	if isDel {
		model.DB.Delete(&form)
	}
	ctx.Write(model.NewResult(1, 0, "删除成功", ""))
}

func DockerPortList(ctx iris.Context) {
	portList := []model.DockerPort{}
	model.DB.Find(&portList)
	ctx.Write(model.NewResult(1, 0, "成功", portList))
}

func DockerRemove(ctx iris.Context) {
	id, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	name := ctx.PostValue("name")
	if name == "" {
		ctx.Write(model.NewResult(0, 0, "容器名必须", ""))
		return
	}
	machine := model.Machine{ID: id}
	model.DB.First(&machine)
	if machine.Name == "" {
		ctx.Write(model.NewResult(0, 0, "没发现机器", ""))
		return
	}
	output, err := service.ContainerRemove(machine, name)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	ctx.Write(model.NewResult(1, 0, "删除成功", output))
}

func DockerImages(ctx iris.Context) {
	id, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	machine := model.Machine{ID: id}
	model.DB.First(&machine)
	if machine.Name == "" {
		ctx.Write(model.NewResult(0, 0, "没发现机器", ""))
		return
	}
	images, err := service.ImagesList(machine)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	ctx.Write(model.NewResult(1, 0, "获取成功", images))

}
func DockerNetworkRemove(ctx iris.Context) {
	id, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	name := ctx.PostValue("name")
	if name == "" {
		ctx.Write(model.NewResult(0, 0, "网络名称不能为空", ""))
		return
	}
	machine := model.Machine{ID: id}
	model.DB.First(&machine)
	if machine.Name == "" {
		ctx.Write(model.NewResult(0, 0, "没发现机器", ""))
		return
	}
	output, err := service.NetworkRemove(machine, name)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	ctx.Write(model.NewResult(1, 0, "删除成功", output))
}
