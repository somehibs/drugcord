package drugcord

import (
	"fmt"
)

type FatalError struct {
	Cause error
	Msg   string
}

func (f FatalError) Error() string {
	if f.Cause == f {
		return fmt.Sprintf("FATAL! %s\n", f.Msg)
	}
	return fmt.Sprintf("FATAL! %s %s\n", f.Msg, f.Cause.Error())
}

func Fatal(cause error, message string) FatalError {
	return FatalError{cause, message}
}
