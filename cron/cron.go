package cron

import (
	"gopub3.0/cmd"
	"gopub3.0/mlog"
	"gopub3.0/mssh"

	cron "gopkg.in/robfig/cron.v2"
	"gopub3.0/model"
)

var Cron = cron.New()

func init() {

	items := []model.Cron{}
	model.DB.Find(&items)
	for _, item := range items {
		machine := model.Machine{}
		model.DB.Where("name = ?", item.Machine).First(&machine)
		command := item.Cmd
		logName := item.Name + ".log"

		Cron.AddFunc(item.Schedule, func() {
			conn, err := mssh.Connect(machine)
			if err != nil {
				mlog.Mlog.Println(err)
				return
			}
			mlog.Flog(logName, "[remote task run]", command)
			output, err := cmd.RunRemote(conn, command)
			if err != nil {
				mlog.Mlog.Println(err)
				return
			}
			mlog.Flog(logName, "[remote task result]", "run remote result:"+output)
			conn.Close()
		})

	}

	Cron.Start()
}
