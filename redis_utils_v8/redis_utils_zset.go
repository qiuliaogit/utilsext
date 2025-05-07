package redisv8

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/qiuliaogit/commonutils/commonutils"
)

const (
	MAX_VALUE = "+inf"
	MIN_VALUE = "-inf"
)

// 基于redis的有序集合集合工具类
type RedisZSetUtils struct {
	key         string
	cli         *redis.Client
	expire      int32 // 超时时间单位秒
	auto_expire bool  // 是否自动更新超期时间  否则手动更新，默认自动更新
}

/*
创建一个有序集合操作工具类

  - paramCli
  - paramKey 集合的key
  - paramExpire 超时时间，<=0 时表示没有超时， 单位秒
*/
func CreateRedisZSetUtils(paramCli *redis.Client, paramKey string, paramExpire int32) *RedisZSetUtils {
	return &RedisZSetUtils{
		key:         paramKey,
		cli:         paramCli,
		expire:      paramExpire,
		auto_expire: false,
	}
}

// 计算超时时间
func (m *RedisZSetUtils) calcExpire() time.Duration {
	if m.expire <= 0 {
		return -1
	} else {
		return time.Duration(m.expire) * time.Second
	}
}

func (m *RedisZSetUtils) SetKey(paramKey string) {
	m.key = paramKey
}

// 设置字段更新超时标志
func (m *RedisZSetUtils) SetAutoExpire(paramValue bool) {
	m.auto_expire = paramValue
}

// 设置超时 -1表示设为不过期
func (m *RedisZSetUtils) ExpireSecond(ctx context.Context, paramSeconds int) {
	if paramSeconds < 0 {
		m.cli.Persist(ctx, m.key)
	} else {
		m.cli.Expire(ctx, m.key, time.Duration(paramSeconds)*time.Second)
	}
}

// 获取集合的数量
func (m *RedisZSetUtils) Count(ctx context.Context) *redis.IntCmd {
	return m.cli.ZCard(ctx, m.key)
}

// 获取指定分数范围内成员的数量
func (m *RedisZSetUtils) CountByScore(ctx context.Context, paramMinScore, paramMaxScore float64) *redis.IntCmd {
	return m.cli.ZCount(ctx, m.key, commonutils.Float2Str(paramMinScore), commonutils.Float2Str(paramMaxScore))
}

// 获取指定分数范围内成员的数量（指定最小分数）
func (m *RedisZSetUtils) CountByMinScore(ctx context.Context, paramMinScore float64) *redis.IntCmd {
	return m.cli.ZCount(ctx, m.key, commonutils.Float2Str(paramMinScore), MAX_VALUE)
}

// 获取指定分数范围内成员的数量（指定最大分数）
func (m *RedisZSetUtils) CountByMaxScore(ctx context.Context, paramMaxScore float64) *redis.IntCmd {
	return m.cli.ZCount(ctx, m.key, MIN_VALUE, commonutils.Float2Str(paramMaxScore))
}

// 获取指定分数范围内成员的数量（整数分数）
func (m *RedisZSetUtils) CountByIntScore(ctx context.Context, paramMinScore, paramMaxScore int64) *redis.IntCmd {
	return m.CountByScore(ctx, float64(paramMinScore), float64(paramMaxScore))
}

// 获取指定分数范围内成员的数量（指定最大分数）
func (m *RedisZSetUtils) CountByIntMaxScore(ctx context.Context, paramMaxScore int64) *redis.IntCmd {
	return m.CountByMaxScore(ctx, float64(paramMaxScore))
}

// 清除分数
func (m *RedisZSetUtils) ZeroScore(ctx context.Context) *commonutils.BaseRet {
	r := commonutils.NewBaseRet()

	for range [1]int{} {
		list, err := m.MemberListByMinScore(ctx, 0, false).Result()
		if err != nil && err != redis.Nil {
			r.SetError(commonutils.ERR_FAIL, "获取分数大于0的成员列表失败："+m.key+" err:"+err.Error())
			break
		}
		if len(list) == 0 {
			break
		}

		updateList := make([]*redis.Z, 0, 500)
		for _, member := range list {
			updateList = append(updateList, &redis.Z{Member: member, Score: 0})
			if len(updateList) >= 500 {
				err = m.cli.ZAdd(ctx, m.key, updateList...).Err()
				if err != nil {
					r.SetError(commonutils.ERR_FAIL, "更新分数失败："+m.key+" err:"+err.Error())
					break
				}
				updateList = make([]*redis.Z, 0, 500)
			}
		}
		if r.IsNotOK() {
			break
		}
		if len(updateList) > 0 {
			err := m.cli.ZAdd(ctx, m.key, updateList...).Err()
			if err != nil {
				r.SetError(commonutils.ERR_FAIL, "更新分数失败："+m.key+" err:"+err.Error())
				break
			}
		}
	}
	return r
}

