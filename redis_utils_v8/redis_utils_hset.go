package redisv8

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
)

// 基于Redis的Hash集合工具类
type RedisHSetUtils struct {
	HSetKey     string
	cli         *redis.Client
	expire      int32 // 超时时间 单位秒
	auto_expire bool  // 是否自动更新超期时间  否则手动更新，默认自动更新
}

/*
创建一个HSet操作工具类

  - paramCli
  - paramHSetKey 集合的key
  - paramExpire 超时时间，<=0 时表示没有超时， 单位秒
*/
func CreateRedisHSetUtils(paramCli *redis.Client, paramHSetKey string, paramExpire int32) *RedisHSetUtils {
	return &RedisHSetUtils{
		HSetKey:     paramHSetKey,
		cli:         paramCli,
		expire:      paramExpire,
		auto_expire: true,
	}
}

// 计算超时时间
func (m *RedisHSetUtils) calcExpire() time.Duration {
	if m.expire <= 0 {
		return -1
	} else {
		return time.Duration(m.expire) * time.Second
	}
}

// 设置某个字段
func (m *RedisHSetUtils) Set(ctx context.Context, paramFieldName string, paramValue interface{}) *redis.IntCmd {
	retCmd := m.cli.HSet(ctx, m.HSetKey, paramFieldName, paramValue)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 设置多个字段
func (m *RedisHSetUtils) MultSet(ctx context.Context, paramFieldValues ...interface{}) *redis.BoolCmd {
	retCmd := m.cli.HMSet(ctx, m.HSetKey, paramFieldValues...)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 获取某个字段
func (m *RedisHSetUtils) Get(ctx context.Context, paramFieldName string) *redis.StringCmd {
	retCmd := m.cli.HGet(ctx, m.HSetKey, paramFieldName)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 设置超时 -1表示设为不过期
func (m *RedisHSetUtils) ExpireSecond(ctx context.Context, paramSeconds int) {
	if paramSeconds < 0 {
		m.cli.Persist(ctx, m.HSetKey)
	} else {
		m.cli.Expire(ctx, m.HSetKey, time.Duration(paramSeconds)*time.Second)
	}
}

// 删除某个字段
func (m *RedisHSetUtils) Del(ctx context.Context, paramFieldName string) *redis.IntCmd {
	retCmd := m.cli.HDel(ctx, m.HSetKey, paramFieldName)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 判断某个字段是否存在
func (m *RedisHSetUtils) Exists(ctx context.Context, paramFieldName string) *redis.BoolCmd {
	retCmd := m.cli.HExists(ctx, m.HSetKey, paramFieldName)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 获取所有字段与值的映射
func (m *RedisHSetUtils) GetAll(ctx context.Context) *redis.StringStringMapCmd {
	retCmd := m.cli.HGetAll(ctx, m.HSetKey)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 获取所有字段
func (m *RedisHSetUtils) GetFields(ctx context.Context) *redis.StringSliceCmd {
	retCmd := m.cli.HKeys(ctx, m.HSetKey)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 获取所有值
func (m *RedisHSetUtils) GetValues(ctx context.Context) *redis.StringSliceCmd {
	retCmd := m.cli.HVals(ctx, m.HSetKey)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 获取字段数量
func (m *RedisHSetUtils) Count(ctx context.Context) *redis.IntCmd {
	retCmd := m.cli.HLen(ctx, m.HSetKey)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 增减指定数值
func (m *RedisHSetUtils) Incr(ctx context.Context, paramFieldName string, paramIncrement int64) *redis.IntCmd {
	retCmd := m.cli.HIncrBy(ctx, m.HSetKey, paramFieldName, paramIncrement)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 自减一
func (m *RedisHSetUtils) Dec(ctx context.Context, paramFieldName string) *redis.IntCmd {
	retCmd := m.cli.HIncrBy(ctx, m.HSetKey, paramFieldName, -1)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

// 自增一
func (m *RedisHSetUtils) Inc(ctx context.Context, paramFieldName string) *redis.IntCmd {
	retCmd := m.cli.HIncrBy(ctx, m.HSetKey, paramFieldName, 1)
	m.afterExpire(ctx, retCmd.Err())
	return retCmd
}

func (m *RedisHSetUtils) afterExpire(ctx context.Context, paramErr error) {
	if paramErr == nil && m.auto_expire && m.expire > 0 {
		m.cli.Expire(ctx, m.HSetKey, time.Duration(m.expire)*time.Second)
	}
}
