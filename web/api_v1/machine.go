package api_v1

import (
	"time"

	"github.com/kataras/iris"
	"golang.org/x/crypto/ssh"
	"gopub3.0/model"
)

func MachineAdd(ctx iris.Context) {
	machine := model.Machine{}
	err := ctx.ReadForm(&machine)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), []byte("")))
	}
	// model.DB.NewRecord(machine)
	model.DB.Create(&machine)
	ctx.Write(model.NewResult(1, 0, "成功", machine))
}

func MachineList(ctx iris.Context) {
	machine := []model.Machine{}
	model.DB.Find(&machine)
	ctx.Write(model.NewResult(1, 0, "成功", machine))
}

func MatchineTest(ctx iris.Context) {
	id, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	machine := model.Machine{}
	machine.ID = id
	model.DB.First(&machine)

	b := []byte(machine.Rsa)
	pKey, err := ssh.ParsePrivateKey(b)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	config := ssh.ClientConfig{
		User: machine.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(pKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	c, err := ssh.Dial("tcp", machine.Ip+":"+machine.Port, &config)
	if err != nil {
		ctx.Write(model.NewResult(0, 0, err.Error(), ""))
		return
	}
	info := c.ClientVersion()
	ctx.Write(model.NewResult(1, 0, "成功", string(info)))

}
