package cmd

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
	"gopub3.0/mlog"

	"io/ioutil"
	"os/exec"
	"strings"
)

// run command with local
func RunLocal(command string) (output string, err error) {
	cmd := exec.Command("/bin/bash", "-c", command)

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	mlog.Flog("localCommand", "[local command run]", command)

	if err = cmd.Start(); err != nil {
		fmt.Println("Execute failed when Start:" + err.Error())
		return "", err
	}

	stdin.Close()

	outBytes, _ := ioutil.ReadAll(stdout)
	stdout.Close()

	errBytes, _ := ioutil.ReadAll(stderr)
	stderr.Close()

	if err = cmd.Wait(); err != nil {
		mlog.Flog("localCommand", "[local command result]", strings.TrimSpace(string(errBytes)))
		return "", errors.New(strings.TrimSpace(string(errBytes)))
	}
	mlog.Flog("localCommand", "[local command result]", string(outBytes))

	// fmt.Println("Execute finished:" + string(outBytes))
	return string(outBytes), nil
}

// run remote command
func RunRemote(session *ssh.Session, command string) (result string, err error) {
	defer recoverName()
	// stdin, _ := session.StdinPipe()
	stdout, _ := session.StdoutPipe()
	stderr, _ := session.StderrPipe()
	mlog.Flog("remoteCommand", "[remote command run]", command)

	// if err = session.Run(command); err != nil {
	// 	errBytes, _ := ioutil.ReadAll(stderr)
	// 	mlog.Flog("remoteCommand", "[remote command result]", string(errBytes))
	// 	return err.Error() + ":" + string(errBytes), err
	// }

	err = session.Run(command)
	if err != nil {
		errBytes, _ := ioutil.ReadAll(stderr)
		mlog.Flog("remoteCommand", "[remote command err]", err.Error()+":"+string(errBytes))
		return err.Error() + ":" + string(errBytes), err
	}
	outBytes, _ := ioutil.ReadAll(stdout)
	// stdin.Close()
	mlog.Flog("remoteCommand", "[remote command result]", string(outBytes))

	return string(outBytes), nil
}

func recoverName() {
	if r := recover(); r != nil {
		fmt.Println("远程命令崩溃，延迟5秒重启", r)
		time.Sleep(1 * time.Second)
	}
}
