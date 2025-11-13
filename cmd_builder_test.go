package rdb

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
)

// TestBuildCmd 测试 BuildCmd 方法 - 构建命令但不执行
func TestBuildCmd(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	var StringCmd = RdCmd{
		Key: "string:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			GET: {
				Params: "",
			},
		},
	}

	// 构建命令但不执行
	cmd := client.BuildCmd(context.Background(), StringCmd, GET, map[string]any{
		"keyName": "test_build",
	})

	// BuildCmd 现在返回 *redis.Cmd（默认类型）
	if _, ok := cmd.(*redis.Cmd); !ok {
		t.Errorf("Expected *redis.Cmd, got %T", cmd)
		return
	}

	fmt.Printf("BuildCmd 返回类型: %T\n", cmd)
	fmt.Printf("命令已构建，但尚未执行\n")

	// 现在可以自己决定如何执行
	// 例如：client.Client.Process(ctx, cmd)
	// 或者：直接使用链式调用 String()
}

// TestExecuteStringCmd 测试链式调用 String() - 直接返回 *redis.StringCmd
func TestExecuteStringCmd(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	var StringCmd = RdCmd{
		Key: "string:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			GET: {
				Params: "",
			},
			SET: {
				Params: "{{value}}",
			},
		},
	}

	// 先设置一个值
	setCmd := client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": "test_execute",
		"value":   "Hello World",
	}).String()
	if setCmd.Err() != nil {
		t.Errorf("SET failed: %v", setCmd.Err())
		return
	}

	// 使用链式调用 String() 直接获取 *redis.StringCmd
	cmd := client.Get(context.Background(), StringCmd, map[string]any{
		"keyName": "test_execute",
	}).String()
	if cmd.Err() != nil {
		t.Errorf("String() failed: %v", cmd.Err())
		return
	}

	// 直接使用类型安全的方法，无需类型断言
	val, err := cmd.Result()
	if err != nil {
		t.Errorf("Result() failed: %v", err)
		return
	}

	fmt.Printf("String() 返回类型: %T\n", cmd)
	fmt.Printf("值: %s\n", val)
	fmt.Printf("可以直接调用 cmd.Text(): %s\n", cmd.String())
}

// TestExecuteIntCmd 测试链式调用 Int() - 直接返回 *redis.IntCmd
func TestExecuteIntCmd(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	var IntCmd = RdCmd{
		Key: "int:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			INCR: {
				Params: "",
			},
		},
	}

	// 先设置一个值
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": "test_execute_int",
		"value":   "10",
	})

	// 使用链式调用 Int() 直接获取 *redis.IntCmd
	cmd := client.Incr(context.Background(), IntCmd, map[string]any{
		"keyName": "test_execute_int",
	}).Int()
	if cmd.Err() != nil {
		t.Errorf("Int() failed: %v", cmd.Err())
		return
	}

	// 直接使用类型安全的方法
	val, err := cmd.Result()
	if err != nil {
		t.Errorf("Result() failed: %v", err)
		return
	}

	fmt.Printf("Int() 返回类型: %T\n", cmd)
	fmt.Printf("值: %d\n", val)
	fmt.Printf("可以直接调用 cmd.Int64(): %d\n", cmd.Val())
}

// TestExecuteSliceCmd 测试链式调用 Slice() - 直接返回 *redis.SliceCmd
func TestExecuteSliceCmd(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	var HashCmd = RdCmd{
		Key: "hash:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			HGETALL: {
				Params: "",
			},
		},
	}

	// 先设置一些数据
	client.HMSet(context.Background(), HashCmd, map[string]any{
		"keyName": "test_execute_slice",
	}, "name", "Alice", "age", "30")

	// 使用链式调用 Slice() 直接获取 *redis.SliceCmd
	cmd := client.HGetAll(context.Background(), HashCmd, map[string]any{
		"keyName": "test_execute_slice",
	}).Slice()
	if cmd.Err() != nil {
		t.Errorf("Slice() failed: %v", cmd.Err())
		return
	}

	// 直接使用类型安全的方法
	val, err := cmd.Result()
	if err != nil {
		t.Errorf("Result() failed: %v", err)
		return
	}

	fmt.Printf("ExecuteSliceCmd 返回类型: %T\n", cmd)
	fmt.Printf("值: %v\n", val)
	fmt.Printf("可以直接调用 cmd.Slice(): %v\n", cmd.Val())
}

