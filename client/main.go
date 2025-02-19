package client

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func Main() {
	input := "$5\r\nyagya\r\n"
	reader := bufio.NewReader(strings.NewReader(input))

	b, err := reader.ReadByte()
	if err != nil {
		log.Fatalf("error reading the byte %v", err)
	}

	if string(b) != "$" {
		fmt.Println("invalid type , expecting bulk strings only ")
		os.Exit(1)
	}
	size, err := reader.ReadByte()
	if err != nil {
		log.Fatal("unable to read the size")
	}
	strSize, err := strconv.ParseInt(string(size), 10, 64)
	if err != nil {
		log.Fatal("unable to parse the size ")
	}
	reader.ReadByte()
	reader.ReadByte()
	name := make([]byte, strSize)
	reader.Read(name)

	fmt.Println(string(name))
	fmt.Println(string(name))

}
