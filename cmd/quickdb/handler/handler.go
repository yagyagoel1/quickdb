package handler

import (
	"fmt"

	"github.com/yagyagoel1/quickdb/utils"
)

var Handlers = map[string]func([]utils.Value) utils.Value{
	"PING": Ping,
}

func Ping(args []utils.Value) utils.Value {
	fmt.Println("ARSG", args)
	return utils.Value{Typ: "string", Str: "PONG"}

}
