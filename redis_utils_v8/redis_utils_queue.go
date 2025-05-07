package redisv8

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
)

// 基于redis的队列工具类
type RedisQueueUtils struct {
	key         string
	cli         *redis.Client
	expire      int32 // 超时时间单位秒
	auto_expire bool  // 是否自动更新超期时间  否则手动更新，默认自动更新
	max_size    int64 // 队列最大长度 0表示不限制
}

/*
创建一个Queue操作工具类

  - paramCli
  - paramQueueKey 队列的key
  - paramExpire 超时时间，<=0 时表示没有超时， 单位秒
*/
func CreateRedisQueueUtils(paramCli *redis.Client, paramQueueKey string, paramExpire int32) *RedisQueueUtils {
	return &RedisQueueUtils{
		key:         paramQueueKey,
		cli:         paramCli,
		expire:      paramExpire,
		auto_expire: true,
		max_size:    0,
	}
}
func CreateRedisQueueUtilsMax(paramCli *redis.Client, paramQueueKey string, paramExpire int32, paramMaxSize int64) *RedisQueueUtils {
	return &RedisQueueUtils{
		key:         paramQueueKey,
		cli:         paramCli,
		expire:      paramExpire,
		auto_expire: true,
		max_size:    paramMaxSize,
	}
}

// 计算超时时间
func (m *RedisQueueUtils) calcExpire() time.Duration {
	if m.expire <= 0 {
		return -1
	} else {
		return time.Duration(m.expire) * time.Second
	}
}

// 设置字段更新超时标志
func (m *RedisQueueUtils) SetAutoExpire(paramValue bool) {
	m.auto_expire = paramValue
}

// 设置超时 -1表示设为不过期
func (m *RedisQueueUtils) ExpireSecond(ctx context.Context, paramSeconds int) {
	if paramSeconds < 0 {
		m.cli.Persist(ctx, m.key)
	} else {
		m.cli.Expire(ctx, m.key, time.Duration(paramSeconds)*time.Second)
	}
}

// 取队列的数量
func (m *RedisQueueUtils) Count(ctx context.Context) *redis.IntCmd {
	retCmd := m.cli.LLen(ctx, m.key)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 更新超时时间
func (m *RedisQueueUtils) afterExpire(ctx context.Context, paramErr error) {
	if paramErr == nil && m.auto_expire && m.expire > 0 {
		m.cli.Expire(ctx, m.key, time.Duration(m.expire)*time.Second)
	}
}

// 压入队列
func (m *RedisQueueUtils) Push(ctx context.Context, paramValue ...interface{}) *redis.IntCmd {
	retCmd := m.cli.RPush(ctx, m.key, paramValue...)
	if retCmd.Err() == nil {
		m.afterExpire(ctx, retCmd.Err())
		if m.max_size > 0 {
			lastPos := retCmd.Val()
			if lastPos >= m.max_size {
				m.cli.LTrim(ctx, m.key, lastPos-m.max_size, lastPos)
			}
		}
	}
	return retCmd
}

// 弹出队列
func (m *RedisQueueUtils) Pop(ctx context.Context) *redis.StringCmd {
	retCmd := m.cli.LPop(ctx, m.key)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 弹出多个元素
func (m *RedisQueueUtils) PopCount(ctx context.Context, paramCount int) *redis.StringSliceCmd {
	retCmd := m.cli.LPopCount(ctx, m.key, paramCount)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}
