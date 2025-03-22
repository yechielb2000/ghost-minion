package apps

import (
	"fmt"
	"ghostminion/db"
	"ghostminion/executor"
	"sync"
	"time"
)

type PeriodicCommandApp struct {
	Command  string
	Timeout  uint // in seconds default is 120
	Interval uint // in seconds
}

func (c *PeriodicCommandApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		fmt.Println("Running command: ", c.Command)
		commandOutput, err := executor.RunCommand(c.Command)
		if err != nil {
			fmt.Println("error: ", err)
		}
		err = db.WriteDataRow("reqId", db.ScreenshotsDataType, commandOutput) // replace requestId
		time.Sleep(time.Duration(c.Interval) * time.Second)
	}
}

func (c *PeriodicCommandApp) Stop() error {
	fmt.Println("Stopping PeriodicCommand app.")
	return nil
}

func (c *PeriodicCommandApp) Validate() error {
	return nil
}
