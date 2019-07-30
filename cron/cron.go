package cron

import (
	"fmt"
	"time"

	"gopub3.0/cmd"
	"gopub3.0/mlog"
	"gopub3.0/mssh"

	cron "gopkg.in/robfig/cron.v2"
	"gopub3.0/model"
)

var Cron = cron.New()

var startQue = make(chan bool, 1)
var stopQue = make(chan bool, 1)
var restartQue = make(chan bool, 1)

func init() {
	go func() {
		for {
			select {
			case <-startQue:
				items := []model.Cron{}
				model.DB.Find(&items)
				for _, item := range items {
					machine := model.Machine{}
					model.DB.Where("name = ?", item.Machine).First(&machine)
					command := item.Cmd
					logName := item.Name
					if item.Status == 1 {
						Cron.AddFunc(item.Schedule, func() {
							defer recoverName()
							session, err := mssh.Connect(machine)
							if err != nil {
								mlog.Mlog.Println(err)
							} else {
								mlog.Flog(logName, "[remote task run]", command)
								output, err := cmd.RunRemote(session, command)
								if err != nil {
									mlog.Mlog.Println(err)
									return
								}
								mlog.Flog(logName, "[remote task result]", output)
							}

						})
					}

				}
				Cron.Start()
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-stopQue:
				Cron.Stop()
			default:
				time.Sleep(1 * time.Second)

			}
		}
	}()
}

func Start() {
	mlog.Mlog.Println("start cron")
	startQue <- true
}

func Stop() {
	mlog.Mlog.Println("stop cron")
	entries := Cron.Entries()
	for _, entry := range entries {
		Cron.Remove(entry.ID)
	}
	stopQue <- true
}

func Restart() {
	Stop()
	time.Sleep(3 * time.Second)
	Start()
}

func recoverName() {
	if r := recover(); r != nil {
		fmt.Println("计划任务崩溃，延迟1秒重启", r)
		time.Sleep(1 * time.Second)
	}
}
