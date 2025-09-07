package core

import (
	"context"
	"fmt"
	"ghostminion/apps"
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
	apps         *apps.AppManager
}

// NewCore creates a Core that owns the cancellation context and WaitGroup.
func NewCore() *Core {
	ctx, cancel := context.WithCancel(context.Background())
	return &Core{
		cancelReason: make(chan string, 1),
		ctx:          ctx,
		cancel:       cancel,
		apps:         apps.NewAppManager(),
	}
}

func (c *Core) Context() context.Context {
	return c.ctx
}

// RegisterApp registers an app into the manager (does not start it).
func (c *Core) RegisterApp(a apps.App) {
	c.apps.Register(a)
}

// Start starts all registered apps and blocks until a shutdown reason or external cancel.
// After shut down is triggered, it waits for all apps to finish.
func (c *Core) Start() {
	// Start all registered apps under the core context and wg.
	c.apps.StartAll(c.ctx, &c.wg)

	// Wait for a shutdown reason or external cancel
	select {
	case msg := <-c.cancelReason:
		lgr.Info(fmt.Sprintf("[core] shutdown requested: %s", msg))
		c.Shutdown(msg)
	case <-c.ctx.Done():
		lgr.Info("[core] context canceled externally")
		// ensure we still stop apps
		c.apps.StopAll()
	}

	// Wait for all apps to finish their cleanup
	c.wg.Wait()
	lgr.Info("[core] all apps stopped, core exiting")
}

// Shutdown triggers a global shutdown. It's important.
func (c *Core) Shutdown(reason string) {
	// log early (useful for C2 ack before we cancel)
	lgr.Warn(reason)

	// cancel context (idempotent)
	c.cancel()

	// ask apps to stop (non-blocking). The actual goroutines will return and the wg will be released.
	c.apps.StopAll()
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
