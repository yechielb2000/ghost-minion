package apps

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"ghostminion/db"
	"ghostminion/db/dbDataTypes"
	"os"
	"sync"
	"time"
)

type PeriodicGetFileApp struct {
	Path     string `json:"Path"`
	MaxSize  int    `json:"MaxSize"`
	Interval int    `json:"Interval"`
	CheckMD5 bool   `json:"CheckMD5"`
	MaxRuns  int    `json:"MaxRuns"` // -1 for infinite

	stop      chan struct{}
	eventChan chan []byte
	gfwg      sync.WaitGroup
	once      sync.Once
}

func NewPeriodicGetFileApp(path string, maxSize, interval int, checkMD5 bool, maxRuns int) (*PeriodicGetFileApp, error) {
	app := &PeriodicGetFileApp{
		Path:      path,
		MaxSize:   maxSize,
		Interval:  interval,
		CheckMD5:  checkMD5,
		MaxRuns:   maxRuns,
		stop:      make(chan struct{}),
		eventChan: make(chan []byte, 100),
	}

	if err := app.validateParams(); err != nil {
		return nil, err
	}

	return app, nil
}

func (c *PeriodicGetFileApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	c.gfwg.Add(2)
	go c.startProducer()
	go c.startConsumer()
	c.gfwg.Wait()
}

func (c *PeriodicGetFileApp) Stop() {
	c.once.Do(func() {
		close(c.stop)
		c.gfwg.Wait()
		close(c.eventChan)
	})
}

func (c *PeriodicGetFileApp) validateParams() error {
	if c.Path == "" {
		return errors.New("path must be provided")
	}
	if c.MaxSize < 0 {
		return errors.New("max size must be 0 or greater")
	}
	if c.Interval <= 0 {
		return errors.New("interval must be greater than 0")
	}
	if c.MaxRuns < -1 {
		return errors.New("MaxRuns must be -1 or greater")
	}
	return nil
}

func (c *PeriodicGetFileApp) startProducer() {
	defer c.gfwg.Done()
	var currentMD5 string
	runCount := 0

	ticker := time.NewTicker(time.Duration(c.Interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.stop:
			return
		case <-ticker.C:
			if c.MaxRuns != -1 && runCount >= c.MaxRuns {
				lgr.Info("Reached MaxRuns, stopping producer")
				return
			}

			data, err := os.ReadFile(c.Path)
			if err != nil {
				lgr.Error("Error reading file %s: %v", c.Path, err)
				continue
			}

			if c.MaxSize > 0 && len(data) > c.MaxSize {
				lgr.Warn("File %s exceeds MaxSize (%d bytes), skipping", c.Path, c.MaxSize)
				continue
			}

			if c.CheckMD5 {
				hash := md5.Sum(data)
				fileMD5 := hex.EncodeToString(hash[:])
				if fileMD5 == currentMD5 {
					continue
				}
				currentMD5 = fileMD5
			}

			select {
			case c.eventChan <- data:
			default:
				lgr.Warn("Event channel full, dropping file data")
			}
			runCount++
		}
	}
}

func (c *PeriodicGetFileApp) startConsumer() {
	defer c.gfwg.Done()
	for {
		select {
		case <-c.stop:
			return
		case data, ok := <-c.eventChan:
			if !ok {
				return
			}
			if err := db.GetInstance().WriteData("", dbDataTypes.Files, data); err != nil {
				lgr.Error("Error writing file data: %v", err)
			}
		}
	}
}
