package rdb

import (
	"context"
)

// HSET key field value
func (b builder) HSet(ctx context.Context, cmd RdCmd, args map[string]any, includeArgs ...any) *CommandBuilder {
	return b(ctx, cmd, HSET, args, includeArgs...)
}

// HGET key field
func (b builder) HGet(ctx context.Context, cmd RdCmd, args map[string]any, includeArgs ...any) *CommandBuilder {
	return b(ctx, cmd, HGET, args, includeArgs...)
}

// HDEL key field [field2 ...], 删除字段，可以同时删除多个
func (b builder) HDel(ctx context.Context, cmd RdCmd, args map[string]any, includeArgs ...any) *CommandBuilder {
	return b(ctx, cmd, HDEL, args, includeArgs...)
}

// HGETALL key
func (b builder) HGetAll(ctx context.Context, cmd RdCmd, args map[string]any, includeArgs ...any) *CommandBuilder {
	return b(ctx, cmd, HGETALL, args, includeArgs...)
}

// HMSET key field1 value1 field2 value2
func (b builder) HMSet(ctx context.Context, cmd RdCmd, args map[string]any, includeArgs ...any) *CommandBuilder {
	return b(ctx, cmd, HMSET, args, includeArgs...)
}

// HMGET key field1  field2
func (b builder) HMGet(ctx context.Context, cmd RdCmd, args map[string]any, includeArgs ...any) *CommandBuilder {
	return b(ctx, cmd, HMGET, args, includeArgs...)
}

// HSETNX key field value , 设置键下字段的值，存在则不操作返回0，不存在并创建成功则返回1
func (b builder) HSetNx(ctx context.Context, cmd RdCmd, args map[string]any, includeArgs ...any) *CommandBuilder {
	return b(ctx, cmd, HSETNX, args, includeArgs...)
}

// HINCRBY key field1  value   , 指定键指定字段自增指定的整数
func (b builder) HIncrBy(ctx context.Context, cmd RdCmd, args map[string]any, includeArgs ...any) *CommandBuilder {
	return b(ctx, cmd, HINCRBY, args, includeArgs...)
}

// HINCRBYFLOAT key field1  value   , 指定键指定字段自增指定的浮点数
func (b builder) HIncrByFloat(ctx context.Context, cmd RdCmd, args map[string]any, includeArgs ...any) *CommandBuilder {
	return b(ctx, cmd, HINCRBYFLOAT, args, includeArgs...)
}

// HKEYS key  , 获取键下的所有字段列表
func (b builder) HKeys(ctx context.Context, cmd RdCmd, args map[string]any, includeArgs ...any) *CommandBuilder {
	return b(ctx, cmd, HKEYS, args, includeArgs...)
}

// HLEN key  , 获取键下字段的数量
func (b builder) HLen(ctx context.Context, cmd RdCmd, args map[string]any, includeArgs ...any) *CommandBuilder {
	return b(ctx, cmd, HLEN, args, includeArgs...)
}

// HVALS key  , 返回哈希表所有的值
func (b builder) HVals(ctx context.Context, cmd RdCmd, args map[string]any, includeArgs ...any) *CommandBuilder {
	return b(ctx, cmd, HVALS, args, includeArgs...)
}

// HEXISTS key field, 键下是否存在指定的字段
func (b builder) HExists(ctx context.Context, cmd RdCmd, args map[string]any, includeArgs ...any) *CommandBuilder {
	return b(ctx, cmd, HEXISTS, args, includeArgs...)
}
