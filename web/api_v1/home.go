package api_v1

import (
	"github.com/kataras/iris"
	"github.com/shirou/gopsutil/mem"
	"gopub3.0/model"
)

func Welcome(ctx iris.Context) {
	v, _ := mem.VirtualMemory()
	ctx.Write(model.NewResult(1, 0, "success", v))
}
