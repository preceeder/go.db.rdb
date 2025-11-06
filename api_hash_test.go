package rdb

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// 测试 Hash 操作的 RdCmd 定义
var HashCmd = RdCmd{
	Key: "hash:{{keyName}}",
	CMD: map[Command]RdSubCmd{
		HSET: {
			Params: "{{field}} {{value}}",
			Exp: func() time.Duration {
				return time.Second * 30
			},
		},
		HGET: {
			Params: "{{field}}",
		},
		HDEL: {
			Params: "{{field}}",
		},
		HGETALL: {
			Params: "",
		},
		HMSET: {
			Params: "",
		},
		HMGET: {
			Params: "{{field}}",
		},
		HSETNX: {
			Params: "{{field}} {{value}}",
		},
		HINCRBY: {
			Params: "{{field}} {{increment}}",
		},
		HINCRBYFLOAT: {
			Params: "{{field}} {{increment}}",
		},
		HKEYS: {
			Params: "",
		},
		HLEN: {
			Params: "",
		},
		HVALS: {
			Params: "",
		},
		HEXISTS: {
			Params: "{{field}}",
		},
	},
}

// TestRedisClient_HSet 测试 HSET 命令
func TestRedisClient_HSet(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	cmd := client.HSet(context.Background(), HashCmd, map[string]any{
		"keyName": "test1",
		"field":   "name",
		"value":   "John",
	})

	if cmd.Err() != nil {
		t.Errorf("HSet failed: %v", cmd.Err())
		return
	}

	fmt.Printf("HSET result: %d\n", cmd.Val())
}

// TestRedisClient_HGet 测试 HGET 命令
func TestRedisClient_HGet(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先设置一个字段
	client.HSet(context.Background(), HashCmd, map[string]any{
		"keyName": "test2",
		"field":   "email",
		"value":   "test@example.com",
	})

	// 获取字段值
	cmd := client.HGet(context.Background(), HashCmd, map[string]any{
		"keyName": "test2",
		"field":   "email",
	})

	if cmd.Err() != nil {
		t.Errorf("HGet failed: %v", cmd.Err())
		return
	}

	fmt.Printf("HGET result: %s\n", cmd.Val())
}

// TestRedisClient_HMSet 测试 HMSET 命令（批量设置）
func TestRedisClient_HMSet(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// HMSET 需要多个 field-value 对
	cmd := client.HMSet(context.Background(), HashCmd, map[string]any{
		"keyName": "test3",
	}, "field1", "value1", "field2", "value2", "field3", "value3")

	if cmd.Err() != nil {
		t.Errorf("HMSet failed: %v", cmd.Err())
		return
	}

	fmt.Printf("HMSET result: %v\n", cmd.Val())
}

// TestRedisClient_HMGet 测试 HMGET 命令（批量获取）
func TestRedisClient_HMGet(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先设置多个字段
	client.HMSet(context.Background(), HashCmd, map[string]any{
		"keyName": "test4",
	}, "name", "Alice", "age", "30", "city", "Beijing")

	// 批量获取
	cmd := client.HMGet(context.Background(), HashCmd, map[string]any{
		"keyName": "test4",
		"field":   "name",
	}, "age", "city", "email") // email 不存在

	if cmd.Err() != nil {
		t.Errorf("HMGet failed: %v", cmd.Err())
		return
	}

	values := cmd.Val()
	fmt.Printf("HMGET result: %v\n", values)
}

// TestRedisClient_HGetAll 测试 HGETALL 命令
func TestRedisClient_HGetAll(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先设置多个字段
	client.HMSet(context.Background(), HashCmd, map[string]any{
		"keyName": "test5",
	}, "name", "Bob", "age", "25", "department", "Engineering")

	// 获取所有字段和值
	cmd := client.HGetAll(context.Background(), HashCmd, map[string]any{
		"keyName": "test5",
	})

	if cmd.Err() != nil {
		t.Errorf("HGetAll failed: %v", cmd.Err())
		return
	}

	result := cmd.Val()
	fmt.Printf("HGETALL result: %v\n", result)
}

