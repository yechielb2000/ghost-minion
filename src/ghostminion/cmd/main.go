package main

import (
	"fmt"
	"ghostminion/executor"
	"log"
)

func main() {
	fileContent, err := executor.GetFile("/mnt/c/Users/yechi/Desktop/test.py")
	if err != nil {
		log.Printf("File error: %v", err.Error())
	} else {
		fmt.Println(string(fileContent))
	}
	commandOutput, err := executor.RunCommand("pwd")
	if err != nil {
		log.Printf("Command error: %v\nstatus: %v", string(commandOutput), err.Error())
	} else {
		fmt.Println(string(commandOutput))
	}
}
