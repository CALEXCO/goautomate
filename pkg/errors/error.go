package errors

import "fmt"

// TODO: Define errors codes and use exit to finish the program
const (
	ConnectionFailure = iota * 125
	CommandTimeout
)

type GoAutomateErrors struct {
	Code int
	Msg  string
}

func (ga *GoAutomateErrors) Error() string {
	return fmt.Sprintf("ERROR: Exit Code: %d, msg: %s", ga.Code, ga.Msg)
}
