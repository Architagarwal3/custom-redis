package main

import (
	"bufio"
	"custom-redis/services"
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Client connected:", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		response, err := services.HandleMessage(scanner.Text())
		if err != nil {
			response = "ERR " + err.Error()
		}

		_, err = fmt.Fprintln(conn, "Echo: "+response)
		if err != nil {
			fmt.Println("Write error:", err)
			return
		}
	}

	// Check for errors after the Scan loop
	if err := scanner.Err(); err != nil {
		fmt.Println("Read error:", err)
	}

	fmt.Println("Client disconnected:", conn.RemoteAddr())
}
