package api_v1

import (
	"github.com/kataras/iris"
	"gopub3.0/model"
	"gopub3.0/service"
)

func DomainList(ctx iris.Context) {
	id, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), id))
		return
	}
	domain := model.DomainAccess{ID: id}
	model.DB.Find(&domain)
	list, err := service.InfoHostList(domain, 1)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), domain))
		return
	}
	ctx.Write(model.NewResult(1, 0, "获取成功", list.DomainRecords))

}

func DomainNew(ctx iris.Context) {
	domain := ctx.PostValue("Domain")
	subdomain := ctx.PostValue("SubDomain")
	ip := ctx.PostValue("Ip")
	if domain == "" || subdomain == "" || ip == "" {
		ctx.Write(model.NewResult(0, 0, "某条件为空", ""))
		return
	}
	ok, err := service.NewHost(domain, subdomain, ip)
	if !ok {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	ctx.Write(model.NewResult(1, 0, "新增成功", ""))

}

func DomainRemove(ctx iris.Context) {
	subdomianId := ctx.PostValue("subdomianId")
	if subdomianId == "" {
		ctx.Write(model.NewResult(0, 0, "域名id不能为空", ""))
		return
	}
	id, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), id))
		return
	}
	domain := model.DomainAccess{ID: id}
	model.DB.Find(&domain)

	ok, err := service.DelHost(domain, subdomianId)
	if !ok {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	ctx.Write(model.NewResult(1, 0, "删除成功", ""))
}

func DomainSettingAdd(ctx iris.Context) {
	domain := model.DomainAccess{}
	err := ctx.ReadForm(&domain)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), domain))
		return
	}
	model.DB.Create(&domain)
	ctx.Write(model.NewResult(1, 0, "新增成功", domain))
}

func DomainSettingRemove(ctx iris.Context) {
	id, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), id))
		return
	}
	domain := model.DomainAccess{ID: id}
	model.DB.Find(&domain)
	if domain.Domain != "" {
		model.DB.Delete(&domain)
	}
	ctx.Write(model.NewResult(1, 0, "删除成功", domain))
}

func DomainSettingList(ctx iris.Context) {
	list := []model.DomainAccess{}
	model.DB.Find(&list)
	ctx.Write(model.NewResult(1, 0, "获取成功", list))
}