// TestRedisClient_HDel 测试 HDEL 命令
func TestRedisClient_HDel(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先设置多个字段
	client.HMSet(context.Background(), HashCmd, map[string]any{
		"keyName": "test6",
	}, "field1", "value1", "field2", "value2", "field3", "value3")

	// 删除一个字段
	cmd := client.HDel(context.Background(), HashCmd, map[string]any{
		"keyName": "test6",
		"field":   "field1",
	})

	if cmd.Err() != nil {
		t.Errorf("HDel failed: %v", cmd.Err())
		return
	}

	fmt.Printf("HDEL deleted count: %d\n", cmd.Val())

	// 删除多个字段（使用 includeArgs）
	cmd2 := client.HDel(context.Background(), HashCmd, map[string]any{
		"keyName": "test6",
		"field":   "field2",
	}, "field3")

	if cmd2.Err() != nil {
		t.Errorf("HDel multiple failed: %v", cmd2.Err())
		return
	}

	fmt.Printf("HDEL multiple deleted count: %d\n", cmd2.Val())
}

// TestRedisClient_HSetNx 测试 HSETNX 命令（仅当字段不存在时设置）
func TestRedisClient_HSetNx(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	keyName := "test7"

	// 第一次设置，应该成功
	cmd1 := client.HSetNx(context.Background(), HashCmd, map[string]any{
		"keyName": keyName,
		"field":   "unique_field",
		"value":   "First value",
	})

	if cmd1.Err() != nil {
		t.Errorf("First HSetNx failed: %v", cmd1.Err())
		return
	}

	fmt.Printf("First HSETNX result: %d (should be 1)\n", cmd1.Val())

	// 第二次设置，应该失败（字段已存在）
	cmd2 := client.HSetNx(context.Background(), HashCmd, map[string]any{
		"keyName": keyName,
		"field":   "unique_field",
		"value":   "Second value",
	})

	if cmd2.Err() != nil {
		t.Errorf("Second HSetNx failed: %v", cmd2.Err())
		return
	}

	fmt.Printf("Second HSETNX result: %d (should be 0)\n", cmd2.Val())
}

// TestRedisClient_HIncrBy 测试 HINCRBY 命令
func TestRedisClient_HIncrBy(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先设置一个数字字段
	client.HSet(context.Background(), HashCmd, map[string]any{
		"keyName": "test8",
		"field":   "count",
		"value":   "10",
	})

	// 增加 5
	cmd := client.HIncrBy(context.Background(), HashCmd, map[string]any{
		"keyName":  "test8",
		"field":   "count",
		"increment": 5,
	})

	if cmd.Err() != nil {
		t.Errorf("HIncrBy failed: %v", cmd.Err())
		return
	}

	fmt.Printf("HINCRBY 5 result: %d\n", cmd.Val())
}

// TestRedisClient_HIncrByFloat 测试 HINCRBYFLOAT 命令
func TestRedisClient_HIncrByFloat(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先设置一个浮点数字段
	client.HSet(context.Background(), HashCmd, map[string]any{
		"keyName": "test9",
		"field":   "score",
		"value":   "10.5",
	})

	// 增加 2.3
	cmd := client.HIncrByFloat(context.Background(), HashCmd, map[string]any{
		"keyName":  "test9",
		"field":   "score",
		"increment": 2.3,
	})

	if cmd.Err() != nil {
		t.Errorf("HIncrByFloat failed: %v", cmd.Err())
		return
	}

	fmt.Printf("HINCRBYFLOAT 2.3 result: %s\n", cmd.Val())
}

// TestRedisClient_HKeys 测试 HKEYS 命令
func TestRedisClient_HKeys(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先设置多个字段
	client.HMSet(context.Background(), HashCmd, map[string]any{
		"keyName": "test10",
	}, "name", "Charlie", "age", "28", "city", "Shanghai")

	// 获取所有字段名
	cmd := client.HKeys(context.Background(), HashCmd, map[string]any{
		"keyName": "test10",
	})

	if cmd.Err() != nil {
		t.Errorf("HKeys failed: %v", cmd.Err())
		return
	}

	keys := cmd.Val()
	fmt.Printf("HKEYS result: %v\n", keys)
}

