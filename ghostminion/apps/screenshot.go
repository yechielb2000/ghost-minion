package apps

import (
	"bytes"
	"fmt"
	"ghostminion/db"
	"github.com/kbinani/screenshot"
	"image"
	"image/jpeg"
	"log"
	"strconv"
	"sync"
	"time"
)

type ScreenshotApp struct {
	Interval uint // in seconds
}

func (c *ScreenshotApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		c.runScreenshot()
		time.Sleep(time.Duration(c.Interval) * time.Second)
	}
}

func (c *ScreenshotApp) Stop() error {
	fmt.Println("Stopping Screenshot app.")
	return nil
}

func (c *ScreenshotApp) Validate() error {
	return nil
}

func (c *ScreenshotApp) runScreenshot() {
	numOfActiveDisplays := screenshot.NumActiveDisplays()
	if numOfActiveDisplays <= 0 {
		log.Fatalf("Active display not found")
		return
	}

	all := image.Rect(0, 0, 0, 0)

	for i := 0; i < numOfActiveDisplays; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		all = bounds.Union(all)
		captureTime := time.Now().Unix()
		fmt.Println("Screenshot captured at ", captureTime, " for display ", i)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			fmt.Printf("error: %v", err)
			continue
		}
		var buf bytes.Buffer
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
		if err != nil {
			fmt.Printf("error: %v", err)
		}
		err = db.WriteDataRow(strconv.Itoa(i), db.ScreenshotsDataType, buf.Bytes()) // replace requestId
		if err != nil {
			fmt.Printf("error: %v", err)
		}
	}
}
