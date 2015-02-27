package command

import (
	"10gen.com/mci"
	"fmt"
	"github.com/10gen-labs/slogger/v1"
	"io"
	"os/exec"
	"strings"
)

type RemoteCommand struct {
	CmdString string

	Stdout io.Writer
	Stderr io.Writer

	// info necessary for sshing into the remote host
	RemoteHostName string
	User           string
	Options        []string
	Background     bool

	// set after the command is started
	Cmd *exec.Cmd
}

func (self *RemoteCommand) Run() error {
	err := self.Start()
	if err != nil {
		return err
	}
	return self.Cmd.Wait()
}

func (self *RemoteCommand) Wait() error {
	return self.Cmd.Wait()
}

func (self *RemoteCommand) Start() error {

	// build the remote connection, in user@host format
	remote := self.RemoteHostName
	if self.User != "" {
		remote = fmt.Sprintf("%v@%v", self.User, remote)
	}

	// build the command
	cmdArray := append(self.Options, remote)

	// set to the background, if necessary
	cmdString := self.CmdString
	if self.Background {
		cmdString = fmt.Sprintf("nohup %v >& /tmp/start &", cmdString)
	}
	cmdArray = append(cmdArray, cmdString)

	mci.Logger.Logf(slogger.WARN, "Remote command executing: '%#v'",
		strings.Join(cmdArray, " "))

	// set up execution
	cmd := exec.Command("ssh", cmdArray...)
	cmd.Stdout = self.Stdout
	cmd.Stderr = self.Stderr

	// cache the command running
	self.Cmd = cmd
	return cmd.Start()
}

func (self *RemoteCommand) Stop() error {
	return self.Cmd.Process.Kill()
}