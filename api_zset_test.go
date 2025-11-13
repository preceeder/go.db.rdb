package rdb

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

func requireInt64(t *testing.T, cmd *redis.Cmd) int64 {
	t.Helper()
	if err := cmd.Err(); err != nil {
		t.Fatalf("expected int64 result, got error: %v", err)
	}
	val, err := cmd.Result()
	if err != nil {
		t.Fatalf("expected int64 result, got error: %v", err)
	}
	switch v := val.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	default:
		t.Fatalf("expected int64, got %T", val)
	}
	return 0
}

func requireFloat64(t *testing.T, cmd *redis.Cmd) float64 {
	t.Helper()
	if err := cmd.Err(); err != nil {
		t.Fatalf("expected float64 result, got error: %v", err)
	}
	val, err := cmd.Result()
	if err != nil {
		t.Fatalf("expected float64 result, got error: %v", err)
	}
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	default:
		t.Fatalf("expected float64, got %T", val)
	}
	return 0
}

func requireString(t *testing.T, cmd *redis.Cmd) string {
	t.Helper()
	if err := cmd.Err(); err != nil {
		t.Fatalf("expected string result, got error: %v", err)
	}
	val, err := cmd.Result()
	if err != nil {
		t.Fatalf("expected string result, got error: %v", err)
	}
	return fmt.Sprint(val)
}

func requireStringSlice(t *testing.T, cmd *redis.Cmd) []string {
	t.Helper()
	if err := cmd.Err(); err != nil {
		t.Fatalf("expected string slice result, got error: %v", err)
	}
	val, err := cmd.Result()
	if err != nil {
		t.Fatalf("expected string slice result, got error: %v", err)
	}
	switch v := val.(type) {
	case []string:
		return v
	case []interface{}:
		result := make([]string, len(v))
		for i, item := range v {
			result[i] = fmt.Sprint(item)
		}
		return result
	default:
		t.Fatalf("expected []string, got %T", val)
	}
	return nil
}

var ade = RdCmd{
	Key: "haha",
	CMD: map[Command]RdSubCmd{
		ZADD: {
			Params: "{{score1}} {{key1}} {{score2}} {{key2}} {{score3}} {{key3}}",
			Exp: func() time.Duration {
				return time.Second * 30
			},
		},
		ZRANGE: {
			Params: "{{start}} {{stop}}", //  WITHSCORES
		},
	},
}

func TestRedisClient_ZRange(t *testing.T) {
	client := InitRedis()
	cmd := client.ZAdd(context.Background(), ade, map[string]any{
		"score1": 3, "key1": 2, "score2": 5, "key2": 4, "score3": 5, "key3": 6,
	})
	if cmd.Err() != nil {
		fmt.Println(cmd.Err())
		return
	}

	fmt.Println(cmd.Val())
	cmd = client.ZRange(context.Background(), ade, map[string]any{"start": 0, "stop": -1})
	fmt.Println()
}

// 完整的 ZSet 测试命令定义
var ZSetCmd = RdCmd{
	Key: "zset:{{keyName}}",
	CMD: map[Command]RdSubCmd{
		ZADD: {
			Params: "{{score1}} {{member1}} {{score2}} {{member2}} {{score3}} {{member3}}",
			Exp: func() time.Duration {
				return time.Second * 30
			},
		},
		ZCARD: {
			Params: "",
		},
		ZCOUNT: {
			Params: "{{min}} {{max}}",
		},
		ZINCRBY: {
			Params: "{{increment}} {{member}}",
		},
		ZRANGE: {
			Params: "{{start}} {{stop}}",
		},
		ZREVRANGE: {
			Params: "{{start}} {{stop}}",
		},
		ZRANGEBYSCORE: {
			Params: "{{min}} {{max}}",
		},
		ZREVRANGEBYSCORE: {
			Params: "{{max}} {{min}}",
		},
		ZRANK: {
			Params: "{{member}}",
		},
		ZREVRANK: {
			Params: "{{member}}",
		},
		ZREM: {
			Params: "{{member}}",
		},
		ZREMRANGEBYRANK: {
			Params: "{{start}} {{stop}}",
		},
		ZREMRANGEBYSCORE: {
			Params: "{{min}} {{max}}",
		},
		ZSCORE: {
			Params: "{{member}}",
		},
		ZLEXCOUNT: {
			Params: "{{min}} {{max}}",
		},
		ZREMRANGEBYLEX: {
			Params: "{{min}} {{max}}",
		},
		ZINTERSTORE: {
			Params:   "",
			NoUseKey: true,
		},
		ZINTER: {
			Params:   "",
			NoUseKey: true,
		},
		ZUNIONSTORE: {
			Params:   "",
			NoUseKey: true,
		},
		ZUNION: {
			Params:   "",
			NoUseKey: true,
		},
	},
}

