package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/donetkit/gin-contrib/utils/uuid"
	"github.com/go-redis/redis/v8"
	"time"
)

func (r *Client) Get(db int, key string, c ...context.Context) interface{} {
	var ctx = r.config.ctx
	if len(c) > 0 {
		ctx = c[0]
	}
	data, err := r.Client[db].Get(ctx, key).Bytes()
	if err != nil {
		return nil
	}
	var reply interface{}
	if err = json.Unmarshal(data, &reply); err != nil {
		return nil
	}
	return reply
}

func (r *Client) GetString(db int, key string, c ...context.Context) (string, error) {
	var ctx = r.config.ctx
	if len(c) > 0 {
		ctx = c[0]
	}
	data, err := r.Client[db].Get(ctx, key).Bytes()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

//Set 设置一个值
func (r *Client) Set(db int, key string, val interface{}, timeout time.Duration, c ...context.Context) (err error) {
	var ctx = r.config.ctx
	if len(c) > 0 {
		ctx = c[0]
	}
	var data []byte
	if data, err = json.Marshal(val); err != nil {
		return
	}
	err = r.Client[db].Set(ctx, key, data, timeout).Err()
	if err != nil {
		panic(err)
	}
	return
}

//IsExist 判断key是否存在
func (r *Client) IsExist(db int, key string, c ...context.Context) bool {
	var ctx = r.config.ctx
	if len(c) > 0 {
		ctx = c[0]
	}
	i := r.Client[db].Exists(ctx, key).Val()
	return i > 0
}

//Delete 删除
func (r *Client) Delete(db int, key string, c ...context.Context) (int64, error) {
	var ctx = r.config.ctx
	if len(c) > 0 {
		ctx = c[0]
	}
	cmd := r.Client[db].Del(ctx, key)
	if cmd.Err() != nil {
		return 0, cmd.Err()
	}
	return cmd.Val(), nil
}

// LPush 左进
func (r *Client) LPush(db int, key string, values interface{}, c ...context.Context) (int64, error) {
	var ctx = r.config.ctx
	if len(c) > 0 {
		ctx = c[0]
	}
	cmd := r.Client[db].LPush(ctx, key, values)
	if cmd.Err() != nil {
		return 0, cmd.Err()
	}
	return cmd.Val(), nil
}

// RPop 右出
func (r *Client) RPop(db int, key string, c ...context.Context) interface{} {
	var ctx = r.config.ctx
	if len(c) > 0 {
		ctx = c[0]
	}
	cmd := r.Client[db].RPop(ctx, key)
	if cmd.Err() != nil {
		return nil
	}
	var reply interface{}
	if err := json.Unmarshal([]byte(cmd.Val()), &reply); err != nil {
		return nil
	}
	return reply
}

func (r *Client) XRead(db int, key string, count int64, c ...context.Context) ([]redis.XStream, error) {
	var ctx = r.config.ctx
	if len(c) > 0 {
		ctx = c[0]
	}
	if count <= 0 {
		count = 10
	}
	msg, err := r.Client[db].XRead(ctx, &redis.XReadArgs{
		Streams: []string{key, "0"},
		Count:   count,
		Block:   500 * time.Millisecond,
	}).Result()
	//msg, err := r.Client[db].XReadStreams(ctx, key, fmt.Sprintf("%d", count)).Result()
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (r *Client) XAdd(db int, key, id string, values []string, c ...context.Context) (string, error) {
	var ctx = r.config.ctx
	if len(c) > 0 {
		ctx = c[0]
	}
	id, err := r.Client[db].XAdd(ctx, &redis.XAddArgs{
		Stream: key,
		ID:     id,
		Values: values,
	}).Result()
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *Client) XDel(db int, key string, id string, c ...context.Context) (int64, error) {
	var ctx = r.config.ctx
	if len(c) > 0 {
		ctx = c[0]
	}
	n, err := r.Client[db].XDel(ctx, key, id).Result()
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (r *Client) GetLock(db int, lockName string, acquireTimeout, lockTimeOut time.Duration, c ...context.Context) (string, error) {
	var ctx = r.config.ctx
	if len(c) > 0 {
		ctx = c[0]
	}
	code := uuid.NewUUID()
	//endTime := util.FwTimer.CalcMillis(time.Now().Add(acquireTimeout))
	endTime := time.Now().Add(acquireTimeout).UnixNano()
	//for util.FwTimer.CalcMillis(time.Now()) <= endTime {
	for time.Now().UnixNano() <= endTime {
		if success, err := r.Client[db].SetNX(ctx, lockName, code, lockTimeOut).Result(); err != nil && err != redis.Nil {
			return "", err
		} else if success {
			return code, nil
		} else if r.Client[db].TTL(ctx, lockName).Val() == -1 {
			r.Client[db].Expire(ctx, lockName, lockTimeOut)
		}
		time.Sleep(time.Millisecond)
	}
	return "", fmt.Errorf("timeout")
}

func (r *Client) ReleaseLock(db int, lockName, code string, c ...context.Context) bool {
	var ctx = r.config.ctx
	if len(c) > 0 {
		ctx = c[0]
	}
	txf := func(tx *redis.Tx) error {
		if v, err := tx.Get(ctx, lockName).Result(); err != nil && err != redis.Nil {
			return err
		} else if v == code {
			_, err := tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
				//count++
				//fmt.Println(count)
				pipe.Del(ctx, lockName)
				return nil
			})
			return err
		}
		return nil
	}

	for {
		if err := r.Client[db].Watch(ctx, txf, lockName); err == nil {
			return true
		} else if err == redis.TxFailedErr {
			fmt.Println("watch key is modified, retry to release lock. err:", err.Error())
		} else {
			fmt.Println("err:", err.Error())
			return false
		}
	}
}
