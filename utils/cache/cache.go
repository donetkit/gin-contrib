package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type ICache interface {
	Get(db int, key string, c ...context.Context) interface{}
	GetString(db int, key string, c ...context.Context) (string, error)
	Set(db int, key string, val interface{}, timeout time.Duration, c ...context.Context) error
	IsExist(db int, key string, c ...context.Context) bool
	Delete(db int, key string, c ...context.Context) (int64, error)
	LPush(db int, key string, values interface{}, c ...context.Context) (int64, error)
	RPop(db int, key string, c ...context.Context) interface{}
	XRead(db int, key string, count int64, c ...context.Context) ([]redis.XStream, error)
	XAdd(db int, key, id string, values []string, c ...context.Context) (string, error)
	XDel(db int, key string, id string, c ...context.Context) (int64, error)
	GetLock(db int, lockName string, acquireTimeout, lockTimeOut time.Duration, c ...context.Context) (string, error)
	ReleaseLock(db int, lockName, code string, c ...context.Context) bool
}
