package cluster

import (
	"context"
	"errors"
	pool "github.com/jolestar/go-commons-pool/v2"
	"go-redis/resp/client"
)

type connectionFactory struct {
	Peer string // 对方的ip地址
}

func (cf connectionFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	// 创建客户端根据peer
	c, err := client.MakeClient(cf.Peer)
	if err != nil {
		return nil, err
	}
	// 登录
	c.Start()
	// 池子里面的对象是client
	return pool.NewPooledObject(c), nil
}

func (cf connectionFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	c, ok := object.Object.(*client.Client)
	if !ok {
		return errors.New("type mismatch")
	}
	c.Close()
	return nil
}

func (cf connectionFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	return true
}

func (cf connectionFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

func (cf connectionFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}
