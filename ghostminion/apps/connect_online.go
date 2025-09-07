package apps

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"sync"
)

type ConnectOnlineParams struct {
	Port     int    `json:"Port,omitempty"`
	Password string `json:"Password,omitempty"`
}

type ConnectOnlineApp struct {
	baseApp  *BaseApp
	params   *ConnectOnlineParams
	listener net.Listener
	wg       sync.WaitGroup
}

func NewConnectOnlineApp(appData AppData) (*ConnectOnlineApp, error) {
	params := &ConnectOnlineParams{}
	if err := appData.UnmarshalParams(params); err != nil {
		return nil, err
	}
	app := &ConnectOnlineApp{
		baseApp: &BaseApp{
			stop:    make(chan struct{}, 1),
			AppData: &appData,
		},
		params: params,
	}

	if err := app.validateParams(); err != nil {
		return nil, err
	}

	return app, nil
}

func (c *ConnectOnlineApp) Name() string {
	return c.baseApp.Name
}

// Start runs the TCP server until ctx is canceled
func (c *ConnectOnlineApp) Start(ctx context.Context) error {
	address := fmt.Sprintf(":%d", c.params.Port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		lgr.Error("Error starting server: " + err.Error())
		return err
	}
	defer ln.Close()

	lgr.Info("Server is listening on port:", c.params.Port)

	// Accept loop
	connChan := make(chan net.Conn)
	errChan := make(chan error)

	// Goroutine to accept connections
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					errChan <- err
					continue
				}
			}
			select {
			case connChan <- conn:
			case <-ctx.Done():
				conn.Close()
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			lgr.Info("Server shutting down")
			return nil
		case err := <-errChan:
			lgr.Error("Error accepting connection: " + err.Error())
		case conn := <-connChan:
			go c.handleConnection(ctx, conn)
		}
	}
}

// Stop is a no-op because context controls shutdown
func (c *ConnectOnlineApp) Stop() error {
	return nil
}

// validateParams ensures port and password are valid
func (c *ConnectOnlineApp) validateParams() error {
	if c.params.Port < 1 || c.params.Port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}
	if c.params.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

// handleConnection handles a single TCP client until exit or ctx cancellation
func (c *ConnectOnlineApp) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	lgr.Info("New connection from:", conn.RemoteAddr().String())

	reader := bufio.NewReader(conn)

	conn.Write([]byte("Enter password: "))
	pass, _ := reader.ReadString('\n')
	pass = strings.TrimSpace(pass)
	if pass != c.params.Password {
		conn.Write([]byte("Invalid password. Closing connection.\n"))
		lgr.Warn("Authentication failed from", conn.RemoteAddr().String())
		return
	}
	conn.Write([]byte("Authentication successful.\n"))

	for {
		conn.Write([]byte("> "))
		command, err := reader.ReadString('\n')
		if err != nil {
			lgr.Error("Error reading command: " + err.Error())
			return
		}

		command = strings.TrimSpace(command)

		select {
		case <-ctx.Done():
			return
		default:
		}

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
