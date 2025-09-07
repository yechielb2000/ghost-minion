package apps

import (
	"context"
	"errors"
	"ghostminion/db"
	"ghostminion/db/dbDataTypes"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type PeriodicCommandParams struct {
	Command  string `json:"Command"`
	Timeout  int    `json:"Timeout"`
	Interval int    `json:"Interval"`
	MaxRuns  int    `json:"MaxRuns"` // -1 = infinite
}

type PeriodicCommandApp struct {
	baseApp     *BaseApp
	params      *PeriodicCommandParams
	commandChan chan []byte
	wg          sync.WaitGroup
	once        sync.Once
}

func NewPeriodicCommand(appData AppData) (*PeriodicCommandApp, error) {
	params := PeriodicCommandParams{}
	if err := appData.UnmarshalParams(&params); err != nil {
		return nil, err
	}
	app := &PeriodicCommandApp{
		baseApp: &BaseApp{
			stop:    make(chan struct{}, 1),
			AppData: &appData,
		},
		params:      &params,
		commandChan: make(chan []byte, 100),
	}

	if err := app.validateParams(); err != nil {
		return nil, err
	}

	return app, nil
}

func (c *PeriodicCommandApp) Name() string {
	return c.baseApp.Name
}

func (c *PeriodicCommandApp) Start(ctx context.Context) error {
	// run store
	done := make(chan struct{})
	go func() {
		defer close(done)
		c.store(ctx)
	}()

	// run produce in this goroutine (blocking)
	c.produce(ctx)

	<-done
	return nil
}
func (c *PeriodicCommandApp) Stop() error {
	close(c.commandChan)
	return nil
}

func (c *PeriodicCommandApp) validateParams() error {
	if c.params.Command == "" {
		return errors.New("command must be provided")
	}
	if c.params.Interval <= 0 {
		return errors.New("interval must be greater than 0")
	}
	if c.params.MaxRuns < -1 {
		return errors.New("MaxRuns must be -1 or greater")
	}
	return nil
}

func (c *PeriodicCommandApp) produce(ctx context.Context) {
	runCount := 0
	ticker := time.NewTicker(time.Duration(c.params.Interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if c.params.MaxRuns != -1 && runCount >= c.params.MaxRuns {
				lgr.Info("Reached MaxRuns, stopping command produce")
				return
			}

			lgr.Info("Running periodic command:", c.params.Command)
			output, err := RunCommand(c.params.Command)
			if err != nil {
				lgr.Error("Error running command: " + err.Error())
			}

			select {
			case c.commandChan <- output:
			case <-ctx.Done():
				return
			}

			runCount++
		}
	}
}

func (c *PeriodicCommandApp) store(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case output, ok := <-c.commandChan:
			if !ok {
				return
			}
			if err := db.GetInstance().WriteData("", dbDataTypes.Commands, output); err != nil {
				lgr.Error("Error writing command output: " + err.Error())
			}
		}
	}
}

func RunCommand(command string) ([]byte, error) {
	cmd := exec.Command("sh", "-c", command) // more portable
	cmd.SysProcAttr = &syscall.SysProcAttr{ParentProcess: 0}
	output, err := cmd.CombinedOutput()
	return output, err
}
