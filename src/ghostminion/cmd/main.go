package main

import (
	"fmt"
	"ghostminion/core/executor"
	"log"
)

func main() {
	commandOutput, err := executor.RunCommand("cat /etc/passwd")
	if err != nil {
		log.Printf("Command error: %v\nstatus: %v", string(commandOutput), err.Error())
	} else {
		fmt.Println(string(commandOutput))
	}
}
