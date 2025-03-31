package persistence

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"os"
	"sync"
)

var (
	targetID string
	once     sync.Once
)

func GenerateTargetID() string {
	once.Do(getTargetID)
	return targetID
}

func getTargetID() {
	var id string

	file, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		fmt.Println(err)
		file, err = os.ReadFile("/var/lib/dbus/machine-id")
		if err != nil {
			fmt.Println(err)
		}
	}

	if file == nil {
		id = rand.Text()
	} else {
		id = string(file)
	}

	md5id := md5.Sum([]byte(id))
	targetID = string(md5id[:])
}
