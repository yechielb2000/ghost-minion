package executor

import (
	"os/exec"
)

const ShellToUse = "bash"

func RunCommand(command string) ([]byte, error) {
	cmd := exec.Command(ShellToUse, "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, err
	}
	return output, nil
}
