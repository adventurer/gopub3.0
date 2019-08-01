package api_v1

import (
	"strings"
	"time"

	"github.com/kataras/iris"
	"gopub3.0/cmd"
	"gopub3.0/model"
	"gopub3.0/mssh"
)

type project struct {
	ID         int    `gorm:"AUTO_INCREMENT"`
	Name       string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	Repo       string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	Deploy     string `gorm:"size:255"`
	DeployName string `gorm:"size:255"`
	Host       string `gorm:"size:255"`
	Audit      int
	RepoReal   string `gorm:"size:255"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeployStep []DeployStep `form:"DeployStep"`
}

type DeployStep struct {
	Title  string
	Action string
}

func ProjectStepRemove(ctx iris.Context) {
	id, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	step := model.DeployStep{ID: id}
	model.DB.First(&step)
	if step.Title == "" {
		ctx.Write(model.NewResult(0, 0, "未找到该节点", ""))
		return
	}
	model.DB.Delete(&step)
	ctx.Write(model.NewResult(1, 0, "已删除节点:"+step.Title, ""))
}

func ProjectStepEdit(ctx iris.Context) {
	projectID, err := ctx.PostValueInt("project")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}

	step := ctx.PostValueIntDefault("step", 0)

	action := ctx.PostValue("action")
	title := ctx.PostValue("title")

	stepRecord := model.DeployStep{ID: step}
	model.DB.First(&stepRecord)
	if step == 0 {
		newRecord := model.DeployStep{ProjectID: projectID, Title: title, Action: action}
		model.DB.Create(&newRecord)
		ctx.Write(model.NewResult(1, 0, "新增成功", ""))
		return
	}
	stepRecord.Action = action
	stepRecord.Title = title
	model.DB.Save(&stepRecord)
	ctx.Write(model.NewResult(1, 0, "修改成功", ""))
}

func ProjectSteps(ctx iris.Context) {
	id, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), []byte("")))
		return
	}
	steps := []model.DeployStep{}
	model.DB.Where("project_id = ?", id).Find(&steps)
	ctx.Write(model.NewResult(1, 0, "成功", steps))
}

func ProjectList(ctx iris.Context) {

	projects := []model.Project{}
	model.DB.Find(&projects)

	innerProjects := []project{}
	for _, item := range projects {
		deploySteps := []DeployStep{}
		model.DB.Where("project_id = ?", item.ID).Find(&deploySteps)

		p := project{ID: item.ID, Name: item.Name, Deploy: item.Deploy, DeployName: item.DeployName, Repo: item.Repo, RepoReal: item.RepoReal, Host: item.Host, Audit: item.Audit, DeployStep: deploySteps}
		innerProjects = append(innerProjects, p)
	}

	ctx.Write(model.NewResult(1, 0, "成功", innerProjects))
}

func ProjectAdd(ctx iris.Context) {
	innerProject := project{}
	err := ctx.ReadJSON(&innerProject)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), []byte("")))
		return
	}

	tx := model.DB.Begin()

	// new project
	project := model.Project{Name: innerProject.Name, Repo: innerProject.Repo, RepoReal: innerProject.RepoReal, Deploy: innerProject.Deploy, DeployName: innerProject.DeployName, Host: innerProject.Host, Audit: innerProject.Audit}
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), []byte("")))
	}
	tx.Create(&project)
	if project.ID <= 0 {
		ctx.Write(model.NewResult(0, 0, "创建项目失败", []byte("")))
		tx.Rollback()
		return
	}
	// new local repository
	repoNameIndex := strings.LastIndex(project.Repo, "/")
	if repoNameIndex < 0 {
		ctx.Write(model.NewResult(1, 0, "无法初始化本地仓库，非法仓库地址？", []byte("")))
		return
	}
	repoName := project.Repo[repoNameIndex:]
	cmd.RunLocal("git clone " + project.RepoReal + " repository/" + repoName)
	// new deploy step

	for _, step := range innerProject.DeployStep {
		projectStep := model.DeployStep{Title: step.Title, Action: step.Action, ProjectID: project.ID}
		tx.Create(&projectStep)
		if projectStep.ID <= 0 {
			tx.Rollback()
			ctx.Write(model.NewResult(0, 0, "创建部署步骤时失败", []byte("")))
			return
		}
	}
	tx.Commit()

	ctx.Write(model.NewResult(1, 0, "成功", project))
}

func ProjectRemove(ctx iris.Context) {
	project := model.Project{}
	err := ctx.ReadForm(&project)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	model.DB.Delete(project)
	ctx.Write(model.NewResult(1, 0, "成功", ""))
}

func HostAdd(ctx iris.Context) {
	projectId, err := ctx.PostValueInt("Project")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	host := ctx.PostValue("Host")

	project := model.Project{ID: projectId}
	model.DB.Find(&project)
	model.DB.Model(&project).Update("host", host)
	ctx.Write(model.NewResult(1, 0, "成功", ""))
}

func ProjectInit(ctx iris.Context) {
	projectID, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}

	project := model.Project{ID: projectID}
	model.DB.First(&project)
	if project.Name == "" {
		ctx.Write(model.NewResult(0, 0, "未找到此项目", []byte("")))
		return
	}

	m := model.Machine{}
	machines := m.FindHost(project)
	if len(machines) <= 0 {
		ctx.Write(model.NewResult(0, 0, "此项目未关联主机", []byte("")))
		return
	}
	for _, machine := range machines {
		conn, err := mssh.Connect(machine)
		if err != nil {
			ctx.Write(model.NewResult(0, 0, err.Error(), ""))
			return
		}
		action := "git clone " + project.Repo + " " + project.Deploy + project.DeployName
		output, err := cmd.RunRemote(conn, action)
		if err != nil {
			ctx.Write(model.NewResult(0, 0, output, ""))
			return
		}
	}
	ctx.Write(model.NewResult(1, 0, "初始化成功", []byte("")))

}

func ProjectChangeAudit(ctx iris.Context) {
	projectID, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}

	project := model.Project{ID: projectID}
	model.DB.First(&project)
	if project.Name == "" {
		ctx.Write(model.NewResult(0, 0, "未找到此项目", []byte("")))
		return
	}
	project.Audit = (project.Audit + 1) % 2
	model.DB.Save(&project)
	ctx.Write(model.NewResult(1, 0, "修改成功", ""))

}
