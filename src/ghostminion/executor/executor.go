package executor

import (
	"os"
	"os/exec"
	"syscall"
)

func RunCommand(command string) ([]byte, error) {
	const ShellToUse = "bash"
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{ParentProcess: 0}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, err
	}
	return output, nil
}

func GetFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}