// TestExecuteCmd_Generic 测试泛型方法 ExecuteCmd
func TestExecuteCmd_Generic(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	var StringCmd = RdCmd{
		Key: "string:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			GET: {
				Params: "",
			},
		},
	}

	// 使用泛型方法 ExecuteCmd
	cmd := ExecuteCmd[*redis.StringCmd](&client, context.Background(), StringCmd, GET, map[string]any{
		"keyName": "test_generic",
	})
	if cmd.Err() != nil {
		t.Errorf("ExecuteCmd failed: %v", cmd.Err())
		return
	}

	// 直接使用类型安全的方法
	val, err := cmd.Result()
	if err != nil {
		t.Errorf("Result() failed: %v", err)
		return
	}

	fmt.Printf("ExecuteCmd[*redis.StringCmd] 返回类型: %T\n", cmd)
	fmt.Printf("值: %s\n", val)
}

// TestBuildCmd_ThenExecute 演示构建命令后自己执行
func TestBuildCmd_ThenExecute(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	var StringCmd = RdCmd{
		Key: "string:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			GET: {
				Params: "",
			},
		},
	}

	// 构建命令
	cmd := client.BuildCmd(context.Background(), StringCmd, GET, map[string]any{
		"keyName": "test_build_then_execute",
	})

	// 验证类型（BuildCmd 返回 *redis.Cmd）
	_, ok := cmd.(*redis.Cmd)
	if !ok {
		t.Errorf("Expected *redis.Cmd, got %T", cmd)
		return
	}

	// 自己执行命令
	err := client.Client.Process(context.Background(), cmd)
	if err != nil {
		t.Errorf("Process failed: %v", err)
		return
	}

	// 获取结果
	val, err := cmd.(*redis.Cmd).Result()
	if err != nil {
		t.Errorf("Result() failed: %v", err)
		return
	}

	fmt.Printf("BuildCmd + 自己执行，返回类型: %T\n", cmd)
	fmt.Printf("值: %v\n", val)
}

// TestExecuteCmd_AllTypes 测试所有类型的 Execute 方法
func TestExecuteCmd_AllTypes(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 测试 String()
	var StringCmd = RdCmd{
		Key: "string:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			GET: {
				Params: "",
			},
		},
	}
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": "test_all_types",
		"value":   "test",
	})
	strCmd, _ := client.Get(context.Background(), StringCmd, map[string]any{"keyName": "test_all_types"}).String()
	fmt.Printf("String(): %T\n", strCmd)

	// 测试 Int()
	var IntCmd = RdCmd{
		Key: "int:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			INCR: {
				Params: "",
			},
		},
	}
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": "test_all_types_int",
		"value":   "10",
	})
	intCmd, _ := client.Incr(context.Background(), IntCmd, map[string]any{"keyName": "test_all_types_int"}).Int()
	fmt.Printf("Int(): %T, value: %d\n", intCmd, intCmd.Val())

	// 测试 Slice()
	var HashCmd = RdCmd{
		Key: "hash:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			HGETALL: {
				Params: "",
			},
		},
	}
	client.HMSet(context.Background(), HashCmd, map[string]any{
		"keyName": "test_all_types_slice",
	}, "field1", "value1")
	sliceCmd, _ := client.HGetAll(context.Background(), HashCmd, map[string]any{"keyName": "test_all_types_slice"}).Slice()
	fmt.Printf("Slice(): %T\n", sliceCmd)

	// 测试 Float()
	var FloatCmd = RdCmd{
		Key: "float:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			INCRBYFLOAT: {
				Params: "{{increment}}",
			},
		},
	}
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": "test_all_types_float",
		"value":   "10.5",
	})
	floatCmd, _ := client.IncrByFloat(context.Background(), FloatCmd, map[string]any{
		"keyName":   "test_all_types_float",
		"increment": 2.5,
	}).Float()
	fmt.Printf("Float(): %T, value: %f\n", floatCmd, floatCmd.Val())

	// 测试 Bool()
	var BoolCmd = RdCmd{
		Key: "bool:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			SETNX: {
				Params: "{{value}}",
			},
		},
	}
	boolCmd, _ := client.SetNx(context.Background(), BoolCmd, map[string]any{
		"keyName": "test_all_types_bool",
		"value":   "test",
	}).Bool()
	fmt.Printf("Bool(): %T, value: %v\n", boolCmd, boolCmd.Val())
}
