package apps

import (
	"errors"
	"ghostminion/db"
	"ghostminion/db/dbDataTypes"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type PeriodicCommandApp struct {
	Command  string `json:"Command"`
	Timeout  int    `json:"Timeout"`
	Interval int    `json:"Interval"`
	MaxRuns  int    `json:"MaxRuns"` // -1 = infinite

	stop        chan struct{}
	commandChan chan []byte
	wg          sync.WaitGroup
	once        sync.Once
}

func NewPeriodicCommandApp(command string, interval int, timeout int, maxRuns int) (*PeriodicCommandApp, error) {
	app := &PeriodicCommandApp{
		Command:     command,
		Interval:    interval,
		Timeout:     timeout,
		MaxRuns:     maxRuns,
		stop:        make(chan struct{}),
		commandChan: make(chan []byte, 100),
	}

	if err := app.validateParams(); err != nil {
		return nil, err
	}

	return app, nil
}

func (c *PeriodicCommandApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	c.wg.Add(2)
	go c.startProducer()
	go c.startConsumer()
	c.wg.Wait()
}

func (c *PeriodicCommandApp) Stop() {
	c.once.Do(func() {
		close(c.stop)
		c.wg.Wait()
		close(c.commandChan)
	})
}

func (c *PeriodicCommandApp) validateParams() error {
	if c.Command == "" {
		return errors.New("command must be provided")
	}
	if c.Interval == 0 {
		return errors.New("interval must be greater than 0")
	}
	if c.MaxRuns < -1 {
		return errors.New("MaxRuns must be -1 or greater")
	}
	return nil
}

func (c *PeriodicCommandApp) startProducer() {
	defer c.wg.Done()
	runCount := 0

	ticker := time.NewTicker(time.Duration(c.Interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.stop:
			return
		case <-ticker.C:
			if c.MaxRuns != -1 && runCount >= c.MaxRuns {
				lgr.Info("Reached MaxRuns, stopping command producer")
				return
			}

			lgr.Info("Running periodic command:", c.Command)
			output, err := RunCommand(c.Command)
			if err != nil {
				lgr.Error("Error running command: ", err)
			}

			select {
			case c.commandChan <- output:
			default:
				lgr.Warn("Command channel full, dropping output")
			}

			runCount++
		}
	}
}

func (c *PeriodicCommandApp) startConsumer() {
	defer c.wg.Done()
	for {
		select {
		case <-c.stop:
			return
		case output, ok := <-c.commandChan:
			if !ok {
				return
			}
			if err := db.GetInstance().WriteData("", dbDataTypes.Commands, output); err != nil {
				lgr.Error("Error writing command output: ", err)
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
