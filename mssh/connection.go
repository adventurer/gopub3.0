package mssh

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"gopub3.0/mlog"
	"gopub3.0/model"
)

var clientPool = make(map[string]*ssh.Client, 10)

func init() {
	// go func() {
	// 	time.Sleep(1 * time.Second)
	// 	for {
	// 		for k, client := range clientPool {
	// 			log.Println(k)
	// 			delete(clientPool, k)

	// 			_, _, err := client.SendRequest("ping", true, []byte("ping"))
	// 			if err != nil {
	// 				delete(clientPool, k)
	// 				mlog.Mlog.Println("ssh connected failed at machine:", k)
	// 				continue
	// 			}
	// 		}
	// 		time.Sleep(5 * time.Second)
	// 	}
	// }()

}

func reconnect() {

}

func Connect(machine model.Machine) (session *ssh.Session, err error) {
	defer recoverConnect()
	pKey, err := ssh.ParsePrivateKey([]byte(machine.Rsa))
	if err != nil {
		log.Println(err)
	}
	config := ssh.ClientConfig{
		User: machine.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(pKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	client, ok := clientPool[machine.Name]
	if ok {
		session, err = client.NewSession()
		if err != nil {
			delete(clientPool, machine.Name)
		} else {
			return session, nil
		}
	}
	clientPool[machine.Name], err = ssh.Dial("tcp", machine.Ip+":"+machine.Port, &config)
	if err != nil {
		mlog.Mlog.Println(err.Error())
		return nil, err
	}
	mlog.Mlog.Println("ssh connected to :", machine.Ip)

	session, err = clientPool[machine.Name].NewSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func ScpCopy(machine model.Machine, localFilePath, remoteDir string) error {
	var (
		sftpClient *sftp.Client
		err        error
	)
	// 这里换成实际的 SSH 连接的 用户名，密码，主机名或IP，SSH端口
	sftpClient, err = sftpconnect(machine)
	if err != nil {
		mlog.Mlog.Println("scpCopy:", err)
		return err
	}
	defer sftpClient.Close()
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		mlog.Mlog.Println("scpCopy:", err)
		return err
	}
	defer srcFile.Close()

	var remoteFileName = path.Base(localFilePath)
	dstFile, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err != nil {
		mlog.Mlog.Println("scpCopy:", err)
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 1024)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf[0:n])
	}
	return nil
}

func sftpconnect(machine model.Machine) (*sftp.Client, error) {
	pKey, err := ssh.ParsePrivateKey([]byte(machine.Rsa))
	if err != nil {
		log.Println(err)
	}
	config := ssh.ClientConfig{
		User: machine.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(pKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}
	sshClient, err := ssh.Dial("tcp", machine.Ip+":"+machine.Port, &config)
	if err != nil {
		return nil, err
	}
	// create sftp client
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}
	return sftpClient, nil
}

func recoverConnect() (session *ssh.Session, err error) {
	if r := recover(); r != nil {
		mlog.Mlog.Println("ssh连接崩溃，延迟5秒重启,并清空连接池", r)
		clientPool = make(map[string]*ssh.Client, 10)
		time.Sleep(5 * time.Second)
	}
	return session, fmt.Errorf("ssh连接崩溃，延迟5秒重启")
}
