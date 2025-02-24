package main

import (
	"fmt"
	"net"

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

		_ = value

		writer := utils.NewWriter(conn)
		writer.Write(utils.Value{Typ: "string", Str: "OK"})
		// ignore request and send back a PONG
		conn.Write([]byte("+OK\r\n"))
	}

}
