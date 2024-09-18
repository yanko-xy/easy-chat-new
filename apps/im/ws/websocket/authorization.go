/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package websocket

import (
	"fmt"
	"net/http"
	"time"
)

type Authorization interface {
	Auth(w http.ResponseWriter, r *http.Request) bool
	UserId(r *http.Request) string
}

type authorization struct{}

func NewDefaultAuthorization() *authorization {
	return &authorization{}
}

func (*authorization) Auth(w http.ResponseWriter, r *http.Request) bool {
	return true
}

func (*authorization) UserId(r *http.Request) string {
	query := r.URL.Query()
	if query != nil && query["userId"] != nil {
		return fmt.Sprintf("%v", query["userId"])
	}
	return fmt.Sprintf("%v", time.Now().UnixMilli())
}
