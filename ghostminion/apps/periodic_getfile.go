package apps

import (
	"errors"
	"fmt"
	"ghostminion/db"
	"os"
	"sync"
	"time"
)

type PeriodicGetFileApp struct {
	Path     string `json:"Path"`
	MaxSize  int    `json:"MaxSize"`
	Interval int    `json:"Interval"`
	CheckMD5 bool   `json:"CheckMD5"`
}

var stopPeriodicGetFileApp = false

func (c *PeriodicGetFileApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	currentFileMD5 := ""
	for stopPeriodicGetFileApp != true {
		fileContent, err := GetFile(c.Path)
		if err != nil {
			fmt.Println("error: ", err)
			continue
		}
		if c.CheckMD5 {
			command := fmt.Sprintf("md5sum %v", c.Path)
			fileMD5Output, err := RunCommand(command)
			if err != nil {
				fmt.Println("error calculating MD5: ", err)
				continue
			}
			fileMD5 := string(fileMD5Output) // Convert output to a string
			if currentFileMD5 != fileMD5 {
				fmt.Println("File content changed. MD5: ", fileMD5)
				currentFileMD5 = fileMD5
				err = db.WriteDataRow("", db.FilesDataType, fileContent) // replace requestId
				if err != nil {
					fmt.Printf("error: %v", err)
				}
			}
		}
		time.Sleep(time.Duration(c.Interval) * time.Second)
	}
}

func (c *PeriodicGetFileApp) Stop() error {
	stopPeriodicGetFileApp = true
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

func GetFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}
