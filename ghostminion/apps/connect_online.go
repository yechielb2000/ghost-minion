package apps

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type ConnectOnlineApp struct {
	Port     int    `json:"Port,omitempty"`
	Password string `json:"Password,omitempty"`

	listener net.Listener
	stop     chan struct{}
	wg       sync.WaitGroup
}

func NewConnectOnlineApp(port int, password string) (*ConnectOnlineApp, error) {
	app := &ConnectOnlineApp{
		Port:     port,
		Password: password,
	}

	if err := app.validateParams(); err != nil {
		return nil, err
	}

	return app, nil
}

func (c *ConnectOnlineApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	address := fmt.Sprintf(":%d", c.Port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		lgr.Error("Error starting server:", err)
		return
	}
	c.listener = ln
	c.stop = make(chan struct{})

	lgr.Info("Server is listening on port:", strconv.Itoa(c.Port))

	for {
		select {
		case <-c.stop:
			return
		default:
			conn, err := ln.Accept()
			if err != nil {
				select {
				case <-c.stop:
					return
				default:
					lgr.Error("Error accepting connection:", err)
					continue
				}
			}
			c.wg.Add(1)
			go func() {
				defer c.wg.Done()
				c.handleConnection(conn)
			}()
		}
	}
}

func (c *ConnectOnlineApp) Stop() {
	close(c.stop)
	if c.listener != nil {
		c.listener.Close()
	}
	c.wg.Wait()
	lgr.Info("Server stopped")
}

func (c *ConnectOnlineApp) validateParams() error {
	if c.Port < 1 || c.Port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}
	if c.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (c *ConnectOnlineApp) handleConnection(conn net.Conn) {
	defer conn.Close()
	lgr.Info("New connection from:", conn.RemoteAddr().String())

	reader := bufio.NewReader(conn)

	conn.Write([]byte("Enter password: "))
	pass, _ := reader.ReadString('\n')
	pass = strings.TrimSpace(pass)
	if pass != c.Password {
		conn.Write([]byte("Invalid password. Closing connection.\n"))
		lgr.Warn("Authentication failed from", conn.RemoteAddr().String())
		return
	}
	conn.Write([]byte("Authentication successful.\n"))

	for {
		conn.Write([]byte("> "))
		command, err := reader.ReadString('\n')
		if err != nil {
			lgr.Error("Error reading command:", err)
			break
		}

		command = strings.TrimSpace(command)

		switch command {
		case "":
			continue
		case "help":
			conn.Write([]byte("Available commands:\n"))
			conn.Write([]byte(" - run shell commands (e.g., ls, whoami)\n"))
			conn.Write([]byte(" - exit : close the connection\n"))
		case "exit":
			lgr.Info("Closing connection with", conn.RemoteAddr().String())
			return
		default:
			parts := strings.Fields(command)
			if len(parts) == 0 {
				continue
			}
			cmd := exec.Command(parts[0], parts[1:]...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				conn.Write([]byte(fmt.Sprintf("error executing command: %v\n", err)))
			}
			conn.Write(output)
		}
	}
}
