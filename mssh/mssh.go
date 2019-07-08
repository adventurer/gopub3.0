package mssh

import (
	"io"
	"io/ioutil"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

var Client *ssh.Client

func Init() {
	var err error
	b, err := ioutil.ReadFile("/Users/wuyang/.ssh/id_rsa")
	if err != nil {
		log.Println(err)
	}
	pKey, err := ssh.ParsePrivateKey(b)
	if err != nil {
		log.Println(err)
	}
	config := ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(pKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	Client, err = ssh.Dial("tcp", "120.79.161.99:22", &config)
	if err != nil {
		log.Println(err)
	}
	log.Println("ssh连接服务器成功")

	server, err := Client.Listen("tcp", "127.0.0.1:3000")

	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("开始接受连接")
	for {
		log.Println("接受连接前")
		client_s, err := server.Accept()
		log.Println("接受连接后")

		if err == nil {
			log.Println("开始处理连接")
			go handleClientRequest(client_s)
		}
		if err != nil {
			log.Println("处理连接错误：", err.Error())
		}
	}

}

// func handleClientRequest(client net.Conn) {
// 	defer client.Close()
// 	client.Write([]byte("hello word"))
// }

func handleClientRequest(client net.Conn) {
	defer client.Close()

	remote, err := net.Dial("tcp", ":3000")
	if err != nil {
		log.Println("转发错误：", err.Error())
		return
	}
	defer remote.Close()

	go io.Copy(remote, client)
	io.Copy(client, remote)

}
