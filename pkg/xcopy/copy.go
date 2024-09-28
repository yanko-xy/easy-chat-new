/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package xcopy

import (
	"database/sql"
	"github.com/jinzhu/copier"
)

// 定义两个结构体
type Source struct {
	Timestamp sql.NullTime
}

type Destination struct {
	Timestamp int64
}

func Copy(toValue interface{}, fromValue interface{}) {
	copier.CopyWithOption(toValue, fromValue, copier.Option{
		Converters: []copier.TypeConverter{
			{
				SrcType: sql.NullTime{},
				DstType: int64(0),
				Fn: func(src interface{}) (interface{}, error) {
					nullTime := src.(sql.NullTime)
					if nullTime.Valid {
						return nullTime.Time.Unix(), nil
					}
					return int64(0), nil
				},
			},
		},
		FieldNameMapping: nil,
	})
}
