package apps

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
)

type SecurityGuardApp struct {
	isSafe         bool
	mu             sync.Mutex
	FilesExistence []string
}

func (ctx *SecurityGuardApp) IsSafe() bool {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	return ctx.isSafe
}

func (ctx *SecurityGuardApp) SetIsSafe(isSafe bool) {
	ctx.mu.Lock()
	ctx.isSafe = isSafe
	ctx.mu.Unlock()
}

func (ctx *SecurityGuardApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Starting SecurityGuard app.")
	ctx.checkFileExistence(ctx.FilesExistence)
	/*
		isSafe should be false on one of these terms
		- it has been too much time without communicating with the C2 (3 days)
		- unknown process or user has touched the db or config file
		- someone wrote the process name in its bash history
		- any of the files that was supposed to be in its place is not anymore
		- the cpu of the target has increase too much because of our process
	*/
}

func (ctx *SecurityGuardApp) Stop() error {
	return nil
}

func (ctx *SecurityGuardApp) Validate() error {
	ctx.isSafe = true
	return nil
}

func (ctx *SecurityGuardApp) checkFileExistence(files []string) {
	for _, file := range files {
		if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
			log.Printf("File %s does not exist.", file)
			ctx.SetIsSafe(false)
			break
		}
	}
}
