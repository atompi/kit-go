package os

import (
	"bytes"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"strconv"
	"syscall"

	"go.uber.org/zap"
)

func ExecCmd(command string, username string) (result string, err error) {
	cmd := exec.Command("/bin/sh", "-c", command)

	stdoutBuf, stderrBuf := new(bytes.Buffer), new(bytes.Buffer)
	cmd.Stdout = stdoutBuf
	cmd.Stderr = stderrBuf

	userObj, err := user.Lookup(username)
	if err != nil {
		return
	}

	uid, _ := strconv.Atoi(userObj.Uid)
	gid, _ := strconv.Atoi(userObj.Gid)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}

	if err = cmd.Run(); err != nil {
		return
	}

	result = "stdout:\n" + stdoutBuf.String() + "\nstderr:\n" + stderrBuf.String()
	return
}

func GracefulExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-c
	zap.S().Warnf("a %v signal is received, exiting...", s)
	os.Exit(0)
}
