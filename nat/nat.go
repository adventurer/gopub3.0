package nat

import (
	"fmt"
	"strings"
	"time"

	"gopub3.0/cmd"
	"gopub3.0/mlog"
	"gopub3.0/model"
	"gopub3.0/mssh"
)

var DockerPort []model.DockerPort
var Machines []model.Machine

func init() {
	model.DB.Find(&DockerPort)
	model.DB.Find(&Machines)
	go keepNat()
}

func keepNat() {
	for {
		for _, dp := range DockerPort {
			machine, err := getMachine(dp.MachineName)
			if err != nil {
				mlog.Mlog.Println(err.Error())
				continue
			}
			if !isExist(machine, dp) {
				addRule(machine, dp)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func getMachine(machineName string) (machine model.Machine, err error) {
	for _, machine := range Machines {
		if machine.Name == machineName {
			return machine, nil
		}
	}
	return machine, fmt.Errorf("未发现主机：%s", machineName)
}

func isExist(machine model.Machine, dockerPort model.DockerPort) bool {
	conn, err := mssh.Connect(machine)
	if err != nil {
		return false
	}
	action := fmt.Sprintf(`iptables -t nat -L -n -v`)
	mlog.Flog("nat", "[nat command run]", action)
	output, err := cmd.RunRemote(conn, action)
	mlog.Flog("nat", "[nat command result]", output)
	if err != nil {
		mlog.Flog("nat", "[nat command result]", err.Error())
		return true
	}

	// search := fmt.Sprintf("tcp dpt:%s to:%s:%s", dockerPort.Port, dockerPort.ToIp, dockerPort.ToPort)
	// log.Println(output, ":", search)
	// log.Println(strings.Contains(output, search))

	if strings.Index(output, fmt.Sprintf("tcp dpt:%s to:%s:%s", dockerPort.Port, dockerPort.ToIp, dockerPort.ToPort)) < 0 {
		return false
	}
	return true
}

func addRule(machine model.Machine, dockerPort model.DockerPort) bool {
	conn, err := mssh.Connect(machine)
	if err != nil {
		return false
	}
	action := fmt.Sprintf(`iptables -t nat -A PREROUTING -p tcp --dport %s -j DNAT --to-destination %s:%s`, dockerPort.Port, dockerPort.ToIp, dockerPort.ToPort)
	mlog.Flog("nat", "[nat command run]", action)
	output, err := cmd.RunRemote(conn, action)
	mlog.Flog("nat", "[nat command result]", output)
	if err != nil {
		return false
	}
	return true
}
