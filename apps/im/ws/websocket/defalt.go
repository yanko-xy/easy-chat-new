/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package websocket

import (
	"math"
	"time"
)

var (
	defaultAuthorization = new(authorization)
	defaultPattrn        = "/ws"

	defaultAckTimeout = time.Second * 30

	defaultMaxConnectionIdle = time.Duration(math.MaxInt64)
)
