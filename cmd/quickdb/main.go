package main

import (
	"fmt"
	"net"
	"strings"

	aof "github.com/yagyagoel1/quickdb/internal/AOF"
	"github.com/yagyagoel1/quickdb/internal/handler"
	"github.com/yagyagoel1/quickdb/utils"
)

func handleConnection(conn net.Conn, aofInstance *aof.Aof) {
	defer conn.Close()
	for {
		resp := utils.NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println("Connection error:", err)
			break
		}

		fmt.Println("Received:", value)
		if value.Typ != "array" {
			fmt.Println("invalid request, expected array")
			continue
		}
		if len(value.Array) == 0 {
			fmt.Println("invalid request, expected array length more than 0")
			continue
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]
		writer := utils.NewWriter(conn)
		cmdHandler, ok := handler.Handlers[command]
		if !ok {
			fmt.Println("invalid command", command)
			writer.Write(utils.Value{Typ: "string", Str: ""})
			continue
		}
		if command == "SET" || command == "HSET" {
			aofInstance.Write(value)
		}

		result := cmdHandler(args)
		fmt.Println("result", result)
		writer.Write(result)
	}
}

func main() {
	fmt.Println("Listening on port 6379")
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Listener error:", err)
		return
	}
	defer l.Close()

	aofInstance, err := aof.NewAof("database.aof")
	if err != nil {
		fmt.Println("AOF error:", err)
		return
	}
	defer aofInstance.Close()

	aofInstance.Read(func(value utils.Value) {
		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]
		cmdHandler, ok := handler.Handlers[command]
		if !ok {
			fmt.Println("Invalid command:", command)
			return
		}
		cmdHandler(args)
	})

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		go handleConnection(conn, aofInstance)
	}
}
