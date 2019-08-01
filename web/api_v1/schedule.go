package api_v1

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/kataras/iris"
	"gopub3.0/model"
)

func ScheduleList(ctx iris.Context) {
	cron := []model.Cron{}
	model.DB.Find(&cron)
	ctx.Write(model.NewResult(1, 0, "成功", cron))

}

func ScheduleAdd(ctx iris.Context) {
	cron := model.Cron{}
	err := ctx.ReadForm(&cron)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), []byte("")))
		return
	}

	md5 := md5.New()
	md5.Write([]byte(time.Now().String()))
	cron.Unique = hex.EncodeToString(md5.Sum(nil))
	model.DB.Create(&cron)
	if cron.ID <= 0 {
		ctx.Write(model.NewResult(0, 0, "添加失败", ""))
		return
	}
	ctx.Write(model.NewResult(1, 0, "添加成功", ""))

}

func ScheduleRemove(ctx iris.Context) {
	id, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	cron := model.Cron{ID: id}
	model.DB.First(&cron)
	if cron.Name == "" {
		ctx.Write(model.NewResult(0, 0, "未找到此任务", ""))
		return
	}
	model.DB.Delete(&cron)
	ctx.Write(model.NewResult(1, 0, "删除成功", ""))
}

func ScheduleLog(ctx iris.Context) {
}

func ScheduleEdit(ctx iris.Context) {
	cronForm := model.Cron{}
	ctx.ReadForm(&cronForm)
	if cronForm.ID <= 0 {
		ctx.Write(model.NewResult(0, 0, "无效的表单", ""))
		return
	}
	cron := model.Cron{ID: cronForm.ID}
	model.DB.First(&cron)
	if cronForm.Name == "" {
		ctx.Write(model.NewResult(0, 0, "未找到相关修改数据", ""))
		return
	}
	cron.Machine = cronForm.Machine
	cron.Cmd = cronForm.Cmd
	cron.Schedule = cronForm.Schedule
	cron.Status = cronForm.Status
	cron.Name = cronForm.Name
	model.DB.Save(&cron)
	ctx.Write(model.NewResult(1, 0, "修改成功", ""))

}
