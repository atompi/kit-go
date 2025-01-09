package os

import (
	"bytes"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
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
