/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
 **/

package job

import "time"

type (
	RetryOption func(opts *retryOption)

	retryOption struct {
		timeout    time.Duration
		retryNums  int
		retryFunc  RetryFunc
		retryJetLa RetryJetLagFunc
	}
)

func newOption(opts ...RetryOption) *retryOption {
	opt := &retryOption{
		timeout:    DefaultRetryTimeout,
		retryNums:  DefaultRetryNums,
		retryFunc:  RetryAlways,
		retryJetLa: RetryJetLagAlways,
	}

	for _, option := range opts {
		option(opt)
	}

	return opt
}

func WithRetryTimeout(timeout time.Duration) RetryOption {
	return func(opts *retryOption) {
		if timeout > 0 {
			opts.timeout = timeout
		}
	}
}

func WithRetryNums(nums int) RetryOption {
	return func(opts *retryOption) {
		opts.retryNums = 1
		if nums > 1 {
			opts.retryNums = nums
		}
	}
}

func WithRetryFunc(retryFunc RetryFunc) RetryOption {
	return func(opts *retryOption) {
		if retryFunc != nil {
			opts.retryFunc = retryFunc
		}
	}
}

func WithRetryJetLagFunc(retryJetLagFunc RetryJetLagFunc) RetryOption {
	return func(opts *retryOption) {
		if retryJetLagFunc != nil {
			opts.retryJetLa = retryJetLagFunc
		}
	}
}
