package mssh

import (
	"log"
	"time"

	"golang.org/x/crypto/ssh"
	"gopub3.0/mlog"
	"gopub3.0/model"
)

var clientPool = make(map[string]*ssh.Client, 10)

func init() {
	go func() {
		time.Sleep(1 * time.Second)
		for {
			for k, client := range clientPool {
				_, _, err := client.SendRequest("ping", true, []byte("ping"))
				if err != nil {
					delete(clientPool, k)
					mlog.Mlog.Println("ssh connected failed at machine:", k)
					continue
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()

}

func reconnect() {

}

func Connect(machine model.Machine) (session *ssh.Session, err error) {
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
		Timeout:         30 * time.Second,
	}

	client, ok := clientPool[machine.Name]
	if ok {
		session, err = client.NewSession()
		return session, nil
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
