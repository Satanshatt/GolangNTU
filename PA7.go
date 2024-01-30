package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		panic(err)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read the HTTP request
	request, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading from connection: %v\n", err)
		return
	}

	// Parse the request to get the path and file name
	requestParts := strings.Fields(request)
	if len(requestParts) < 2 {
		fmt.Println("Malformed HTTP request")
		return
	}

	filePath := strings.TrimPrefix(requestParts[1], "/")

	// Check if the file exists
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Println("File not found")
		return
	} else if err != nil {
		fmt.Printf("Error getting file information: %v\n", err)
		return
	}

	// Print the file size
	fmt.Printf("File size of %s: %d bytes\n", filePath, fileInfo.Size())
}

func main() {
	port := "127.0.0.1:12013"

	fmt.Println("Launching server...")

	listener, err := net.Listen("tcp", port)
	check(err)
	defer listener.Close()

	for {
		fmt.Printf("Listening on port %s...\n", port)
		conn, err := listener.Accept()
		check(err)

		go handleConnection(conn)
	}
}
