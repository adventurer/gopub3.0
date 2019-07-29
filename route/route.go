package route

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"gopub3.0/model"
	"gopub3.0/web/api_v1"
)

func Init(app *iris.Application) {
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, //允许通过的主机名称
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "HEAD", "OPTIONS", "PUT", "DELETE"},
		AllowCredentials: true,
		Debug:            false,
	})
	v1 := app.Party("/api/v1/", crs).AllowMethods(iris.MethodOptions)
	{
		v1.Post("user/login", api_v1.Login)

		v1.Use(middwareAuth)
		v1.Get("user/list", api_v1.UserList)
		v1.Post("user/add", api_v1.UserAdd)
		v1.Post("user/remove", api_v1.UserRemove)
		v1.Post("user/repass", api_v1.UserResetPass)

		v1.Get("welcome", api_v1.Welcome)
		v1.Post("machine/add", api_v1.MachineAdd)
		v1.Get("machine/list", api_v1.MachineList)
		v1.Post("machine/test", api_v1.MatchineTest)
		v1.Post("proxy/off", api_v1.ProxyOff)
		v1.Post("proxy/on", api_v1.ProxyOn)

		v1.Post("service/add", api_v1.ServiceAdd)
		v1.Post("service/remove", api_v1.ServiceRemove)
		v1.Get("service/list", api_v1.ServiceList)

		v1.Post("project/add", api_v1.ProjectAdd)
		v1.Post("project/remove", api_v1.ProjectRemove)
		v1.Get("project/list", api_v1.ProjectList)
		v1.Post("project/hostadd", api_v1.HostAdd)
		v1.Post("project/init", api_v1.ProjectInit)
		v1.Post("project/chaudit", api_v1.ProjectChangeAudit)

		v1.Post("task/getversion", api_v1.GetVersions)
		v1.Post("task/getversioninfo", api_v1.GetVersionInfo)

		v1.Post("task/info", api_v1.TaskInfo)
		v1.Get("task/list", api_v1.TaskList)
		v1.Post("task/add", api_v1.TaskAdd)
		v1.Post("task/remove", api_v1.TaskRemove)
		v1.Post("task/audit", api_v1.TaskAudit)
		v1.Post("task/deploy", api_v1.TaskDeploy)
		v1.Post("task/deploymessage", api_v1.DeployMessage)

		v1.Post("cron/add", api_v1.ScheduleAdd)
		v1.Get("cron/list", api_v1.ScheduleList)

	}

}

// func midwareAuth(ctx iris.Context) {
// 	ctx.WriteString("auth need")
// 	ctx.Next()
// }

func midwareCrs(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Request-Headers", "*")

	// Access-Control-Request-Headers: x-requested-with

	ctx.Next()
}

func middwareAuth(ctx iris.Context) {
	passwordHash := ctx.GetHeader("token")
	userID := model.ValidatePasswordHash(passwordHash)
	if userID == 0 {
		ctx.Write(model.NewResult(0, 400, "登录超时，请刷新后请重新登录", []byte("")))
		return
	}
	ctx.Values().Set("user_id", userID)
	ctx.Next()
}
