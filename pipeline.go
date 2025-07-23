package rdb

import (
	base "github.com/preceeder/go.base"
	"github.com/redis/go-redis/v9"
)

type RedisPipeline struct {
	lua
	builder
	Client redis.Pipeliner
}

func newPipeline(client RedisClient) *RedisPipeline {
	pip := RedisPipeline{
		Client: client.Client.Pipeline(),
	}
	pip.builder = pip.Handler
	pip.lua = pip.ExecScript
	return &pip
}

func (pip RedisPipeline) Handler(ctx base.BaseContext, cmd RdCmd, cmdName Command, args map[string]any, includeArgs ...any) *redis.Cmd {
	cmdList, key, subCmd := Build(ctx, cmd, cmdName, args, includeArgs...)
	resultCmd := pip.Client.Do(ctx, cmdList...)
	pip.setExp(ctx, key, subCmd)
	return resultCmd
}

func (pip RedisPipeline) setExp(ctx base.BaseContext, key string, subCmd RdSubCmd) {
	if subCmd.Exp != nil {
		exp := subCmd.Exp()
		pip.Client.Expire(ctx, key, exp)
	}
}

// 这一步才是真正的执行命令， 之前的所有步骤都是在往数组中添加命令， 实际没有发送到redis中
func (pip RedisPipeline) Exec(ctx base.BaseContext) ([]redis.Cmder, error) {
	return pip.Client.Exec(ctx)
}