// TestRedisClient_ZAdd 测试 ZADD 命令
func TestRedisClient_ZAdd(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	cmd := client.ZAdd(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test1",
		"score1":  1.0,
		"member1": "member1",
		"score2":  2.0,
		"member2": "member2",
		"score3":  3.0,
		"member3": "member3",
	})

	if cmd.Err() != nil {
		t.Errorf("ZAdd failed: %v", cmd.Err())
		return
	}

	fmt.Printf("ZADD result: %v members added\n", cmd.Val())
}

// TestRedisClient_ZCard 测试 ZCARD 命令
func TestRedisClient_ZCard(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先添加成员
	client.ZAdd(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test2",
		"score1":  1.0,
		"member1": "a",
		"score2":  2.0,
		"member2": "b",
		"score3":  3.0,
		"member3": "c",
	}).String()

	// 获取成员数量
	cmd := client.ZCard(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test2",
	})

	if cmd.Err() != nil {
		t.Errorf("ZCard failed: %v", cmd.Err())
		return
	}

	fmt.Printf("ZCARD result: %d\n", cmd.Val())
}

// TestRedisClient_ZCount 测试 ZCOUNT 命令
func TestRedisClient_ZCount(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先添加成员
	client.ZAdd(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test3",
		"score1":  1.0,
		"member1": "a",
		"score2":  2.0,
		"member2": "b",
		"score3":  5.0,
		"member3": "e",
	})

	// 统计分数在 [1, 3] 之间的成员数
	cmd := client.ZCount(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test3",
		"min":     1,
		"max":     3,
	})

	if cmd.Err() != nil {
		t.Errorf("ZCount failed: %v", cmd.Err())
		return
	}

	fmt.Printf("ZCOUNT [1,3] result: %d\n", cmd.Val())
}

// TestRedisClient_ZRange 测试 ZRANGE 命令
func TestRedisClient_ZRange_New(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先添加成员
	client.ZAdd(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test4",
		"score1":  1.0,
		"member1": "a",
		"score2":  2.0,
		"member2": "b",
		"score3":  3.0,
		"member3": "c",
	})

	// 获取所有成员（按分数升序）
	cmd := client.ZRange(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test4",
		"start":   0,
		"stop":    -1,
	})

	if cmd.Err() != nil {
		t.Errorf("ZRange failed: %v", cmd.Err())
		return
	}

	values := cmd.Val()
	fmt.Printf("ZRANGE result: %v\n", values)
}

// TestRedisClient_ZRevRange 测试 ZREVRANGE 命令
func TestRedisClient_ZRevRange(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先添加成员
	client.ZAdd(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test5",
		"score1":  1.0,
		"member1": "a",
		"score2":  2.0,
		"member2": "b",
		"score3":  3.0,
		"member3": "c",
	})

	// 获取所有成员（按分数降序）
	cmd := client.ZRevRange(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test5",
		"start":   0,
		"stop":    -1,
	})

	if cmd.Err() != nil {
		t.Errorf("ZRevRange failed: %v", cmd.Err())
		return
	}

	values := cmd.Val()
	fmt.Printf("ZREVRANGE result: %v\n", values)
}

