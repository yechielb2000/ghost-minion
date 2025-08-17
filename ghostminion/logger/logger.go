package logger

import (
	"log"
	"os"
	"sync"
)

type Logger struct {
	instance *log.Logger
	file     *os.File
}

var (
	logger *Logger
	once   sync.Once
)

func Init(filePath string) (*Logger, error) {
	var err error
	once.Do(func() {
		var f *os.File
		f, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return
		}

		logger = &Logger{
			instance: log.New(f, "", log.LstdFlags|log.Lshortfile),
			file:     f,
		}
	})

	return logger, err
}

func GetLogger() *Logger {
	return logger
}

func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Example usage
// func main() {
//     l, err := logger.Init("app.log")
//     if err != nil {
//         panic(err)
//     }
//     defer l.Close()
//
//     log := logger.GetLogger()
//     log.instance.Println("Application started")
// }
