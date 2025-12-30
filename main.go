package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("Listening on port: 6379")

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

	// for {
	// 	buf := make([]byte, 1024)

	// 	// read message from client
	// 	_, err = conn.Read(buf)
	// 	fmt.Println(string(buf))
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		fmt.Println("error reading from client: ", err.Error())
	// 		os.Exit(1)
	// 	}

	// 	// ignore request and send back a PONG
	// 	conn.Write([]byte("+OK\r\n"))
	// }

	for {
		resp := NewResp(conn)
		value, err := resp.Read()

		if err != nil {
			fmt.Println(err)
			return
		}
		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}
		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]
		writer := NewWriter(conn)
		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command:", command)
			writer.Writer(Value{typ: "string", str: ""})
			continue
		}
		result := handler(args)
		writer.Writer(result)
	}

}