// TestRedisClient_ZRangeByScore 测试 ZRANGEBYSCORE 命令
func TestRedisClient_ZRangeByScore(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先添加成员
	client.ZAdd(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test6",
		"score1":  1.0,
		"member1": "a",
		"score2":  2.0,
		"member2": "b",
		"score3":  5.0,
		"member3": "e",
	})

	// 获取分数在 [1, 3] 之间的成员
	cmd := client.ZRangeByScore(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test6",
		"min":     1,
		"max":     3,
	})

	if cmd.Err() != nil {
		t.Errorf("ZRangeByScore failed: %v", cmd.Err())
		return
	}

	values := cmd.Val()
	fmt.Printf("ZRANGEBYSCORE [1,3] result: %v\n", values)
}

// TestRedisClient_ZRevRangeByScore 测试 ZREVRANGEBYSCORE 命令
func TestRedisClient_ZRevRangeByScore(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先添加成员
	client.ZAdd(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test7",
		"score1":  1.0,
		"member1": "a",
		"score2":  2.0,
		"member2": "b",
		"score3":  5.0,
		"member3": "e",
	})

	// 获取分数在 [2, 5] 之间的成员（按分数降序）
	cmd := client.ZRevRangeByScore(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test7",
		"max":     5,
		"min":     2,
	})

	if cmd.Err() != nil {
		t.Errorf("ZRevRangeByScore failed: %v", cmd.Err())
		return
	}

	values := cmd.Val()
	fmt.Printf("ZREVRANGEBYSCORE [2,5] result: %v\n", values)
}

// TestRedisClient_ZRank 测试 ZRANK 命令
func TestRedisClient_ZRank(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先添加成员
	client.ZAdd(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test8",
		"score1":  10.0,
		"member1": "x",
		"score2":  20.0,
		"member2": "y",
		"score3":  30.0,
		"member3": "z",
	})

	// 获取成员排名（升序）
	cmd := client.ZRank(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test8",
		"member":  "y",
	})

	if cmd.Err() != nil {
		t.Errorf("ZRank failed: %v", cmd.Err())
		return
	}

	fmt.Printf("ZRANK y: %d\n", cmd.Val())
}

// TestRedisClient_ZRevRank 测试 ZREVRANK 命令
func TestRedisClient_ZRevRank(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先添加成员
	client.ZAdd(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test9",
		"score1":  10.0,
		"member1": "x",
		"score2":  20.0,
		"member2": "y",
		"score3":  30.0,
		"member3": "z",
	})

	// 获取成员排名（降序）
	cmd := client.ZRevRank(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test9",
		"member":  "y",
	})

	if cmd.Err() != nil {
		t.Errorf("ZRevRank failed: %v", cmd.Err())
		return
	}

	fmt.Printf("ZREVRANK y: %d\n", cmd.Val())
}

// TestRedisClient_ZScore 测试 ZSCORE 命令
func TestRedisClient_ZScore(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先添加成员
	client.ZAdd(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test10",
		"score1":  25.0,
		"member1": "member1",
		"score2":  50.0,
		"member2": "member2",
	})

	// 获取成员分数
	cmd := client.ZScore(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test10",
		"member":  "member2",
	})

	if cmd.Err() != nil {
		t.Errorf("ZScore failed: %v", cmd.Err())
		return
	}

	fmt.Printf("ZSCORE member2: %s\n", cmd.Val())
}

// TestRedisClient_ZIncrBy 测试 ZINCRBY 命令
func TestRedisClient_ZIncrBy(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先添加成员
	client.ZAdd(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test11",
		"score1":  10.0,
		"member1": "member1",
	})

	// 增加分数
	cmd := client.ZIncrBy(context.Background(), ZSetCmd, map[string]any{
		"keyName":   "test11",
		"increment": 5.0,
		"member":    "member1",
	})

	if cmd.Err() != nil {
		t.Errorf("ZIncrBy failed: %v", cmd.Err())
		return
	}

	fmt.Printf("ZINCRBY member1 +5: %s\n", cmd.Val())
}

