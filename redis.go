package rdb

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

// 普通指令
// 现在返回 *CommandBuilder，它实现了 redis.Cmder 接口，同时支持链式调用
type builder func(ctx context.Context, cmd RdCmd, cmdName Command, args map[string]any, includeArgs ...any) *CommandBuilder

// lua脚本
type lua func(ctx context.Context, lua LuaScript, keyInfo map[string]string, valueInfo map[string]any) *redis.Cmd
type Config struct {
	Host        string `json:"host" yaml:"host"`
	Port        string `json:"port" yaml:"port"`
	Password    string `json:"password" yaml:"password"`
	UserName    string `json:"username" yaml:"username"`
	Db          int    `json:"db" yaml:"db"`
	MaxIdle     int    `json:"maxIdle" yaml:"maxIdle"`
	MinIdle     int    `json:"minIdle" yaml:"minIdle"`
	IdleTimeout int    `json:"idleTimeout" yaml:"idleTimeout"`
	PoolSize    int    `json:"poolSize" yaml:"poolSize"`
}

type RedisClient struct {
	lua
	builder
	Config Config
	Client *redis.Client
}

func NewRedisClient(config Config) *RedisClient {
	client := RedisClient{Client: initRedis(config), Config: config}
	client.builder = client.Handler // Handler 现在返回 *CommandBuilder
	client.lua = client.ExecScript
	return &client
}

func initRedis(c Config) *redis.Client {
	slog.Info("redisDb connect", "info", c)
	addr := c.Host + ":" + c.Port
	redisOpt := &redis.Options{
		Addr:         addr,
		Password:     c.Password,
		Username:     c.UserName,
		DB:           c.Db,
		PoolSize:     c.PoolSize,
		MaxIdleConns: c.MaxIdle,
		MinIdleConns: c.MinIdle,
	}
	rdb := redis.NewClient(redisOpt)
	//rdb.AddHook(RKParesHook{})
	cmd := rdb.Ping(context.Background())
	if cmd.Err() != nil {
		panic("redis connect fail, " + cmd.Err().Error())
	}

	return rdb
}

func (rdm RedisClient) RedisClose() {
	err := rdm.Client.Close()
	if err != nil {
		slog.Error("close redisDb", "index", rdm.Config.Db, "error", err.Error())
	} else {
		slog.Info("close redisDb", "index", rdm.Config.Db)
	}
}

func (rdm RedisClient) Handler(ctx context.Context, cmd RdCmd, cmdName Command, args map[string]any, includeArgs ...any) *CommandBuilder {
	// 返回 CommandBuilder，支持链式调用
	// CommandBuilder 实现了 redis.Cmder 接口，可以直接作为 redis.Cmder 使用
	return NewCommandBuilder(&rdm, ctx, cmd, cmdName, args, includeArgs...)
}

func (rdm RedisClient) PipeLine() *RedisPipeline {
	return newPipeline(rdm)
}
