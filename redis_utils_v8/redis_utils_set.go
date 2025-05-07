package redisv8

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
)

// 基于Redis的Set工具类
type RedisSetUtils struct {
	key         string
	cli         *redis.Client
	expire      int32 // 超时时间 单位秒
	auto_expire bool  // 是否自动更新超期时间  否则手动更新，默认自动更新
}

/*
创建一个集合操作工具类

  - paramCli
  - paramKey 集合的key
  - paramExpire 超时时间，<=0 时表示没有超时， 单位秒
  - paramAutoExpire 是否自动更新超时时间
*/
func CreateSetUtils(paramCli *redis.Client, paramKey string, paramExpire int32, paramAutoExpire bool) *RedisSetUtils {
	return &RedisSetUtils{
		key:         paramKey,
		cli:         paramCli,
		expire:      paramExpire,
		auto_expire: paramAutoExpire,
	}
}

// 计算超时时间
func (u *RedisSetUtils) calcExpire() time.Duration {
	if u.expire <= 0 {
		return -1
	} else {
		return time.Duration(u.expire) * time.Second
	}
}

// 取集合的key
func (u *RedisSetUtils) GetKey() string {
	return u.key
}

// 取集合的超时时间
func (u *RedisSetUtils) GetExpire() int32 {
	return u.expire
}

func (u *RedisSetUtils) afterExpire(ctx context.Context, paramErr error) {
	if paramErr == nil && u.auto_expire && u.expire > 0 {
		u.cli.Expire(ctx, u.key, time.Duration(u.expire)*time.Second)
	}
}

// 增加一个元素
func (u *RedisSetUtils) Add(paramCtx context.Context, paramValue ...interface{}) (int64, error) {
	ret := u.cli.SAdd(paramCtx, u.key, paramValue...)
	u.afterExpire(paramCtx, ret.Err())
	return ret.Result()
}

// 删除一个元素
func (u *RedisSetUtils) Del(paramCtx context.Context, paramValue ...interface{}) (int64, error) {
	ret := u.cli.SRem(paramCtx, u.key, paramValue...)
	u.afterExpire(paramCtx, ret.Err())
	return ret.Result()
}

// 判断元素是否存在
func (u *RedisSetUtils) Has(paramCtx context.Context, paramValue interface{}) (bool, error) {
	ret := u.cli.SIsMember(paramCtx, u.key, paramValue)
	return ret.Result()
}

// 获取元素的数量
func (u *RedisSetUtils) Count(paramCtx context.Context) (int64, error) {
	ret := u.cli.SCard(paramCtx, u.key)
	return ret.Result()
}

// 清空集合
func (u *RedisSetUtils) Clean(paramCtx context.Context) error {
	ret := u.cli.Del(paramCtx, u.key)
	return ret.Err()
}

// 获取集合所有元素
func (u *RedisSetUtils) List(paramCtx context.Context) ([]string, error) {
	ret := u.cli.SMembers(paramCtx, u.key)
	return ret.Result()
}

// 设置超时 -1表示设为不过期
func (u *RedisSetUtils) ExpireSecond(ctx context.Context, paramSeconds int) {
	if paramSeconds < 0 {
		u.cli.Persist(ctx, u.key)
	} else {
		u.cli.Expire(ctx, u.key, time.Duration(paramSeconds)*time.Second)
	}
}
