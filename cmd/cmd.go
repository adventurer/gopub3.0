package cmd

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

// run command with local
func RunLocal(command string) {
	cmd := exec.Command("/bin/bash", "-c", command)

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		fmt.Println("Execute failed when Start:" + err.Error())
		return
	}

	stdin.Close()

	outBytes, _ := ioutil.ReadAll(stdout)
	stdout.Close()

	errBytes, _ := ioutil.ReadAll(stderr)
	stderr.Close()

	if err := cmd.Wait(); err != nil {
		fmt.Println("Execute failed when Wait:" + err.Error() + ":{" + strings.TrimSpace(string(errBytes)) + "}")
		return
	}

	fmt.Println("Execute finished:" + string(outBytes))
}

func StandServ() {
	go RunLocal("ssh -NL 3001:localhost:3000 pazu@192.168.1.204")
	go RunLocal("ssh -NR 3000:localhost:3001 root@sz")
}
