package cfpsql

import (
	"code.cloudfoundry.org/cli/cf/errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

//go:generate counterfeiter . PsqlRunner
type PsqlRunner interface {
	RunPsql(hostname string, port int, dbName string, username string, password string, args ...string) error
}

func NewPsqlRunner(execWrapper ExecWrapper) PsqlRunner {
	return &psqlRunner{
		execWrapper: execWrapper,
	}
}

type psqlRunner struct {
	execWrapper ExecWrapper
}

func (self *psqlRunner) RunPsql(hostname string, port int, dbName string, username string, password string, psqlArgs ...string) error {
	path, err := self.execWrapper.LookPath("psql")
	if err != nil {
		return errors.New("'psql' client not found in PATH")
	}
	cmd := self.MakePsqlCommand(path, hostname, port, dbName, username, psqlArgs...)
	additionalEnv := "PGPASSWORD="+password
	newEnv := append(os.Environ(), additionalEnv)
	cmd.Env = newEnv

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = self.execWrapper.Run(cmd)
	if err != nil {
		return fmt.Errorf("error running psql client: %s", err)
	}

	return nil
}

func (self *psqlRunner) MakePsqlCommand(path string, hostname string, port int, dbName string, username string, psqlArgs ...string) *exec.Cmd {
	args := []string{"-U", username, "-h", hostname, "-p", strconv.Itoa(port), "-d", dbName}
	args = append(args, psqlArgs...)
	return exec.Command(path, args...)
}
