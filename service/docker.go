package service

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"gopub3.0/cmd"
	"gopub3.0/mlog"
	"gopub3.0/model"
	"gopub3.0/mssh"
)

type DockerContainer struct {
	ID     string
	Name   string
	IP     string
	Status string
	Ports  string
}
type proxyInfo struct {
	HostPort      string
	ContainerIp   string
	ContrinerPort string
	Stauts        chan bool
}

type DockerNetworks struct {
	ID     string
	Name   string
	IP     string
	Driver string
	Scope  string
}

type NatRules struct {
}

var dockerProxyChan = make(map[string]chan bool, 0)

func ContainerNew(machine model.Machine, containerDeploy model.ContainerDeploy) (output string, err error) {
	conn, err := mssh.Connect(machine)
	if err != nil {
		return "", errors.New("无法连接主机")
	}
	action := fmt.Sprintf("docker run -itd --network %s --ip %s --name %s %s", containerDeploy.Network, containerDeploy.Ip, containerDeploy.Name, containerDeploy.Image)
	mlog.Flog("docker", "[deploy docker command run]", action)
	output, err = cmd.RunRemote(conn, action)
	mlog.Flog("docker", "[deploy docker command result]", output)
	if err != nil {
		return output, errors.New("远程命令错误")
	}
	return
}

func NatTable(machine model.Machine) (natrules []NatRules, err error) {
	conn, err := mssh.Connect(machine)
	if err != nil {
		return natrules, errors.New("无法连接主机")
	}
	action := "iptables -t nat -L PREROUTING -nv --line"
	mlog.Flog("docker", "[iptables command run]", action)
	output, err := cmd.RunRemote(conn, action)
	mlog.Flog("docker", "[iptables command result]", output)
	if err != nil {
		return natrules, errors.New("远程命令错误")
	}
	return
}

func NetworkList(machine model.Machine) (networks []DockerNetworks, err error) {
	conn, err := mssh.Connect(machine)
	// defer conn.Close()
	if err != nil {
		return networks, errors.New("无法连接主机")
	}
	action := "docker network inspect -f='{{.Id}}|{{.Name}}|{{.Driver}}|{{.IPAM.Config}}|{{.Scope}}' $(docker network ls -q)"
	mlog.Flog("docker", "[docker command run]", action)
	output, err := cmd.RunRemote(conn, action)
	mlog.Flog("docker", "[docker command result]", output)

	items := strings.Split(strings.TrimSpace(output), "\n")

	if len(items) < 1 {
		return networks, errors.New("未发现容器")
	}
	for _, item := range items {
		arr := strings.Split(item, "|")
		network := DockerNetworks{ID: arr[0], Name: arr[1], IP: arr[3], Driver: arr[2], Scope: arr[4]}
		networks = append(networks, network)
	}
	return
}

func ContainerList(machine model.Machine) (containers []DockerContainer, err error) {
	conn, err := mssh.Connect(machine)
	// defer conn.Close()
	if err != nil {
		return containers, errors.New("无法连接主机")
	}
	hasdocker, _ := cmd.RunRemote(conn, "docker ps -aq")
	if hasdocker == "" {
		return containers, errors.New("没有发现容器")
	}
	conn, err = mssh.Connect(machine)
	if err != nil {
		return containers, errors.New("无法连接主机")
	}
	action := "docker inspect -f='{{.Id}}|{{.Name}}|{{.NetworkSettings.Networks}}|{{.State.Status}}|{{.NetworkSettings.Ports}}' $(docker ps -aq)"
	mlog.Flog("docker", "[docker command run]", action)
	output, err := cmd.RunRemote(conn, action)
	mlog.Flog("docker", "[docker command result]", output)

	items := strings.Split(strings.TrimSpace(output), "\n")

	if len(items) < 1 {
		return containers, errors.New("未发现容器")
	}
	for _, item := range items {
		arr := strings.Split(item, "|")
		container := DockerContainer{ID: arr[0], Name: arr[1], IP: arr[2], Status: arr[3], Ports: arr[4]}
		containers = append(containers, container)
	}

	return
}

func StartContainer(machine model.Machine, id string) (output string, err error) {
	conn, err := mssh.Connect(machine)
	defer conn.Close()
	if err != nil {
		return "", errors.New("无法连接主机")
	}
	action := "docker start " + id
	mlog.Flog("docker", "[docker command run]", action)
	output, err = cmd.RunRemote(conn, action)
	mlog.Flog("docker", "[docker command result]", output)
	if err != nil {
		return "", err
	}
	return
}

func StopContainer(machine model.Machine, id string) (output string, err error) {
	conn, err := mssh.Connect(machine)
	defer conn.Close()
	if err != nil {
		return "", errors.New("无法连接主机")
	}
	action := "docker stop " + id
	mlog.Flog("docker", "[docker command run]", action)
	output, err = cmd.RunRemote(conn, action)
	mlog.Flog("docker", "[docker command result]", output)
	if err != nil {
		return "", err
	}
	return
}

func proxy() {
	host := "0.0.0.0"
	port := "8888"
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		fmt.Println(err, err.Error())
		os.Exit(0)
	}

	for {
		s_conn, err := l.Accept()
		if err != nil {
			continue
		}

		d_tcpAddr, _ := net.ResolveTCPAddr("tcp4", "172.18.0.2:80")
		d_conn, err := net.DialTCP("tcp", nil, d_tcpAddr)
		if err != nil {
			fmt.Println(err)
			s_conn.Write([]byte("can't connect 172.17.0.2:80"))
			s_conn.Close()
			continue
		}
		go io.Copy(s_conn, d_conn)
		go io.Copy(d_conn, s_conn)
	}
}