// TestRedisClient_ZRem 测试 ZREM 命令
func TestRedisClient_ZRem(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先添加成员
	client.ZAdd(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test12",
		"score1":  1.0,
		"member1": "a",
		"score2":  2.0,
		"member2": "b",
		"score3":  3.0,
		"member3": "c",
	})

	// 删除成员
	cmd := client.ZRem(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test12",
		"member":  "b",
	}, "c") // 可以删除多个

	if cmd.Err() != nil {
		t.Errorf("ZRem failed: %v", cmd.Err())
		return
	}

	fmt.Printf("ZREM result: %d members removed\n", cmd.Val())
}

// TestRedisClient_ZRemRangeByRank 测试 ZREMRANGEBYRANK 命令
func TestRedisClient_ZRemRangeByRank(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先添加成员
	client.ZAdd(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test13",
		"score1":  1.0,
		"member1": "a",
		"score2":  2.0,
		"member2": "b",
		"score3":  3.0,
		"member3": "c",
	})

	// 删除排名在 [0, 1] 的成员
	cmd := client.ZRemRangeByRank(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test13",
		"start":   0,
		"stop":    1,
	})

	if cmd.Err() != nil {
		t.Errorf("ZRemRangeByRank failed: %v", cmd.Err())
		return
	}

	fmt.Printf("ZREMRANGEBYRANK [0,1] result: %d members removed\n", cmd.Val())
}

// TestRedisClient_ZRemRangeByScore 测试 ZREMRANGEBYSCORE 命令
func TestRedisClient_ZRemRangeByScore(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先添加成员
	client.ZAdd(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test14",
		"score1":  1.0,
		"member1": "m1",
		"score2":  5.0,
		"member2": "m2",
		"score3":  10.0,
		"member3": "m3",
	})

	// 删除分数在 [2, 9] 之间的成员
	cmd := client.ZRemRangeByScore(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test14",
		"min":     2,
		"max":     9,
	})

	if cmd.Err() != nil {
		t.Errorf("ZRemRangeByScore failed: %v", cmd.Err())
		return
	}

	fmt.Printf("ZREMRANGEBYSCORE [2,9] result: %d members removed\n", cmd.Val())
}

// TestRedisClient_ZLexCount 测试 ZLEXCOUNT 命令
func TestRedisClient_ZLexCount(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先添加成员（分数相同，按字典序排列）
	client.ZAdd(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test15",
		"score1":  0,
		"member1": "a",
		"score2":  0,
		"member2": "b",
		"score3":  0,
		"member3": "c",
	})

	// 统计字典序范围
	cmd := client.ZLexCount(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test15",
		"min":     "-",
		"max":     "+",
	})

	if cmd.Err() != nil {
		t.Errorf("ZLexCount failed: %v", cmd.Err())
		return
	}

	fmt.Printf("ZLEXCOUNT [-,+] result: %d\n", cmd.Val())
}

// TestRedisClient_ZRemRangeByLex 测试 ZREMRANGEBYLEX 命令
func TestRedisClient_ZRemRangeByLex(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先添加成员（分数相同）
	client.ZAdd(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test16",
		"score1":  0,
		"member1": "alpha",
		"score2":  0,
		"member2": "beta",
		"score3":  0,
		"member3": "gamma",
	})

	// 删除字典序范围
	cmd := client.ZRemRangeByLex(context.Background(), ZSetCmd, map[string]any{
		"keyName": "test16",
		"min":     "[alpha",
		"max":     "[beta",
	})

	if cmd.Err() != nil {
		t.Errorf("ZRemRangeByLex failed: %v", cmd.Err())
		return
	}

	fmt.Printf("ZREMRANGEBYLEX [alpha,[beta result: %d members removed\n", cmd.Val())
}

