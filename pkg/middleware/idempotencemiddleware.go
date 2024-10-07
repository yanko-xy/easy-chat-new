/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
 **/

package middleware

import (
	"github.com/yanko-xy/easy-chat/pkg/interceptor"
	"net/http"
)

type IdempotenceMiddle struct {
}

func NewIdempotenceMiddle() *IdempotenceMiddle {
	return &IdempotenceMiddle{}
}

func (m *IdempotenceMiddle) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		request = request.WithContext(interceptor.ContextWithVal(request.Context()))
		next(writer, request)
	}
}
