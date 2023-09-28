package os

import "os/exec"

func ExecCmd(command string) (output string, err error) {
	cmd := exec.Command("/bin/sh", "-c", command)
	outputBuf, err := cmd.Output()
	if err != nil {
		output = ""
		return
	}
	output = string(outputBuf)
	return
}
