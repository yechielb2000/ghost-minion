package monitors

import (
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"unsafe"
)

func main() {
	fd, err := unix.InotifyInit()
	if err != nil {
		log.Fatal(err)
	}
	defer unix.Close(fd)

	// File you want to monitor
	file := "/etc/myconfig.conf"
	wd, err := unix.InotifyAddWatch(fd, file, unix.IN_MODIFY|unix.IN_ATTRIB|unix.IN_DELETE_SELF)
	if err != nil {
		log.Fatal(err)
	}
	defer unix.InotifyRmWatch(fd, uint32(wd))

	buf := make([]byte, 4096)
	for {
		n, err := unix.Read(fd, buf)
		if err != nil {
			log.Fatal(err)
		}

		// Parse events
		var offset uint32
		for offset <= uint32(n-unix.SizeofInotifyEvent) {
			raw := (*unix.InotifyEvent)(unsafe.Pointer(&buf[offset]))
			name := file
			fmt.Printf("[!] File event: %s mask=0x%x\n", name, raw.Mask)

			offset += unix.SizeofInotifyEvent + raw.Len
		}
	}
}
