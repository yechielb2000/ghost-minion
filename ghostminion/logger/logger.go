package logger

import (
	"ghostminion/config"
	"log"
	"os"
	"sync"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

type Logger struct {
	instance *log.Logger
	file     *os.File
	level    LogLevel
}

var (
	logger *Logger
	once   sync.Once
)

func GetLogger() *Logger {
	once.Do(func() {
		c, _ := config.GetConfig()
		var f *os.File
		var err error

		if c != nil && c.Installation.LogFilePath != "" {
			f, err = os.OpenFile(c.Installation.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		}

		if err != nil || f == nil {
			f = os.Stderr
		}

		logger = &Logger{
			instance: log.New(f, "", log.LstdFlags|log.Lshortfile),
			file:     f,
		}
	})
	return logger
}

func (lgr *Logger) Debug(v ...any) { lgr.instance.SetPrefix("[DEBUG] "); lgr.instance.Println(v...) }
func (lgr *Logger) Info(v ...any)  { lgr.instance.SetPrefix("[INFO] "); lgr.instance.Println(v...) }
func (lgr *Logger) Warn(v ...any)  { lgr.instance.SetPrefix("[WARN] "); lgr.instance.Println(v...) }
func (lgr *Logger) Error(v ...any) { lgr.instance.SetPrefix("[ERROR] "); lgr.instance.Println(v...) }

func (lgr *Logger) Close() error {
	if lgr.file != nil && lgr.file != os.Stderr {
		return lgr.file.Close()
	}
	return nil
}