// 获取指定分数范围内成员的数量（指定最小分数 >= paramMinScore）
func (m *RedisZSetUtils) CountByIntMinScore(ctx context.Context, paramMinScore int64) *redis.IntCmd {
	return m.CountByMinScore(ctx, float64(paramMinScore))
}

// 增加成员或设置成员
func (m *RedisZSetUtils) Add(ctx context.Context, paramMembers ...*redis.Z) *redis.IntCmd {
	return m.cli.ZAdd(ctx, m.key, paramMembers...)
}

// 增加一个成员或设置成员
func (m *RedisZSetUtils) AddOne(ctx context.Context, paramMember string, paramScore float64) *redis.IntCmd {
	return m.cli.ZAdd(ctx, m.key, &redis.Z{Member: paramMember, Score: paramScore})
}

// 增加一个成员或设置成员
func (m *RedisZSetUtils) AddOneIntScore(ctx context.Context, paramMember string, paramScore int64) *redis.IntCmd {
	return m.cli.ZAdd(ctx, m.key, &redis.Z{Member: paramMember, Score: float64(paramScore)})
}

// 移除指定member的元素
func (m *RedisZSetUtils) Remove(ctx context.Context, paramMembers ...interface{}) *redis.IntCmd {
	return m.cli.ZRem(ctx, m.key, paramMembers...)
}

// 移除指定排名范围内的元素 排名值从0开始
func (m *RedisZSetUtils) RemoveRangeByRank(ctx context.Context, paramStartRank, paramStopRank int64) *redis.IntCmd {
	return m.cli.ZRemRangeByRank(ctx, m.key, paramStartRank, paramStopRank)
}

// 移除指定分数范围内的元素（浮点分数）
func (m *RedisZSetUtils) RemoveRangeByScore(ctx context.Context, paramMinScore, paramMaxScore float64) *redis.IntCmd {
	return m.cli.ZRemRangeByScore(ctx, m.key, commonutils.Float2Str(paramMinScore), commonutils.Float2Str(paramMaxScore))
}

// 移除指定分数范围内的元素(整数分数)
func (m *RedisZSetUtils) RemoveRangeByIntScore(ctx context.Context, paramMinScore, paramMaxScore int64) *redis.IntCmd {
	return m.cli.ZRemRangeByScore(ctx, m.key, commonutils.I(paramMinScore), commonutils.I(paramMaxScore))
}

// 给指定成员增加分数
func (m *RedisZSetUtils) IncrementScore(ctx context.Context, paramMember string, paramIncrement float64) *redis.FloatCmd {
	return m.cli.ZIncrBy(ctx, m.key, paramIncrement, paramMember)
}

// 给指定成员增加分数
func (m *RedisZSetUtils) IncrementScoreIntScore(ctx context.Context, paramMember string, paramIncrement int64) *redis.FloatCmd {
	return m.cli.ZIncrBy(ctx, m.key, float64(paramIncrement), paramMember)
}

// 获取指定成员的分数
func (m *RedisZSetUtils) GetScore(ctx context.Context, paramMember string) *redis.FloatCmd {
	return m.cli.ZScore(ctx, m.key, paramMember)
}

// 取指定排名范围内的成员列表 排名值从0开始(按分数重小到大， 顺序)
// 是包括【startRank, stopRank】的
func (m *RedisZSetUtils) MemberListByRank(ctx context.Context, paramStartRank, paramStopRank int64) *redis.StringSliceCmd {
	return m.cli.ZRange(ctx, m.key, paramStartRank, paramStopRank)
}

