package redisutils

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
)

// 基于redis的列表工具类
type RedisListUtils struct {
	key         string
	cli         *redis.Client
	expire      int32 // 超时时间单位秒
	auto_expire bool  // 是否自动更新超期时间  否则手动更新，默认自动更新
}

// 计算超时时间
func (m *RedisListUtils) calcExpire() time.Duration {
	if m.expire <= 0 {
		return -1
	} else {
		return time.Duration(m.expire) * time.Second
	}
}

// 设置字段更新超时标志
func (m *RedisListUtils) SetAutoExpire(paramValue bool) {
	m.auto_expire = paramValue
}

// 设置超时 -1表示设为不过期
func (m *RedisListUtils) ExpireSecond(ctx context.Context, paramSeconds int) {
	if paramSeconds < 0 {
		m.cli.Persist(ctx, m.key)
	} else {
		m.cli.Expire(ctx, m.key, time.Duration(paramSeconds)*time.Second)
	}
}

// 取队列的数量
func (m *RedisListUtils) Count(ctx context.Context) *redis.IntCmd {
	retCmd := m.cli.LLen(ctx, m.key)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 更新超时时间
func (m *RedisListUtils) afterExpire(ctx context.Context, paramErr error) {
	if paramErr == nil && m.auto_expire && m.expire > 0 {
		m.cli.Expire(ctx, m.key, time.Duration(m.expire)*time.Second)
	}
}

// 在列表中添加一个或多个值到列表尾部
func (m *RedisListUtils) RPush(ctx context.Context, paramValue ...interface{}) *redis.IntCmd {
	retCmd := m.cli.RPush(ctx, m.key, paramValue...)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 移除列表的最后一个元素，返回值为移除的元素。
func (m *RedisListUtils) RPop(ctx context.Context) *redis.StringCmd {
	retCmd := m.cli.RPop(ctx, m.key)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 移出并获取列表的第一个元素
func (m *RedisListUtils) LPop(ctx context.Context) *redis.StringCmd {
	retCmd := m.cli.LPop(ctx, m.key)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 弹出多个元素
func (m *RedisListUtils) LPopCount(ctx context.Context, paramCount int) *redis.StringSliceCmd {
	retCmd := m.cli.LPopCount(ctx, m.key, paramCount)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 将一个或多个值插入到列表头部
func (m *RedisListUtils) LPush(ctx context.Context, paramValue ...interface{}) *redis.IntCmd {
	retCmd := m.cli.LPush(ctx, m.key, paramValue...)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 通过索引获取列表中的元素
func (m *RedisListUtils) Get(ctx context.Context, paramIndex int64) *redis.StringCmd {
	retCmd := m.cli.LIndex(ctx, m.key, paramIndex)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 通过索引设置列表元素的值
func (m *RedisListUtils) Set(ctx context.Context, paramIndex int64, paramValue interface{}) *redis.StatusCmd {
	retCmd := m.cli.LSet(ctx, m.key, paramIndex, paramValue)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 移除表中所有与 value 相等的值
func (m *RedisListUtils) Del(ctx context.Context, paramValue interface{}) *redis.IntCmd {
	retCmd := m.cli.LRem(ctx, m.key, 0, paramValue)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 从列表尾部开始删除count个与value相同的值
func (m *RedisListUtils) DelFromTail(ctx context.Context, paramValue interface{}, paramCount uint) *redis.IntCmd {
	retCmd := m.cli.LRem(ctx, m.key, int64(paramCount), paramValue)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 从列表头部开始删除count个与value相同的值
func (m *RedisListUtils) DelFromHead(ctx context.Context, paramValue interface{}, paramCount uint) *redis.IntCmd {
	retCmd := m.cli.LRem(ctx, m.key, -int64(paramCount), paramValue)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 获取列表指定范围内的元素
func (m *RedisListUtils) Range(ctx context.Context, paramStart int64, paramEnd int64) *redis.StringSliceCmd {
	retCmd := m.cli.LRange(ctx, m.key, paramStart, paramEnd)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

/*
创建一个List操作工具类

  - paramCli
  - paramListKey 列表的key
  - paramExpire 超时时间，<=0 时表示没有超时， 单位秒
*/
func CreateRedisListUtils(paramCli *redis.Client, paramListKey string, paramExpire int32) *RedisListUtils {
	return &RedisListUtils{
		key:         paramListKey,
		cli:         paramCli,
		expire:      paramExpire,
		auto_expire: true,
	}
}
