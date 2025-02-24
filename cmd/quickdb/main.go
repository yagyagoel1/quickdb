package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/yagyagoel1/quickdb/cmd/quickdb/handler"
	"github.com/yagyagoel1/quickdb/utils"
)

func main() {
	fmt.Println("listening on port 6379")
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	for {
		resp := utils.NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(value)
		if value.Typ != "array" {
			fmt.Println("invalid request , expected array")
		}
		if len(value.Array) == 0 {
			fmt.Println("invalid requrest expected array length more than 0")
		}
		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		writer := utils.NewWriter(conn)
		handler, ok := handler.Handlers[command]
		if !ok {
			fmt.Println("invalid command", command)
			writer.Write(utils.Value{
				Typ: "string", Str: "",
			})
			continue
		}
		result := handler(args)
		writer.Write(result)
		// ignore request and send back a PONG
		// conn.Write([]byte("+OK\r\n"))
	}

}
