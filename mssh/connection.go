package mssh

import (
	"log"
	"time"

	"golang.org/x/crypto/ssh"
	"gopub3.0/mlog"
	"gopub3.0/model"
)

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
	client, err := ssh.Dial("tcp", machine.Ip+":"+machine.Port, &config)
	if err != nil {
		mlog.Mlog.Println(err.Error())
		return nil, err
	}

	session, err = client.NewSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}
