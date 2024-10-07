/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
 **/

package job

import (
	"context"
	"errors"
	"time"
)

var ErrJobTimeout = errors.New("任务超时")

type RetryJetLagFunc func(ctx context.Context, retryCoun int, lastTime time.Duration) time.Duration

func RetryJetLagAlways(ctx context.Context, retryCount int, lastTime time.Duration) time.Duration {
	return DefaultRetryJetLag
}

type RetryFunc func(ctx context.Context, retryCount int, err error) bool

// RetryAlways always retry on error
func RetryAlways(ctx context.Context, retryCount int, err error) bool {
	return true
}

func WithRetry(ctx context.Context, handler func(context.Context) error,
	opts ...RetryOption) error {
	opt := newOption(opts...)

	// 判断是否设置超时
	_, ok := ctx.Deadline()
	if !ok {
		// no deadline so we need create a new one
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opt.timeout)
		defer cancel()
	}

	var (
		herr        error
		retryJetLag time.Duration
		ch          = make(chan error, 1)
	)

	for i := 0; i < opt.retryNums; i++ {
		go func() {
			// 执行任务
			ch <- handler(ctx)
		}()

		select {
		case herr = <-ch:
			if herr == nil {
				// 成功
				return nil
			}

			if !opt.retryFunc(ctx, i, herr) {
				// 不重试
				return herr
			}

			retryJetLag = opt.retryJetLa(ctx, i, retryJetLag)
			time.Sleep(retryJetLag)
		case <-ctx.Done():
			return ErrJobTimeout
		}
	}

	return herr
}
