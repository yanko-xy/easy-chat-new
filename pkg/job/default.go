/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
 **/

package job

import "time"

const (
	DefaultRetryJetLag  = 1 * time.Second
	DefaultRetryTimeout = 2 * time.Second
	DefaultRetryNums    = 5
)
