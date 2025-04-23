package dec

import (
	"github.com/shopspring/decimal"
)

// 将数据类型转换为decimal.Decimal类型
// 支持：int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string, decimal.Decimal
// 其他类型返回 decimal.Decimal(0)
func D[T any](value T) decimal.Decimal {
	switch v := any(value).(type) {
	case int:
		return decimal.NewFromInt(int64(v))
	case int8:
		return decimal.NewFromInt(int64(v))
	case int16:
		return decimal.NewFromInt(int64(v))
	case int32:
		return decimal.NewFromInt32(v)
	case int64:
		return decimal.NewFromInt(v)
	case uint:
		return decimal.NewFromUint64(uint64(v))
	case uint8:
		return decimal.NewFromUint64(uint64(v))
	case uint16:
		return decimal.NewFromUint64(uint64(v))
	case uint32:
		return decimal.NewFromUint64(uint64(v))
	case uint64:
		return decimal.NewFromUint64(v)
	case float32:
		return decimal.NewFromFloat32(v)
	case float64:
		return decimal.NewFromFloat(v)
	case string:
		if v == "" {
			return decimal.NewFromInt(0)
		}
		r, err := decimal.NewFromString(v)
		if err != nil {
			return decimal.NewFromInt(0)
		} else {
			return r
		}
	case decimal.Decimal:
		return v
	default:
		return decimal.NewFromInt(0)
	}
}

// 两个decimal.Decimal类型相加
func Add(a any, b any) decimal.Decimal {
	return D(a).Add(D(b))
}

// 两个decimal.Decimal类型相减
func Sub(a any, b any) decimal.Decimal {
	return D(a).Sub(D(b))
}

// 两个decimal.Decimal类型相乘
func Mul(a any, b any) decimal.Decimal {
	return D(a).Mul(D(b))
}

// 两个decimal.Decimal类型相除
func Div(a any, b any) decimal.Decimal {
	return D(a).Div(D(b))
}

// a < b 小于
func Lt(a any, b any) bool {
	return D(a).LessThan(D(b))
}

// a <= b 小于等于
func Lte(a any, b any) bool {
	return D(a).LessThanOrEqual(D(b))
}

// a > b 大于
func Gt(a any, b any) bool {
	return D(a).GreaterThan(D(b))
}

// a >= b 大于等于
func Gte(a any, b any) bool {
	return D(a).GreaterThanOrEqual(D(b))
}

// a == b 等于
func Eq(a any, b any) bool {
	return D(a).Equal(D(b))
}

// 比较两个数的大小
//
//	-1 if x <  y
//	 0 if x == y
//	+1 if x >  y
func Cmp(x any, y any) int {
	return D(x).Cmp(D(y))
}

// 现在的decimal.Decimal类型保留14位小数
func R14(a decimal.Decimal) decimal.Decimal {
	return a.Round(14)
}

// 取0值的数字
func Zero() decimal.Decimal {
	return decimal.NewFromInt(0)
}

// 将当前数字转换为分
func YuanToCent(a decimal.Decimal) int64 {
	return a.Mul(decimal.NewFromInt(100)).Round(0).IntPart()
}

// 将分转换为元
func CentToYuan(a int64) decimal.Decimal {
	return decimal.NewFromInt(a).Div(decimal.NewFromInt(100))
}

// 修正金额为分
// 类似于helper.FormatAmount的效果
func FixMoneyForCent(a decimal.Decimal) float64 {
	value, _ := a.Round(2).Float64()
	return value
}

func FixMoneyForJili(a decimal.Decimal) float64 {
	value, _ := a.RoundDown(12).Float64()
	return value
}

func FixMoneyForJiliFloat(a float64) float64 {
	return FixMoneyForJili(D(a))
}

// 修正百分比为小数
func FixPercent(a decimal.Decimal) float64 {
	value, _ := a.Mul(decimal.NewFromInt(100)).Round(2).Float64()
	return value
}

// 截断指定小数位数的浮点数
func TruncatePlaces(a decimal.Decimal, places int) float64 {
	return a.Truncate(int32(places)).InexactFloat64()
}

// 实现取负的效果
// a = -a
func Neg(a decimal.Decimal) decimal.Decimal {
	return a.Neg()
}

// 金额转字符串
func S(a decimal.Decimal) string {
	return a.String()
}

// 计算百分比
func CalcPercent(a any, b any) string {
	bb := D(b)
	aa := D(a)
	if bb.IsZero() || aa.IsZero() {
		return "0"
	}
	return aa.Div(bb).Mul(decimal.NewFromInt(100)).Round(2).StringFixed(2)
}
