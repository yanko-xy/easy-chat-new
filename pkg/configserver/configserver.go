/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package configserver

import (
	"errors"
	"github.com/zeromicro/go-zero/core/conf"
)

var ErrNotSetConfig = errors.New("未设置配置信息")

type OnChange func([]byte) error

type Configserver interface {
	Build() error
	SetOnChange(OnChange)
	FromJsonBytes() ([]byte, error)
}

type configserver struct {
	Configserver
	configFile string
}

func NewConfigserver(configFile string, p Configserver) *configserver {
	return &configserver{
		Configserver: p,
		configFile:   configFile,
	}
}

func (p *configserver) MustLoad(v any, onChange OnChange) error {
	if p.configFile == "" && p.Configserver == nil {
		return ErrNotSetConfig
	}

	if p.Configserver == nil {
		conf.MustLoad(p.configFile, v)
		return nil
	}

	if onChange != nil {
		p.SetOnChange(onChange)
	}

	if err := p.Configserver.Build(); err != nil {
		return err
	}

	data, err := p.Configserver.FromJsonBytes()
	if err != nil {
		return err
	}

	return LoadFromJsonBytes(data, v)
}
func LoadFromJsonBytes(data []byte, v any) error {
	return conf.LoadFromJsonBytes(data, v)
}
