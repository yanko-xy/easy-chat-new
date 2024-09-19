/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package websocket

import "net/http"

type DailOptins func(option *dailOption)

type dailOption struct {
	pattern string
	header  http.Header
}

func newDailOption(opts ...DailOptins) *dailOption {
	o := &dailOption{
		pattern: defaultPattrn,
		header:  nil,
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

func WithClientPattern(pattern string) DailOptins {
	return func(opt *dailOption) {
		opt.pattern = pattern
	}
}

func WithClientHeader(header http.Header) DailOptins {
	return func(opt *dailOption) {
		opt.header = header
	}
}
