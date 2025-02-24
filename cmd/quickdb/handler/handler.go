package handler

import (
	"fmt"
	"sync"

	"github.com/yagyagoel1/quickdb/utils"
)

var Handlers = map[string]func([]utils.Value) utils.Value{
	"PING":    Ping,
	"SET":     Set,
	"GET":     get,
	"HSET":    Hset,
	"HGET":    Hget,
	"HGETALL": HgetAll,
}

func Ping(args []utils.Value) utils.Value {
	fmt.Println("ARSG", args)
	if len(args) == 0 {
		return utils.Value{Typ: "string", Str: "PONG"}
	}
	return utils.Value{Typ: "string", Str: args[0].Bulk}

}

var SETs = map[string]string{}

var SETsMu = sync.RWMutex{}

func Set(args []utils.Value) utils.Value {
	if len(args) != 2 {
		return utils.Value{Typ: "error", Str: "Err wrong number of arguments for set "}
	}
	key := args[0].Bulk
	value := args[1].Bulk
	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return utils.Value{Typ: "string", Str: "OK"}
}

func get(args []utils.Value) utils.Value {
	if len(args) != 1 {
		return utils.Value{Typ: "error", Str: "ERR wrong number of arguments for a get command "}

	}
	key := args[0].Bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	if !ok {
		return utils.Value{Typ: "null"}
	}
	return utils.Value{Typ: "bulk", Bulk: value}

}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func Hset(args []utils.Value) utils.Value {
	if len(args) != 3 {
		return utils.Value{Typ: "error", Str: "ERR wrong no of arguments for HSET"}

	}
	outerKey := args[0].Bulk
	innerKey := args[1].Bulk
	value := args[2].Bulk

	HSETsMu.Lock()
	if _, ok := HSETs[outerKey]; !ok {
		HSETs[outerKey] = map[string]string{}
	}
	HSETs[outerKey][innerKey] = value
	HSETsMu.Unlock()

	return utils.Value{Typ: "string", Str: "OK"}
}

func Hget(args []utils.Value) utils.Value {
	if len(args) != 2 {
		return utils.Value{Typ: "error", Str: "ERR wrong no of arguments for HGET"}

	}
	outerKey := args[0].Bulk
	innerKey := args[1].Bulk
	HSETsMu.RLock()
	value, ok := HSETs[outerKey][innerKey]
	if !ok {
		return utils.Value{Typ: "null"}
	}
	return utils.Value{Typ: "bulk", Str: value}
}

func HgetAll(args []utils.Value) utils.Value {
	if len(args) != 1 {
		return utils.Value{Typ: "error", Str: "ERR wrong no of arguments for HGETALL"}

	}
	outerKey := args[0].Bulk
	HSETsMu.RLock()
	value, ok := HSETs[outerKey]
	if !ok {
		return utils.Value{Typ: "null"}
	}
	arr := utils.Value{Typ: "array", Array: make([]utils.Value, 0)}
	for i := range value {
		arr.Array = append(arr.Array, utils.Value{Typ: "bulk", Str: HSETs[outerKey][i]})
	}
	return arr
}
