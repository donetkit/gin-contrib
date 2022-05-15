package redis

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/donetkit/gin-contrib/utils/strings"
	"github.com/donetkit/gin-contrib/utils/uuid"
	"github.com/go-redis/redis/v8"
	"strconv"
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
func (r *Client) Set(db int, key string, val interface{}, timeout time.Duration, c ...context.Context) error {
	var err error
	var ctx = r.config.ctx
	if len(c) > 0 {
		ctx = c[0]
	}
	var data []byte
	if data, err = json.Marshal(val); err != nil {
		return err
	}
	err = r.Client[db].Set(ctx, key, data, timeout).Err()
	return err
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
			r.config.logger.Errorf("watch key is modified, retry to release lock. err: %s", err.Error())
		} else {
			r.config.logger.Errorf("err: %s", err.Error())
			return false
		}
	}
}

func CmdString(cmd redis.Cmder) string {
	b := make([]byte, 0, 32)
	b = AppendCmd(b, cmd)
	return strings.String(b)
}

func CmdsString(cmds []redis.Cmder) (string, string) {
	const numCmdLimit = 100
	const numNameLimit = 10

	seen := make(map[string]struct{}, numNameLimit)
	unqNames := make([]string, 0, numNameLimit)

	b := make([]byte, 0, 32*len(cmds))

	for i, cmd := range cmds {
		if i > numCmdLimit {
			break
		}

		if i > 0 {
			b = append(b, '\n')
		}
		b = AppendCmd(b, cmd)

		if len(unqNames) >= numNameLimit {
			continue
		}

		name := cmd.FullName()
		if _, ok := seen[name]; !ok {
			seen[name] = struct{}{}
			unqNames = append(unqNames, name)
		}
	}

	summary := strings.Join(unqNames, " ")
	return summary, strings.String(b)
}

func AppendCmd(b []byte, cmd redis.Cmder) []byte {
	const numArgLimit = 32

	for i, arg := range cmd.Args() {
		if i > numArgLimit {
			break
		}
		if i > 0 {
			b = append(b, ' ')
		}
		b = appendArg(b, arg)
	}

	if err := cmd.Err(); err != nil {
		b = append(b, ": "...)
		b = append(b, err.Error()...)
	}

	return b
}

func appendArg(b []byte, v interface{}) []byte {
	const argLenLimit = 64

	switch v := v.(type) {
	case nil:
		return append(b, "<nil>"...)
	case string:
		if len(v) > argLenLimit {
			v = v[:argLenLimit]
		}
		return appendUTF8String(b, strings.Bytes(v))
	case []byte:
		if len(v) > argLenLimit {
			v = v[:argLenLimit]
		}
		return appendUTF8String(b, v)
	case int:
		return strconv.AppendInt(b, int64(v), 10)
	case int8:
		return strconv.AppendInt(b, int64(v), 10)
	case int16:
		return strconv.AppendInt(b, int64(v), 10)
	case int32:
		return strconv.AppendInt(b, int64(v), 10)
	case int64:
		return strconv.AppendInt(b, v, 10)
	case uint:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint8:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint16:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint32:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint64:
		return strconv.AppendUint(b, v, 10)
	case float32:
		return strconv.AppendFloat(b, float64(v), 'f', -1, 64)
	case float64:
		return strconv.AppendFloat(b, v, 'f', -1, 64)
	case bool:
		if v {
			return append(b, "true"...)
		}
		return append(b, "false"...)
	case time.Time:
		return v.AppendFormat(b, time.RFC3339Nano)
	default:
		return append(b, fmt.Sprint(v)...)
	}
}

func appendUTF8String(dst []byte, src []byte) []byte {
	if isSimple(src) {
		dst = append(dst, src...)
		return dst
	}

	s := len(dst)
	dst = append(dst, make([]byte, hex.EncodedLen(len(src)))...)
	hex.Encode(dst[s:], src)
	return dst
}

func isSimple(b []byte) bool {
	for _, c := range b {
		if !isSimpleByte(c) {
			return false
		}
	}
	return true
}

func isSimpleByte(c byte) bool {
	return c >= 0x21 && c <= 0x7e
}
