package main

import (
	"fmt"
	"net"
	"strings"

	aof "github.com/yagyagoel1/quickdb/internal/AOF"
	"github.com/yagyagoel1/quickdb/internal/handler"
	"github.com/yagyagoel1/quickdb/utils"
)

func main() {
	fmt.Println("listening on port 6379")
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	aof, err := aof.NewAof("database.aof")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	aof.Read(func(value utils.Value) {
		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]
		handler, ok := handler.Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args)
	})
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
		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}

		result := handler(args)
		writer.Write(result)
		// ignore request and send back a PONG
		// conn.Write([]byte("+OK\r\n"))
	}

}