func (m *RedisZSetUtils) MemberListWithScore(ctx context.Context, paramStartRank, paramStopRank int64) *redis.ZSliceCmd {
	return m.cli.ZRangeWithScores(ctx, m.key, paramStartRank, paramStopRank)
}

// 倒叙
func (m *RedisZSetUtils) MemberListRevWithScore(ctx context.Context, limit int64) *redis.ZSliceCmd {
	return m.cli.ZRevRangeByScoreWithScores(ctx, m.key, &redis.ZRangeBy{Offset: 0, Count: limit, Min: "-inf", Max: "+inf"})
	//return m.cli.ZRevRangeWithScores(ctx, m.key, paramStartRank, limit)
}

// 倒叙
func (m *RedisZSetUtils) MemberListRevWithScore2(ctx context.Context, paramStartRank, paramStopRank int64) *redis.ZSliceCmd {
	return m.cli.ZRevRangeWithScores(ctx, m.key, paramStartRank, paramStopRank)
}

// 取指定排名范围内的成员列表 排名值从0开始(按分数重大到小， 逆序)
func (m *RedisZSetUtils) MemberListByScore(ctx context.Context, paramMinScore, paramMaxScore float64) *redis.StringSliceCmd {
	opt := &redis.ZRangeBy{
		Min:    commonutils.Float2Str(paramMinScore),
		Max:    commonutils.Float2Str(paramMaxScore),
		Offset: 0,
		Count:  -1,
	}
	return m.cli.ZRangeByScore(ctx, m.key, opt)
}

/*
取分数 > paramMinScore 或 >= paramMinScore 的成员列表
  - paramMinScore 最小分数
  - paramIncludeMin 是否包含最小分数的成员
*/
func (m *RedisZSetUtils) MemberListByMinScore(ctx context.Context, paramMinScore float64, paramIncludeMin bool) *redis.StringSliceCmd {
	MinValue := commonutils.Float2Str(paramMinScore)
	if !paramIncludeMin {
		MinValue = "(" + MinValue
	}

	opt := &redis.ZRangeBy{
		Min:    MinValue,
		Max:    MAX_VALUE,
		Offset: 0,
		Count:  -1,
	}
	return m.cli.ZRangeByScore(ctx, m.key, opt)
}

/*
取分数 < paramMaxScore 或 <= paramMaxScore 的成员列表
  - paramMaxScore 最大分数
  - paramIncludeMax 是否包含最小分数的成员
*/
func (m *RedisZSetUtils) MemberListByMaxScore(ctx context.Context, paramMaxScore float64, paramIncludeMax bool) *redis.StringSliceCmd {
	MaxValue := commonutils.Float2Str(paramMaxScore)
	if !paramIncludeMax {
		MaxValue = "(" + MaxValue
	}

	opt := &redis.ZRangeBy{
		Min:    MIN_VALUE,
		Max:    MaxValue,
		Offset: 0,
		Count:  -1,
	}
	return m.cli.ZRangeByScore(ctx, m.key, opt)
}

// 取指定排名范围内的成员列表 排名值从0开始（按分数重大到小， 逆序）
// 是包括【startRank, stopRank】的
func (m *RedisZSetUtils) MemberListByRankRevScore(ctx context.Context, paramStartRank, paramStopRank int64) *redis.StringSliceCmd {
	return m.cli.ZRevRange(ctx, m.key, paramStartRank, paramStopRank)
}

// 删除小于最大分数的成员
func (m *RedisZSetUtils) RemoveByMaxScore(ctx context.Context, paramMaxScore float64, paramIncludeMax bool) *redis.IntCmd {
	MaxValue := commonutils.Float2Str(paramMaxScore)
	if !paramIncludeMax {
		MaxValue = "(" + MaxValue
	}
	return m.cli.ZRemRangeByScore(ctx, m.key, MIN_VALUE, MaxValue)
}

/*
删除当前有序集合不在另外一个有序集合中的成员，只保留存在的
- paramOtherSetKey 另一个有序集合的key
*/
func (m *RedisZSetUtils) Intersect(ctx context.Context, paramOtherSetKey string) *redis.IntCmd {
	return m.cli.ZInterStore(ctx, m.key, &redis.ZStore{Keys: []string{m.key, paramOtherSetKey}, Weights: []float64{1, 0}})
}