// TestRedisClient_ZInterStore 测试 ZINTERSTORE 命令
func TestRedisClient_ZInterStore(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 创建两个有序集合
	client.ZAdd(context.Background(), ZSetCmd, map[string]any{
		"keyName": "z1",
		"score1":  1.0,
		"member1": "one",
		"score2":  2.0,
		"member2": "two",
		"score3":  3.0,
		"member3": "three",
	})

	client.ZAdd(context.Background(), ZSetCmd, map[string]any{
		"keyName": "z2",
		"score1":  1.0,
		"member1": "one",
		"score2":  3.0,
		"member2": "two",
		"score3":  4.0,
		"member3": "four",
	})

	// 计算交集并存储
	cmd := client.ZInterStore(context.Background(), ZSetCmd, nil,
		"zset:out",
		2,
		"zset:z1",
		"zset:z2",
		"WEIGHTS", 2, 3,
		"AGGREGATE", "SUM",
	)

	if cmd.Err() != nil {
		t.Errorf("ZInterStore failed: %v", cmd.Err())
		return
	}

	fmt.Printf("ZINTERSTORE result: %d members stored\n", cmd.Val())
}

// TestRedisClient_ZUnionStore 测试 ZUNIONSTORE 命令
func TestRedisClient_ZUnionStore(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 创建两个有序集合
	client.ZAdd(context.Background(), ZSetCmd, map[string]any{
		"keyName": "union1",
		"score1":  1.0,
		"member1": "a",
		"score2":  2.0,
		"member2": "b",
	})

	client.ZAdd(context.Background(), ZSetCmd, map[string]any{
		"keyName": "union2",
		"score1":  2.0,
		"member1": "b",
		"score2":  3.0,
		"member2": "c",
	})

	// 计算并集并存储
	cmd := client.ZUnionStore(context.Background(), ZSetCmd, nil,
		"zset:union_out",
		2,
		"zset:union1",
		"zset:union2",
	)

	if cmd.Err() != nil {
		t.Errorf("ZUnionStore failed: %v", cmd.Err())
		return
	}

	fmt.Printf("ZUNIONSTORE result: %d members stored\n", cmd.Val())
}

// TestRedisClient_ZSet_Integration 集成测试：ZSet 操作的完整流程
func TestRedisClient_ZSet_Integration(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	keyName := "integration_zset"

	// 1. ZADD
	client.ZAdd(context.Background(), ZSetCmd, map[string]any{
		"keyName": keyName,
		"score1":  10.0,
		"member1": "member1",
		"score2":  20.0,
		"member2": "member2",
		"score3":  30.0,
		"member3": "member3",
	})

	// 2. ZCARD
	cardCmd := client.ZCard(context.Background(), ZSetCmd, map[string]any{
		"keyName": keyName,
	})
	fmt.Printf("1. ZCARD: %d\n", cardCmd.Val())

	// 3. ZRANGE
	rangeCmd := client.ZRange(context.Background(), ZSetCmd, map[string]any{
		"keyName": keyName,
		"start":   0,
		"stop":    -1,
	})
	fmt.Printf("2. ZRANGE: %v\n", rangeCmd.Val())

	// 4. ZRANK
	rankCmd := client.ZRank(context.Background(), ZSetCmd, map[string]any{
		"keyName": keyName,
		"member":  "member2",
	})
	fmt.Printf("3. ZRANK member2: %d\n", rankCmd.Val())

	// 5. ZSCORE
	scoreCmd := client.ZScore(context.Background(), ZSetCmd, map[string]any{
		"keyName": keyName,
		"member":  "member2",
	})
	fmt.Printf("4. ZSCORE member2: %s\n", scoreCmd.Val())

	// 6. ZINCRBY
	incrCmd := client.ZIncrBy(context.Background(), ZSetCmd, map[string]any{
		"keyName":   keyName,
		"increment": 5.0,
		"member":    "member2",
	})
	fmt.Printf("5. ZINCRBY member2 +5: %s\n", incrCmd.Val())

	// 7. ZCOUNT
	countCmd := client.ZCount(context.Background(), ZSetCmd, map[string]any{
		"keyName": keyName,
		"min":     15,
		"max":     25,
	})
	fmt.Printf("6. ZCOUNT [15,25]: %d\n", countCmd.Val())
}
