package main

import (
	"net/http"

	"github.com/kataras/iris"

	"gopub3.0/cron"
	"gopub3.0/mlog"
	"gopub3.0/mssh"
	_ "gopub3.0/nat"
	"gopub3.0/route"
)

func main() {

	cron.Start()
	mssh.Begin()

	app := iris.New()
	app.Use(mlog.CustomLogger)

	route.Init(app)

	app.Build()
	srv := &http.Server{Handler: app, Addr: ":8088"}
	println("Start a server listening on http://localhost:8088")

	srv.ListenAndServe()

}
