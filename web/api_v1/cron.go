package api_v1

import (
	"github.com/kataras/iris"
	"gopub3.0/cron"
	"gopub3.0/model"
)

func CronStart(ctx iris.Context) {
	cron.Start()
	ctx.Write(model.NewResult(1, 0, "已执行启动", ""))
}

func CronStop(ctx iris.Context) {
	cron.Stop()
	ctx.Write(model.NewResult(1, 0, "已执行停止", ""))
}

func CronRestart(ctx iris.Context) {
	cron.Restart()
	ctx.Write(model.NewResult(1, 0, "已执行重启", ""))
}

func CronOn(ctx iris.Context) {
	id, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	cron := model.Cron{ID: id}
	model.DB.First(&cron)
	if cron.Name == "" {
		ctx.Write(model.NewResult(0, 0, "未找到任务", ""))
		return
	}
	cron.Status = 1
	model.DB.Save(&cron)
	ctx.Write(model.NewResult(1, 0, "成功开启", ""))
}

func CronOff(ctx iris.Context) {
	id, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	cron := model.Cron{ID: id}
	model.DB.First(&cron)
	if cron.Name == "" {
		ctx.Write(model.NewResult(0, 0, "未找到任务", ""))
		return
	}
	cron.Status = 0
	model.DB.Save(&cron)
	ctx.Write(model.NewResult(1, 0, "成功关闭", ""))
}
