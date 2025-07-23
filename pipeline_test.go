package rdb

import (
	"fmt"
	base "github.com/preceeder/go.base"
	"testing"
)

func TestRedisClient_PipeLine(t *testing.T) {
	client := InitRedis()
	pip := client.PipeLine()
	add := pip.ZAdd(base.Context{}, ade, map[string]any{
		"score1": 3, "key1": 2, "score2": 5, "key2": 4, "score3": 5, "key3": 6,
	})

	zer := pip.ZRange(base.Context{}, ade, map[string]any{"start": 0, "stop": -1})

	ds, err := pip.Exec(base.Context{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ds)
	fmt.Println(add.Val())
	fmt.Println(zer.Val())
}
