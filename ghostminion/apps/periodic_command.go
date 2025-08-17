package apps

import (
	"fmt"
	"ghostminion/db"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type PeriodicCommandApp struct {
	Command  string `json:"Command"`
	Timeout  uint   `json:"Timeout"`
	Interval uint   `json:"Interval"`
}

var stopPeriodicCommandApp = false

func (c *PeriodicCommandApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	for stopPeriodicCommandApp != true {
		fmt.Println("Running command: ", c.Command)
		commandOutput, err := RunCommand(c.Command)
		if err != nil {
			fmt.Println("error: ", err)
		}
		err = db.WriteDataRow("", db.CommandsDataType, commandOutput) // replace requestId
		time.Sleep(time.Duration(c.Interval) * time.Second)
	}
}

func (c *PeriodicCommandApp) Stop() error {
	stopPeriodicCommandApp = true
	return nil
}

func (c *PeriodicCommandApp) Validate() error {
	return nil
}

func RunCommand(command string) ([]byte, error) {
	cmd := exec.Command(command)
	cmd.SysProcAttr = &syscall.SysProcAttr{ParentProcess: 0}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, err
	}
	return output, nil
}
