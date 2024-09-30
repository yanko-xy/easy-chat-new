/**
 * @author: yanko/chenyangzhao542@gmail.com
 * @doc:
**/

package websocket

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"net/http"
)

type Discover interface {
	Register(serverAddr string) error               // 注册服务
	BindUser(uid string) error                      // 绑定用户
	RelieveUser(uid string) error                   // 解绑用户
	Transpond(msg interface{}, uid ...string) error // 转发
}

type redisDiscover struct {
	serverAddr   string
	auth         http.Header
	srvKey       string
	boundUserKey string
	redis        *redis.Redis
	cliends      map[string]Client
}

func NewRedisDiscover(auth http.Header, srvKey string, redisCfg redis.RedisConf) *redisDiscover {
	return &redisDiscover{
		auth:         auth,
		srvKey:       srvKey,
		boundUserKey: fmt.Sprintf("%s.%s", srvKey, "boundUserKey"),
		redis:        redis.MustNewRedis(redisCfg),
		cliends:      make(map[string]Client),
	}
}

func (d *redisDiscover) Register(serverAddr string) (err error) {
	d.serverAddr = serverAddr

	// 服务列表： redis存储用set
	go d.redis.Set(d.srvKey, serverAddr)
	return
}

func (d *redisDiscover) BindUser(uid string) (err error) {
	//  用户绑定
	exists, err := d.redis.Hexists(d.boundUserKey, uid)
	if err != nil {
		return err
	}
	if exists {
		// 存在绑定关系
		return nil
	}

	// 绑定
	return d.redis.Hset(d.boundUserKey, uid, d.serverAddr)
}

func (d *redisDiscover) RelieveUser(uid string) (err error) {
	_, err = d.redis.Hdel(d.boundUserKey, uid)
	return
}

func (d *redisDiscover) Transpond(msg interface{}, uids ...string) (err error) {
	for _, uid := range uids {
		srvAddr, err := d.redis.Hget(d.boundUserKey, uid)
		if err != nil {
			return err
		}
		srvClient, ok := d.cliends[srvAddr]
		if !ok {
			srvClient = d.createClient(srvAddr)
		}

		if err := d.send(srvClient, msg, uid); err != nil {
			return err
		}
	}
	return
}

func (d *redisDiscover) send(srvClient Client, msg interface{}, uid string) error {
	return srvClient.Send(Message{
		FrameType:    FrameTranspond,
		TranspondUid: uid,
		Data:         msg,
	})
}

func (d *redisDiscover) createClient(srvAddr string) Client {
	return NewClient(srvAddr, WithClientHeader(d.auth))
}
