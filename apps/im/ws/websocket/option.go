/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package websocket

type ServerOptions func(opt *serverOption)

type serverOption struct {
	Authorization
	pattern string
}

func newServerOptions(opts ...ServerOptions) *serverOption {

	o := &serverOption{
		Authorization: defaultAuthorization,
		pattern:       defaultPattrn,
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
