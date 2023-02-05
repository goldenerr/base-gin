package log

import (
	"fmt"
	"github.com/pkg/errors"
)

type Error struct {
	ErrorCode    int
	ErrorMsg     string
	ErrorUserMsg string
}

func NewError(code int, message string, userMessage string) Error {
	return Error{
		ErrorCode:    code,
		ErrorMsg:     message,
		ErrorUserMsg: userMessage,
	}
}

func (err Error) Error() string {
	return err.ErrorMsg
}

// WrapPrint print Error & User like
// Error {
//     ErrorCode: 400,
//     ErrorMsg: "something wrong with: %v",
//     ErrorUserMsg: "UI Error Tips: %v",
// }
//
// Error {
//     ErrorCode: 400,
//     ErrorMsg: "something wrong with: %v",
//     ErrorUserMsg: "UI Error Tips",
// }
func (err Error) WrapPrint(core error, message string, user ...interface{}) error {
	if core == nil {
		return nil
	}
	ret := err
	SetErrPrintfMsg(&ret, core)
	if len(user) > 0 {
		SetErrPrintfUserMsg(&ret, user...)
	}
	return errors.Wrap(ret, message)
}

// Wrap Error like
// Error {
//     ErrorCode: 400,
//     ErrorMsg: "action error",
//     ErrorUserMsg: "UI Error Tips",
// }
func (err Error) Wrap(core error) error {
	if core == nil {
		return nil
	}

	msg := err.ErrorMsg
	err.ErrorMsg = core.Error()
	return errors.Wrap(err, msg)
}

func SetErrPrintfMsg(err *Error, v ...interface{}) {
	err.ErrorMsg = fmt.Sprintf(err.ErrorMsg, v...)
}

func SetErrPrintfUserMsg(err *Error, v ...interface{}) {
	err.ErrorUserMsg = fmt.Sprintf(err.ErrorMsg, v...)
}
