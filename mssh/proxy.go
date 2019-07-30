package mssh

import (
	"io"
	"net"
	"time"

	"gopub3.0/mlog"

	"golang.org/x/crypto/ssh"
	"gopub3.0/model"
)

type ProxyConn struct {
	Name      string
	User      string
	Host      string
	Port      string
	Rsa       string
	LocalIp   string
	LocalPort string
}

var Servers = make(map[string]ProxyConn)
var cpoll = make(map[string]clientPoll, 100)

var proxyConnChan = make(chan ProxyConn, 10)

type clientPoll struct {
	sshClient *ssh.Client
	status    bool
}

func init() {
	machines := []model.Machine{}
	services := []model.Service{}
	model.DB.Find(&machines)
	model.DB.Find(&services)
	for _, machine := range machines {
		for _, service := range services {
			if service.Auto == 1 {
				proxyMachine := ProxyConn{Name: machine.Name, User: machine.User, Host: machine.Ip, Port: machine.Port, Rsa: machine.Rsa, LocalIp: service.Ip, LocalPort: service.Port}
				proxyConnChan <- proxyMachine
			}
		}
	}

}

func Begin() {
	go func() {
		for {
			proxyCon := <-proxyConnChan
			go proxyStart(proxyCon)
		}
	}()

}

func NewProxy(proxyCon ProxyConn) {
	proxyConnChan <- proxyCon
}

func StopProxy(hostName string) {
	cpoll[hostName].sshClient.Close()
	cpoll[hostName] = clientPoll{sshClient: cpoll[hostName].sshClient, status: false}
}

func GetStatus(hostName string) bool {
	return cpoll[hostName].status
}

// 120.79.161.99:22
// 39.108.114.213:22
func proxyStart(hostInfo ProxyConn) {
	var err error
	pKey, err := ssh.ParsePrivateKey([]byte(hostInfo.Rsa))
	if err != nil {
		mlog.Mlog.Println(err)
	}
	config := ssh.ClientConfig{
		User: hostInfo.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(pKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}
	client, err := ssh.Dial("tcp", hostInfo.Host+":"+hostInfo.Port, &config)
	if err != nil {
		mlog.Mlog.Println(err.Error())
		return
	}
	cpoll[hostInfo.Name] = clientPoll{sshClient: client, status: true}

	// Client, err = ssh.Dial("tcp", "39.108.114.213:22", &config)

	if err != nil {
		mlog.Mlog.Println("sshdial err:", err)
	}

	server, err := cpoll[hostInfo.Name].sshClient.Listen("tcp", "127.0.0.1:"+hostInfo.LocalPort)

	if err != nil {
		mlog.Mlog.Println("listen err:", err.Error())
		cpoll[hostInfo.Name].sshClient.Close()
		cpoll[hostInfo.Name] = clientPoll{sshClient: client, status: false}
		return
	}

	for {
		mlog.Mlog.Println(hostInfo.Host+":"+hostInfo.Port, "发生数据")
		serverConn, err := server.Accept()
		if err == nil {
			go handleClientRequest(serverConn, hostInfo)
		}
		if err != nil {
			mlog.Mlog.Println("处理连接错误：", err.Error())
			break
		}

	}
}

func handleClientRequest(serverConn net.Conn, hostInfo ProxyConn) {
	defer serverConn.Close()

	localConn, err := net.Dial("tcp", hostInfo.LocalIp+":"+hostInfo.LocalPort)
	if err != nil {
		mlog.Mlog.Println("连接git服务器错误：", err.Error())
		return
	}
	defer localConn.Close()

	go io.Copy(localConn, serverConn)
	io.Copy(serverConn, localConn)

}
