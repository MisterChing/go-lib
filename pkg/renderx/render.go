package renderx

import (
	"fmt"
	"strings"
)

type AppError struct {
	ErrCode int
	ErrMsg  string
}

func (ae AppError) Code() int {
	return ae.ErrCode
}

func (ae AppError) Msg() string {
	return ae.ErrMsg
}

func (ae AppError) WithExtendMsg(args ...interface{}) AppError {
	extraCount := strings.Count(ae.ErrMsg, "%")
	if extraCount > 0 {
		if len(args) > extraCount {
			args = args[0:extraCount]
		}
		ae.ErrMsg = fmt.Sprintf(ae.ErrMsg, args...)
	}
	return ae
}
