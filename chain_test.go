package rdb

import (
	"context"
	"fmt"
	"testing"
)

// TestChainCall_Set_ExecuteStringCmd 测试链式调用：Set().ExecuteStringCmd()
func TestChainCall_Set_ExecuteStringCmd(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	var StringCmd = RdCmd{
		Key: "string:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			SET: {
				Params: "{{value}}",
			},
			GET: {
				Params: "",
			},
		},
	}

	// 链式调用：Set().String()
	cmd := client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": "test_chain",
		"value":   "Hello World",
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

	fmt.Printf("链式调用 Set().ExecuteStringCmd() 返回类型: %T\n", cmd)
	fmt.Printf("值: %s\n", val)
	fmt.Printf("可以直接调用 cmd.Text(): %s\n", cmd.String())
}

// TestChainCall_Get_ExecuteStringCmd 测试链式调用：Get().ExecuteStringCmd()
func TestChainCall_Get_ExecuteStringCmd(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	var StringCmd = RdCmd{
		Key: "string:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			SET: {
				Params: "{{value}}",
			},
			GET: {
				Params: "",
			},
		},
	}

	// 先设置一个值
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": "test_chain_get",
		"value":   "Chain Test Value",
	})

	// 链式调用：Get().String()
	cmd := client.Get(context.Background(), StringCmd, map[string]any{
		"keyName": "test_chain_get",
	}).String()

	if cmd.Err() != nil {
		t.Errorf("String() failed: %v", cmd.Err())
		return
	}

	val, err := cmd.Result()
	if err != nil {
		t.Errorf("Result() failed: %v", err)
		return
	}

	fmt.Printf("链式调用 Get().ExecuteStringCmd() 返回类型: %T\n", cmd)
	fmt.Printf("值: %s\n", val)
}

// TestChainCall_Incr_ExecuteIntCmd 测试链式调用：Incr().ExecuteIntCmd()
func TestChainCall_Incr_ExecuteIntCmd(t *testing.T) {
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
		"keyName": "test_chain_int",
		"value":   "10",
	})

	// 链式调用：Incr().Int()
	cmd := client.Incr(context.Background(), IntCmd, map[string]any{
		"keyName": "test_chain_int",
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

	fmt.Printf("链式调用 Incr().ExecuteIntCmd() 返回类型: %T\n", cmd)
	fmt.Printf("值: %d\n", val)
	fmt.Printf("可以直接调用 cmd.Int64(): %d\n", cmd.Val())
}

// TestChainCall_HGetAll_ExecuteSliceCmd 测试链式调用：HGetAll().ExecuteSliceCmd()
func TestChainCall_HGetAll_ExecuteSliceCmd(t *testing.T) {
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
		"keyName": "test_chain_hash",
	}, "name", "Alice", "age", "30", "city", "Beijing")

	// 链式调用：HGetAll().Slice()
	cmd := client.HGetAll(context.Background(), HashCmd, map[string]any{
		"keyName": "test_chain_hash",
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

	fmt.Printf("链式调用 HGetAll().ExecuteSliceCmd() 返回类型: %T\n", cmd)
	fmt.Printf("值: %v\n", val)
	fmt.Printf("可以直接调用 cmd.Slice(): %v\n", cmd.Val())
}

// TestChainCall_MultipleTypes 测试多种类型的链式调用
func TestChainCall_MultipleTypes(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 测试 StringCmd
	var StringCmd = RdCmd{
		Key: "string:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			GET: {
				Params: "",
			},
		},
	}
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": "test_multi",
		"value":   "test",
	})
	strCmd := client.Get(context.Background(), StringCmd, map[string]any{"keyName": "test_multi"}).String()
	fmt.Printf("StringCmd: %T, value: %s\n", strCmd, strCmd.String())

	// 测试 IntCmd
	var IntCmd = RdCmd{
		Key: "int:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			INCR: {
				Params: "",
			},
		},
	}
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": "test_multi_int",
		"value":   "10",
	})
	intCmd := client.Incr(context.Background(), IntCmd, map[string]any{"keyName": "test_multi_int"}).Int()
	fmt.Printf("IntCmd: %T, value: %d\n", intCmd, intCmd.Val())

	// 测试 FloatCmd
	var FloatCmd = RdCmd{
		Key: "float:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			INCRBYFLOAT: {
				Params: "{{increment}}",
			},
		},
	}
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": "test_multi_float",
		"value":   "10.5",
	})
	floatCmd := client.IncrByFloat(context.Background(), FloatCmd, map[string]any{
		"keyName":   "test_multi_float",
		"increment": 2.5,
	}).Float()
	fmt.Printf("FloatCmd: %T, value: %f\n", floatCmd, floatCmd.Val())

	// 测试 BoolCmd
	var BoolCmd = RdCmd{
		Key: "bool:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			SETNX: {
				Params: "{{value}}",
			},
		},
	}
	boolCmd := client.SetNx(context.Background(), BoolCmd, map[string]any{
		"keyName": "test_multi_bool",
		"value":   "test",
	}).Bool()
	fmt.Printf("BoolCmd: %T, value: %v\n", boolCmd, boolCmd.Val())
}

// TestChainCall_Execute 测试通用的 Execute 方法
func TestChainCall_Execute(t *testing.T) {
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

	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": "test_execute",
		"value":   "Hello",
	})

	// 使用链式调用 String() 获取具体类型
	cmd := client.Get(context.Background(), StringCmd, map[string]any{
		"keyName": "test_execute",
	}).String()

	if cmd.Err() != nil {
		t.Errorf("String() failed: %v", cmd.Err())
		return
	}

	val, _ := cmd.Result()
	fmt.Printf("String() 返回类型: %T, 值: %s\n", cmd, val)
}
