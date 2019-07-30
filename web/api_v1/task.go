package api_v1

import (
	"log"
	"strings"
	"time"

	"gopub3.0/mlog"

	"github.com/kataras/iris"
	"gopub3.0/cmd"
	"gopub3.0/model"
	"gopub3.0/mssh"
	"gopub3.0/service"
)

type Commit struct {
	Key   string
	Value string
}

type InnerTask struct {
	ID      int
	Name    string
	Status  int
	Machine []InnerMachine
}

type InnerMachine struct {
	Name       string
	Step       int
	DeployStep []model.DeployStep
}

var innerTaskMap = make(map[int]InnerTask, 10)

func TaskAudit(ctx iris.Context) {
	TaskID, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), []byte("")))
		return
	}
	task := model.Task{ID: TaskID}
	model.DB.First(&task)
	if task.Name == "" {
		ctx.Write(model.NewResult(0, 0, "未找到此上线单", []byte("")))
		return
	}
	model.DB.Model(&task).Update(model.Task{Audit: 1, AuditAt: time.Now()})
	ctx.Write(model.NewResult(1, 0, "审核成功", task))
}

func DeployMessage(ctx iris.Context) {
	TaskID, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), []byte("")))
		return
	}
	result, ok := innerTaskMap[TaskID]
	if ok {
		ctx.Write(model.NewResult(1, 1, "", result))
		return
	}
	ctx.Write(model.NewResult(1, 0, "", []byte("")))

}

func TaskInfo(ctx iris.Context) {
	TaskID, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), []byte("")))
		return
	}
	task := model.Task{ID: TaskID}
	model.DB.First(&task)
	if task.Name == "" {
		ctx.Write(model.NewResult(0, 0, "未找到此上线单", []byte("")))
		return
	}

	project := model.Project{ID: task.ProjectID}
	model.DB.First(&project)
	if project.Name == "" {
		ctx.Write(model.NewResult(0, 0, "未找到此上线单关联项目", []byte("")))
		return
	}

	m := model.Machine{}
	machines := m.FindHost(project)
	if len(machines) <= 0 {
		ctx.Write(model.NewResult(0, 0, "此项目未关联主机", []byte("")))
		return
	}

	steps := []model.DeployStep{}
	model.DB.Where("project_id = ?", project.ID).Find(&steps).Order("id asc")
	if len(steps) <= 0 {
		ctx.Write(model.NewResult(0, 0, "没发现此项目的部署步骤", []byte("")))
		return
	}

	innerMachines := []InnerMachine{}
	for _, machine := range machines {
		// 注意深拷贝
		_step := make([]model.DeployStep, len(steps), len(steps))
		copy(_step, steps)
		innerMachine := InnerMachine{Name: machine.Name, Step: -1, DeployStep: _step}
		innerMachines = append(innerMachines, innerMachine)

	}
	innerTask := InnerTask{ID: TaskID, Name: task.Name, Status: task.Status, Machine: innerMachines}
	// 初始化缓存
	_, ok := innerTaskMap[TaskID]
	if !ok {
		innerTaskMap[TaskID] = innerTask

	}

	ctx.Write(model.NewResult(1, 0, "成功", innerTask))
}

func TaskList(ctx iris.Context) {
	task := []model.Task{}
	model.DB.Order("id desc").Find(&task).Limit(15)
	ctx.Write(model.NewResult(1, 0, "成功", task))
}

func TaskAdd(ctx iris.Context) {
	task := model.Task{}
	err := ctx.ReadForm(&task)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), []byte("")))
		return
	}
	project := model.Project{ID: task.ProjectID}
	model.DB.First(&project)
	if project.ID < 1 {
		ctx.Write(model.NewResult(0, 0, "该项目不存在", []byte("")))
		return
	}
	task.ProjectName = project.Name
	task.TicketAt = time.Now()
	task.Audit = project.Audit
	model.DB.Create(&task)
	if task.ID < 1 {
		ctx.Write(model.NewResult(0, 0, "写入失败", []byte("")))
		return
	}
	ctx.Write(model.NewResult(1, 0, "成功", task))

}