// TestRedisClient_HLen 测试 HLEN 命令
func TestRedisClient_HLen(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先设置多个字段
	client.HMSet(context.Background(), HashCmd, map[string]any{
		"keyName": "test11",
	}, "field1", "value1", "field2", "value2", "field3", "value3")

	// 获取字段数量
	cmd := client.HLen(context.Background(), HashCmd, map[string]any{
		"keyName": "test11",
	})

	if cmd.Err() != nil {
		t.Errorf("HLen failed: %v", cmd.Err())
		return
	}

	fmt.Printf("HLEN result: %d\n", cmd.Val())
}

// TestRedisClient_HVals 测试 HVALS 命令
func TestRedisClient_HVals(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先设置多个字段
	client.HMSet(context.Background(), HashCmd, map[string]any{
		"keyName": "test12",
	}, "name", "David", "age", "32", "role", "Manager")

	// 获取所有字段值
	cmd := client.HVals(context.Background(), HashCmd, map[string]any{
		"keyName": "test12",
	})

	if cmd.Err() != nil {
		t.Errorf("HVals failed: %v", cmd.Err())
		return
	}

	values := cmd.Val()
	fmt.Printf("HVALS result: %v\n", values)
}

// TestRedisClient_HExists 测试 HEXISTS 命令
func TestRedisClient_HExists(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先设置一个字段
	client.HSet(context.Background(), HashCmd, map[string]any{
		"keyName": "test13",
		"field":   "email",
		"value":   "test@example.com",
	})

	// 检查字段是否存在
	cmd1 := client.HExists(context.Background(), HashCmd, map[string]any{
		"keyName": "test13",
		"field":   "email",
	})

	if cmd1.Err() != nil {
		t.Errorf("HExists failed: %v", cmd1.Err())
		return
	}

	fmt.Printf("HEXISTS email: %d (should be 1)\n", cmd1.Val())

	// 检查不存在的字段
	cmd2 := client.HExists(context.Background(), HashCmd, map[string]any{
		"keyName": "test13",
		"field":   "phone",
	})

	if cmd2.Err() != nil {
		t.Errorf("HExists failed: %v", cmd2.Err())
		return
	}

	fmt.Printf("HEXISTS phone: %d (should be 0)\n", cmd2.Val())
}

// TestRedisClient_Hash_Integration 集成测试：Hash 操作的完整流程
func TestRedisClient_Hash_Integration(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	keyName := "integration_hash"

	// 1. HSET 单个字段
	client.HSet(context.Background(), HashCmd, map[string]any{
		"keyName": keyName,
		"field":   "name",
		"value":   "Integration Test",
	})

	// 2. HMSET 多个字段
	client.HMSet(context.Background(), HashCmd, map[string]any{
		"keyName": keyName,
	}, "age", "25", "city", "Beijing", "score", "100")

	// 3. HGETALL
	getAllCmd := client.HGetAll(context.Background(), HashCmd, map[string]any{
		"keyName": keyName,
	})
	fmt.Printf("1. HGETALL: %v\n", getAllCmd.Val())

	// 4. HKEYS
	keysCmd := client.HKeys(context.Background(), HashCmd, map[string]any{
		"keyName": keyName,
	})
	fmt.Printf("2. HKEYS: %v\n", keysCmd.Val())

	// 5. HLEN
	lenCmd := client.HLen(context.Background(), HashCmd, map[string]any{
		"keyName": keyName,
	})
	fmt.Printf("3. HLEN: %d\n", lenCmd.Val())

	// 6. HINCRBY
	incrCmd := client.HIncrBy(context.Background(), HashCmd, map[string]any{
		"keyName":  keyName,
		"field":   "score",
		"increment": 10,
	})
	fmt.Printf("4. HINCRBY score +10: %d\n", incrCmd.Val())

	// 7. HEXISTS
	existsCmd := client.HExists(context.Background(), HashCmd, map[string]any{
		"keyName": keyName,
		"field":   "name",
	})
	fmt.Printf("5. HEXISTS name: %d\n", existsCmd.Val())
}

