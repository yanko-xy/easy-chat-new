/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package configserver

import (
	"context"
	"encoding/json"
	sailClient "github.com/HYY-yu/sail-client"
	"github.com/zeromicro/go-zero/core/logx"
)

type Config struct {
	ETCDEndpoints  string `toml:"etcd_endpoints"` // 逗号分隔的ETCD地址，0.0.0.0:2379,0.0.0.0:12379,0.0.0.0:22379
	ProjectKey     string `toml:"project_key"`
	Namespace      string `toml:"namespace"`
	Configs        string `toml:"configs"`
	ConfigFilePath string `toml:"config_file_path"` // 本地配置文件存放路径，空代表不存储本都配置文件
	LogLevel       string `toml:"log_level"`        // 日志级别(DEBUG\INFO\WARN\ERROR)，默认 WARN
}

type Sail struct {
	logx.Logger
	*sailClient.Sail
	sailClient.OnConfigChange
	cfg *Config
}

func NewSail(cfg *Config) *Sail {
	return &Sail{cfg: cfg, Logger: logx.WithContext(context.Background())}
}

func (s *Sail) Build() error {
	var opts []sailClient.Option

	if s.OnConfigChange != nil {
		opts = append(opts, sailClient.WithOnConfigChange(s.OnConfigChange))
	}

	s.Sail = sailClient.New(&sailClient.MetaConfig{
		ETCDEndpoints:  s.cfg.ETCDEndpoints,
		ProjectKey:     s.cfg.ProjectKey,
		Namespace:      s.cfg.Namespace,
		Configs:        s.cfg.Configs,
		ConfigFilePath: s.cfg.ConfigFilePath,
		LogLevel:       s.cfg.LogLevel,
	}, opts...)

	return s.Sail.Err()
}

func (s *Sail) FromJsonBytes() ([]byte, error) {
	if err := s.Pull(); err != nil {
		return nil, err
	}

	return s.fromJsonBytes(s.Sail)
}

func (s *Sail) fromJsonBytes(sailClient *sailClient.Sail) ([]byte, error) {
	v, err := sailClient.MergeVipers()
	if err != nil {
		return nil, err
	}
	data := v.AllSettings()
	return json.Marshal(data)
}
func (s *Sail) SetOnChange(f OnChange) {
	s.OnConfigChange = func(configFileKey string, sailClient *sailClient.Sail) {
		data, err := s.fromJsonBytes(sailClient)
		if err != nil {
			s.Errorf("config on change fromJsonBytes err %v", err)
			return
		}

		if err = f(data); err != nil {
			s.Errorf("config on change exec f OnChange err %v, data %v", err, string(data))
		}
	}
}
