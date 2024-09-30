/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package websocket

import "time"

type ServerOptions func(opt *serverOption)

type serverOption struct {
	Authorization

	ack        AckType
	ackTimeout time.Duration

	pattern string

	discover Discover

	maxConnectionIdle time.Duration
	concurrency       int
}

func newServerOptions(opts ...ServerOptions) *serverOption {

	o := &serverOption{
		Authorization:     defaultAuthorization,
		ackTimeout:        defaultAckTimeout,
		pattern:           defaultPattrn,
		maxConnectionIdle: defaultMaxConnectionIdle,
		concurrency:       defaultServerHandlerConcurrency,
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

func WithServerAck(ack AckType) ServerOptions {
	return func(opt *serverOption) {
		opt.ack = ack
	}
}

func WithServerAckTimeout(ackTimeout time.Duration) ServerOptions {
	return func(opt *serverOption) {
		opt.ackTimeout = ackTimeout
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

func WithServerWithConcurrency(concurrency int) ServerOptions {
	return func(opt *serverOption) {
		if concurrency > 0 {
			opt.concurrency = concurrency
		}
	}
}

func WithServerDiscover(discover Discover) ServerOptions {
	return func(opt *serverOption) {
		opt.discover = discover
	}
}