func TaskRemove(ctx iris.Context) {
	TaskID, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), []byte("")))
		return
	}
	task := model.Task{}
	model.DB.Where("id = ?", TaskID).First(&task)
	if task.UserName == "" {
		ctx.Write(model.NewResult(0, 0, "未找到上线单", []byte("")))
		return
	}
	if task.Status > 0 {
		ctx.Write(model.NewResult(0, 0, "不能删除已上线的上线单", []byte("")))
		return
	}
	model.DB.Delete(&task)
	ctx.Write(model.NewResult(1, 0, "成功", task))
}

func GetVersions(ctx iris.Context) {

	projectID, err := ctx.PostValueInt("project")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), []byte("")))
		return
	}
	project := model.Project{}
	model.DB.Where("id = ?", projectID).First(&project)
	repoNameIndex := strings.LastIndex(project.Repo, "/")
	repoName := project.Repo[repoNameIndex:]
	versions, err := service.GetVersions("repository" + repoName)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), []byte("")))
		return
	}

	commitList := []Commit{}
	for _, version := range versions {
		if len(version) > 0 {
			log.Println(version)
			commit := Commit{Key: version[0:7], Value: version}
			commitList = append(commitList, commit)
		}

	}
	log.Printf("%#v", commitList)
	ctx.Write(model.NewResult(1, 0, "成功", commitList))

}

func GetVersionInfo(ctx iris.Context) {
	projectID, err := ctx.PostValueInt("project")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), []byte("")))
		return
	}
	hash := ctx.PostValue("hash")
	project := model.Project{}
	model.DB.Where("id = ?", projectID).First(&project)
	repoNameIndex := strings.LastIndex(project.Repo, "/")
	repoName := project.Repo[repoNameIndex:]
	info, err := service.GetVersionInfo("repository"+repoName, hash)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), []byte("")))
		return
	}
	ctx.Write(model.NewResult(1, 0, "成功", info))

}

func TaskDeploy(ctx iris.Context) {
	TaskID, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), []byte("")))
		return
	}
	task := model.Task{ID: TaskID}
	model.DB.First(&task)
	if task.Name == "" {
		ctx.Write(model.NewResult(0, 0, "未找到此上线单", []byte("")))
		return
	}
	model.DB.Model(&task).Update(model.Task{DeployAt: time.Now()})

	project := model.Project{ID: task.ProjectID}
	model.DB.First(&project)
	if project.Name == "" {
		ctx.Write(model.NewResult(0, 0, "未找到此上线单关联项目", []byte("")))
		return
	}

	m := model.Machine{}
	machines := m.FindHost(project)
	if len(machines) <= 0 {
		ctx.Write(model.NewResult(0, 0, "此项目未关联主机", []byte("")))
		return
	}

	deployStep := []model.DeployStep{}
	model.DB.Where("project_id = ?", project.ID).Find(&deployStep)

	for _, machine := range machines {
		for k, step := range deployStep {
			conn, err := mssh.Connect(machine)
			if err != nil {
				ctx.Write(model.NewResult(0, 0, err.Error(), []byte("")))
				return
			}

			// 开始变量替换
			action := strings.Replace(step.Action, "__version__", task.Version, -1)

			mlog.Flog("publish.log", "[publish task run]", action)
			output, err := cmd.RunRemote(conn, action)
			mlog.Flog("publish.log", "[publish task result]", output)

			if err != nil {
				ctx.Write(model.NewResult(0, 0, output, []byte("")))
				return
			}
			info(TaskID, machine.Name, k, output)
			model.DB.Model(&task).Update(model.Task{Step: k})
		}
	}

	model.DB.Model(&task).Update(model.Task{DoneAt: time.Now(), Status: 1})

	ctx.Write(model.NewResult(1, 0, "成功", []byte("")))

}

func info(taskID int, name string, current int, result string) {
	newInnerTask := innerTaskMap

	for k, _ := range newInnerTask[taskID].Machine {
		if newInnerTask[taskID].Machine[k].Name == name {
			mlog.Mlog.Println("修改", name)
			mlog.Mlog.Println("action结果", newInnerTask[taskID].Machine[k].DeployStep[current].Action)

			newInnerTask[taskID].Machine[k].Step = current

			newInnerTask[taskID].Machine[k].DeployStep[current].Action += "#执行结果#" + result

		}
	}

}
