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
}

func ScheduleLog(ctx iris.Context) {
}

func ScheduleEdit(ctx iris.Context) {
}
