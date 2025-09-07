package apps

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"ghostminion/db"
	"ghostminion/db/dbDataTypes"
	"io/ioutil"
	"sync"
	"time"
)

type PeriodicGetFileParams struct {
	Path     string `json:"Path"`
	MaxSize  int    `json:"MaxSize"`
	Interval int    `json:"Interval"`
	CheckMD5 bool   `json:"CheckMD5"`
	MaxRuns  int    `json:"MaxRuns"` // -1 for infinite
}

type PeriodicGetFileApp struct {
	baseApp   *BaseApp
	params    *PeriodicGetFileParams
	eventChan chan []byte
	gfwg      sync.WaitGroup
	once      sync.Once
}

func NewPeriodicGetFileApp(appData AppData) (*PeriodicGetFileApp, error) {
	params := PeriodicGetFileParams{}
	if err := appData.UnmarshalParams(&params); err != nil {
		return nil, err
	}
	app := &PeriodicGetFileApp{
		baseApp: &BaseApp{
			stop:    make(chan struct{}),
			AppData: &appData,
		},
		params:    &params,
		eventChan: make(chan []byte, 100),
	}

	if err := app.validateParams(); err != nil {
		return nil, err
	}

	return app, nil
}

func (c *PeriodicGetFileApp) Name() string {
	return c.baseApp.Name
}

func (c *PeriodicGetFileApp) Start(ctx context.Context) error {
	// run consumer in goroutine
	done := make(chan struct{})
	go func() {
		defer close(done)
		c.store(ctx)
	}()

	// run producer in this goroutine (blocking)
	c.produce(ctx)

	<-done
	return nil
}

func (c *PeriodicGetFileApp) Stop() error {
	close(c.eventChan)
	return nil
}

// validateParams ensures required fields are correct
func (c *PeriodicGetFileApp) validateParams() error {
	if c.params.Path == "" {
		return errors.New("path must be provided")
	}
	if c.params.MaxSize < 0 {
		return errors.New("max size must be 0 or greater")
	}
	if c.params.Interval <= 0 {
		return errors.New("interval must be greater than 0")
	}
	if c.params.MaxRuns < -1 {
		return errors.New("MaxRuns must be -1 or greater")
	}
	return nil
}

// produce reads file periodically and sends data to eventChan
func (c *PeriodicGetFileApp) produce(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(c.params.Interval) * time.Second)
	defer ticker.Stop()

	var currentMD5 string
	runCount := 0

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if c.params.MaxRuns != -1 && runCount >= c.params.MaxRuns {
				lgr.Info("Reached MaxRuns, stopping file producer")
				return
			}

			data, err := ioutil.ReadFile(c.params.Path)
			if err != nil {
				lgr.Error("Error reading file " + c.params.Path + ": " + err.Error())
				continue
			}

			if c.params.MaxSize > 0 && len(data) > c.params.MaxSize {
				lgr.Warn("File " + c.params.Path + " exceeds MaxSize, skipping")
				continue
			}

			if c.params.CheckMD5 {
				hash := md5.Sum(data)
				fileMD5 := hex.EncodeToString(hash[:])
				if fileMD5 == currentMD5 {
					continue
				}
				currentMD5 = fileMD5
			}

			select {
			case c.eventChan <- data: // blocks if channel full
			case <-ctx.Done():
				return
			}
			runCount++
		}
	}
}

// store writes file data from eventChan to DB
func (c *PeriodicGetFileApp) store(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-c.eventChan:
			if !ok {
				return
			}
			if err := db.GetInstance().WriteData("", dbDataTypes.Files, data); err != nil {
				lgr.Error("Error writing file data: " + err.Error())
			}
		}
	}
}
