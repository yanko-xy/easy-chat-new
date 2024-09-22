package xerr

import "github.com/zeromicro/x/errors"

func New(code int, msg string) error {
	return errors.New(code, msg)
}

func NewMsgErr(msg string) error {
	return errors.New(SERVER_COMMON_ERROR, msg)
}

func NewCodeErr(code int) error {
	return errors.New(code, ErrMsg(code))
}

func NewInternalErr() error {
	return errors.New(SERVER_COMMON_ERROR, ErrMsg(SERVER_COMMON_ERROR))
}

func NewIllegalOperationErr() error {
	return errors.New(ILLEGAL_OPERATION_ERROR, ErrMsg(ILLEGAL_OPERATION_ERROR))
}
func NewDBErr() error {
	return errors.New(DB_ERROR, ErrMsg(DB_ERROR))
}

func NewParameterErr(msg string) error {
	return errors.New(PARAMETER_ERROR, msg)
}

func NewReqParamErr() error {
	return errors.New(REQUEST_PARAM_ERROR, ErrMsg(REQUEST_PARAM_ERROR))
}
