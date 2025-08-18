package logger

import (
	"fmt"
	"ghostminion/db"
	"log"
	"os"
	"sync"
)

type logEntry struct {
	level   LogLevel
	message string
}

type Logger struct {
	instance *log.Logger
	file     *os.File
	level    LogLevel
	queue    chan logEntry
	wg       sync.WaitGroup
}

var (
	logger *Logger
	once   sync.Once
)

func GetLogger() *Logger {
	once.Do(func() {
		var f *os.File
		var err error

		if err != nil || f == nil {
			f = os.Stderr
		}

		l := &Logger{
			instance: log.New(f, "", log.LstdFlags|log.Lshortfile),
			file:     f,
			level:    DEBUG,
			queue:    make(chan logEntry, 100),
		}

		l.wg.Add(1)
		go l.worker()

		logger = l
	})
	return logger
}

func (l *Logger) worker() {
	defer l.wg.Done()
	for entry := range l.queue {
		if err := db.WriteLog(entry.level.String(), []byte(entry.message)); err != nil {
			return
		}
	}
}

func (l *Logger) log(level LogLevel, v ...any) {
	if level < l.level {
		return
	}
	msg := fmt.Sprint(v...)

	l.instance.SetPrefix("[" + level.String() + "] ")
	l.instance.Println(msg)

	select {
	case l.queue <- logEntry{level: level, message: msg}:
	default:
		l.instance.Printf("[WARN] Log queue full, dropping log: %s", msg)
	}
}

func (l *Logger) Debug(v ...any) { l.log(DEBUG, v...) }
func (l *Logger) Info(v ...any)  { l.log(INFO, v...) }
func (l *Logger) Warn(v ...any)  { l.log(WARN, v...) }
func (l *Logger) Error(v ...any) { l.log(ERROR, v...) }

func (l *Logger) Close() error {
	close(l.queue)
	l.wg.Wait()
	if l.file != nil && l.file != os.Stderr {
		return l.file.Close()
	}
	return nil
}
