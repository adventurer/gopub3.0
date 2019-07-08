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
		v1.Get("welcome", api_v1.Welcome)

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
		ctx.HTML(`登录超时，请重新登录`)
		return
	}
	ctx.Values().Set("user_id", userID)
	ctx.Next()
}
