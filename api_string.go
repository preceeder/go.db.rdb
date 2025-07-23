package rdb

import (
	"github.com/preceeder/go.base"
	"github.com/redis/go-redis/v9"
)

func (b builder) Set(ctx base.BaseContext, cmd RdCmd, args map[string]any, includeArgs ...any) *redis.Cmd {
	return b(ctx, cmd, SET, args, includeArgs...)
}

func (b builder) MSet(ctx base.BaseContext, cmd RdCmd, args map[string]any, includeArgs ...any) *redis.Cmd {
	return b(ctx, cmd, MSET, args, includeArgs...)
}

// SETRANGE key offset value   , 用 value 参数覆写给定 key 所储存的字符串值，从偏移量 offset 开始。
func (b builder) SetRange(ctx base.BaseContext, cmd RdCmd, args map[string]any, includeArgs ...any) *redis.Cmd {
	return b(ctx, cmd, SETRANGE, args, includeArgs...)
}

// 将值 value 关联到 key ，并将 key 的过期时间设为 seconds (以秒为单位)。
func (b builder) SetEx(ctx base.BaseContext, cmd RdCmd, args map[string]any, includeArgs ...any) *redis.Cmd {
	return b(ctx, cmd, SETEX, args, includeArgs...)
}

// 只有在 key 不存在时设置 key 的值。
func (b builder) SetNx(ctx base.BaseContext, cmd RdCmd, args map[string]any, includeArgs ...any) *redis.Cmd {
	return b(ctx, cmd, SETNX, args, includeArgs...)
}

func (b builder) Del(ctx base.BaseContext, cmd RdCmd, args map[string]any, includeArgs ...any) *redis.Cmd {
	return b(ctx, cmd, DEL, args, includeArgs...)
}

func (b builder) Get(ctx base.BaseContext, cmd RdCmd, args map[string]any, includeArgs ...any) *redis.Cmd {
	return b(ctx, cmd, GET, args, includeArgs...)
}

func (b builder) GetSet(ctx base.BaseContext, cmd RdCmd, args map[string]any, includeArgs ...any) *redis.Cmd {
	return b(ctx, cmd, GETSET, args, includeArgs...)
}

func (b builder) GetRange(ctx base.BaseContext, cmd RdCmd, args map[string]any, includeArgs ...any) *redis.Cmd {
	return b(ctx, cmd, GETRANGE, args, includeArgs...)
}

func (b builder) MGet(ctx base.BaseContext, cmd RdCmd, args map[string]any, includeArgs ...any) *redis.Cmd {
	return b(ctx, cmd, MGET, args, includeArgs...)
}

func (b builder) Incr(ctx base.BaseContext, cmd RdCmd, args map[string]any, includeArgs ...any) *redis.Cmd {
	return b(ctx, cmd, INCR, args, includeArgs...)
}

func (b builder) IncrBy(ctx base.BaseContext, cmd RdCmd, args map[string]any, includeArgs ...any) *redis.Cmd {
	return b(ctx, cmd, INCRBY, args, includeArgs...)
}
func (b builder) IncrByFloat(ctx base.BaseContext, cmd RdCmd, args map[string]any, includeArgs ...any) *redis.Cmd {
	return b(ctx, cmd, INCRBYFLOAT, args, includeArgs...)
}

func (b builder) DecrBy(ctx base.BaseContext, cmd RdCmd, args map[string]any, includeArgs ...any) *redis.Cmd {
	return b(ctx, cmd, DECRBY, args, includeArgs...)
}

func (b builder) Decr(ctx base.BaseContext, cmd RdCmd, args map[string]any, includeArgs ...any) *redis.Cmd {
	return b(ctx, cmd, DECR, args, includeArgs...)
}

func (b builder) StringAppend(ctx base.BaseContext, cmd RdCmd, args map[string]any, includeArgs ...any) *redis.Cmd {
	return b(ctx, cmd, APPEND, args, includeArgs...)
}
