package cmd

import (
	"fmt"

	"golang.org/x/crypto/ssh"

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
		fmt.Println("Execute failed when Wait:" + err.Error() + ":{" + strings.TrimSpace(string(errBytes)) + "}")
		return "", err
	}

	// fmt.Println("Execute finished:" + string(outBytes))
	return string(outBytes), nil
}

// run remote command
func RunRemote(session *ssh.Session, command string) (result string, err error) {
	stdin, _ := session.StdinPipe()
	stdout, _ := session.StdoutPipe()
	stderr, _ := session.StderrPipe()

	if err = session.Run(command); err != nil {
		errBytes, _ := ioutil.ReadAll(stderr)
		return err.Error() + ":" + string(errBytes), err
	}

	stdin.Close()

	outBytes, _ := ioutil.ReadAll(stdout)

	return string(outBytes), nil
}
