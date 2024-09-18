/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package websocket

import "time"

type ServerOptions func(opt *serverOption)

type serverOption struct {
	Authorization
	pattern           string
	maxConnectionIdle time.Duration
}

func newServerOptions(opts ...ServerOptions) *serverOption {

	o := &serverOption{
		Authorization:     defaultAuthorization,
		pattern:           defaultPattrn,
		maxConnectionIdle: defaultMaxConnectionIdle,
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

func WithServerAuthorization(auth Authorization) ServerOptions {
	return func(opt *serverOption) {
		opt.Authorization = auth
	}
}

func WithServerPattern(pattern string) ServerOptions {
	return func(opt *serverOption) {
		opt.pattern = pattern
	}
}

func WithServerMaxConnectionIdle(maxConnectionIdle time.Duration) ServerOptions {
	return func(opt *serverOption) {
		if maxConnectionIdle > 0 {
			opt.maxConnectionIdle = maxConnectionIdle
		}
	}
}
