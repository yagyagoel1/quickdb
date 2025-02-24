package handler

import (
	"fmt"
	"sync"

	"github.com/yagyagoel1/quickdb/utils"
)

var Handlers = map[string]func([]utils.Value) utils.Value{
	"PING": Ping,
	"SET":  Set,
	"GET":  get,
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
