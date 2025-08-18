package dbDataTypes

type DataType string

const (
	Files       DataType = "file"
	Commands    DataType = "command"
	Keyloggers  DataType = "keylogger"
	Screenshots DataType = "screenshot"
)
