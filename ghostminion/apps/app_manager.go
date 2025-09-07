package apps

import (
	"context"
	"errors"
	"sync"
)

type AppManager struct {
	mu   sync.Mutex
	apps map[string]App
}

func NewAppManager() *AppManager {
	return &AppManager{
		apps: make(map[string]App),
	}
}

func (am *AppManager) Register(a App) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.apps[a.Name()] = a
}

// GetApp returns the app by name. Caller should not modify the returned app.
func (am *AppManager) GetApp(name string) (App, error) {
	am.mu.Lock()
	defer am.mu.Unlock()
	a := am.apps[name]
	if a == nil {
		return nil, errors.New("am: app not found")
	}
	return a, nil
}

// ListApps returns a shallow copy of the apps map (safe to iterate).
func (am *AppManager) ListApps() map[string]App {
	am.mu.Lock()
	defer am.mu.Unlock()
	out := make(map[string]App, len(am.apps))
	for k, v := range am.apps {
		out[k] = v
	}
	return out
}

// StartAll starts every registered app in its own goroutine and registers them with the provided WaitGroup.
// The wg.Done() is executed by the wrapper goroutine to ensure correctness.
func (am *AppManager) StartAll(ctx context.Context, wg *sync.WaitGroup) {
	am.mu.Lock()
	names := make([]string, 0, len(am.apps))
	for name := range am.apps {
		names = append(names, name)
	}
	am.mu.Unlock()

	for _, name := range names {
		am.StartApp(name, ctx, wg)
	}
}

// StartApp starts a single app (if found). It's safe to call concurrently.
func (am *AppManager) StartApp(name string, ctx context.Context, wg *sync.WaitGroup) {
	am.mu.Lock()
	app := am.apps[name]
	am.mu.Unlock()

	if app == nil {
		lgr.Warn("am: StartApp: app " + name + " not found")
		return
	}

	// increment wg and run the app in a goroutine; the wrapper is responsible for wg.Done()
	wg.Add(1)
	go func(a App) {
		defer wg.Done()

		if err := a.Start(ctx); err != nil {
			lgr.Error("app [" + a.Name() + "] exited with error: " + err.Error())
		} else {
			lgr.Info("app [" + a.Name() + "] exited cleanly")
		}
	}(app)
}

// StopApp calls Stop() on the app and removes it from the registry.
// It does NOT call wg.Done(): the goroutine running Start must return to release the wg.
func (am *AppManager) StopApp(name string) {
	am.mu.Lock()
	app := am.apps[name]
	if app == nil {
		am.mu.Unlock()
		lgr.Warn("am: StopApp: app " + name + " not found")
		return
	}
	// remove from registry before calling Stop() to prevent races with Start/Stop callers
	delete(am.apps, name)
	am.mu.Unlock()

	if err := app.Stop(); err != nil {
		lgr.Warn("am: StopApp: app " + name + " Stop() returned error: " + err.Error())
	}
}

// StopAll stops all registered apps safely.
func (am *AppManager) StopAll() {
	am.mu.Lock()
	names := make([]string, 0, len(am.apps))
	for name := range am.apps {
		names = append(names, name)
	}
	am.mu.Unlock()

	for _, name := range names {
		am.StopApp(name)
	}
}
