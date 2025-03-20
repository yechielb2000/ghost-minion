package apps

import (
	"errors"
	"fmt"
	"ghostminion/executor"
	"sync"
	"time"
)

type PeriodicGetFileApp struct {
	Path     string
	MaxSize  int  // maximum allowed size in bytes.
	Interval int  // in seconds
	CheckMD5 bool // check if file was changed since last run.
}

func (c *PeriodicGetFileApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	currentFileMD5 := ""
	for {
		fileContent, err := executor.GetFile(c.Path)
		if err != nil {
			fmt.Println("error: ", err)
			continue
		}
		if c.CheckMD5 {
			command := fmt.Sprintf("md5sum %v", c.Path)
			fileMD5Output, err := executor.RunCommand(command)
			if err != nil {
				fmt.Println("error calculating MD5: ", err)
				continue
			}
			fileMD5 := string(fileMD5Output) // Convert output to a string
			if currentFileMD5 != fileMD5 {
				fmt.Println("File content changed. MD5: ", fileMD5)
				currentFileMD5 = fileMD5
				fmt.Println("file content: ", fileContent) // save this to db
			}
		}
		time.Sleep(time.Duration(c.Interval) * time.Second)
	}
}

func (c *PeriodicGetFileApp) Stop() error {
	fmt.Println("Stopping PeriodicGetFile app.")
	return nil
}

func (c *PeriodicGetFileApp) Validate() error {
	if c.Path == "" {
		return errors.New("path must be provided")
	}
	if c.MaxSize <= 0 {
		return errors.New("max size must be greater than 0")
	}
	if c.Interval <= 0 {
		return errors.New("interval must be greater than 0")
	}
	return nil
}
