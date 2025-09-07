package persistence

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"ghostminion/logger"
	"os"
	"sync"
)

var (
	targetID string
	once     sync.Once
	lgr      = logger.GetInstance()
)

func GeTargetID() string {
	once.Do(getTargetID)
	return targetID
}

func getTargetID() {
	var id string

	file, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		lgr.Error("Error reading /etc/machine-id:", err.Error())
		file, err = os.ReadFile("/var/lib/dbus/machine-id")
		if err != nil {
			lgr.Error("Error reading /var/lib/dbus/machine-id:", err.Error())
		}
	}

	if file == nil {
		lgr.Info("Generating random target id")
		id = rand.Text()
	} else {
		id = string(file)
	}

	md5id := md5.Sum([]byte(id))
	targetID = hex.EncodeToString(md5id[:])[:]
}
