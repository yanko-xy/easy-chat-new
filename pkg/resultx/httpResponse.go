package resultx

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/yanko-xy/easy-chat/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
	zrpc "github.com/zeromicro/x/errors"
	"google.golang.org/grpc/status"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(data interface{}) *Response {
	return &Response{
		Code: 200,
		Msg:  "",
		Data: data,
	}
}

func Fail(code int, msg string) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

func OkHandler(_ context.Context, v interface{}) any {
	return Success(v)
}

func ErrHandler(name string) func(ctx context.Context, err error) (int, any) {
	return func(ctx context.Context, err error) (int, any) {
		errcode := xerr.SERVER_COMMON_ERROR
		errmsg := xerr.ErrMsg(errcode)

		causeErr := errors.Cause(err)
		if e, ok := causeErr.(*zrpc.CodeMsg); ok {
			errcode = e.Code
			errmsg = e.Msg
		} else {
			if gstaus, ok := status.FromError(causeErr); ok {
				errcode = int(gstaus.Code())
				errmsg = gstaus.Message()
			}
		}

		// 日志记录
		logx.WithContext(ctx).Errorf("【%s】err %v", name, err)

		return http.StatusBadRequest, Fail(errcode, errmsg)
	}
}
