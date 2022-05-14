package err

import "fmt"

const (
	None int = iota
	SystemError
)

var errs map[int]string

type SxgErr struct {
	code int
	msg  string
}

func (se *SxgErr) Error() string {
	if se.code == None {
		return se.msg
	}
	return fmt.Sprintf("code: [%d], msg: %s", se.code, se.msg)
}

func NewErr(code int) *SxgErr {
	msg, ok := errs[code]
	if !ok {
		panic("找不到指定代码的异常描述")
	}
	return &SxgErr{code, msg}
}

func NewErrMsg(msg string) *SxgErr {
	return &SxgErr{None, msg}
}
