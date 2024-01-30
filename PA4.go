package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

/* UPLOAD SERVER
Having experience programming the file upload client, you are ready to program the
file upload server. Your PA4.go should allow your PA3.go to upload a file. More
specifically, the server:
(1) listens at <your port#> until there’s an upload request
(2) reads from the socket first the file size (just the number in a single line)
(3) reads from the socket one line at a time
(4) prepend the line count to each line and store the new line into a new
file: whatever.txt
(5) repeats (3) and (4) until all lines in the file is processed
(6) sends a message back that tells the client the original file and the new
file size
(7) closes the connections and terminates the program
*/

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		panic(err)
	}
}

func main() {

	port := "127.0.0.1:12013"

	//starting server
	fmt.Println("Launching server...")
	//listening at this port for requests
	listener, err := net.Listen("tcp", port)
	check(err)
	defer listener.Close()

	fmt.Printf("Listening on port %s...\n", port)

	conn, err := listener.Accept()
	check(err)
	defer conn.Close()

	//First reading the file size from the socket
	//reader := bufio.NewReader(conn)
	fileSizeStr, err := bufio.NewReader(conn).ReadString('\n')
	fileSizeStr = strings.TrimSpace(fileSizeStr) // Trim leading/trailing whitespace
	check(err)

	fileSize, err := strconv.Atoi(fileSizeStr)
	check(err)
	//fmt.Printf("%s", fileSize)
	fmt.Printf("file size: %d\n", fileSizeStr)

	//creating file for output
	newFile, err := os.Create("whatever.txt")
	check(err)
	defer newFile.Close()

	lineCount := 1
	tmpFileSize := fileSize
	scanner := bufio.NewScanner(conn)

	for tmpFileSize > 0 {
		scanner.Scan()
		line := scanner.Text()
		newLine := fmt.Sprintf("%d: %s\n", lineCount, line)
		//fmt.Printf("%s", newLine)
		lineCount++
		tmpFileSize -= len(line) + 1
		//fmt.Printf("tmpfileSize: %d\n", tmpFileSize)
		_, err := newFile.WriteString(newLine)
		check(err)
	}

	// Step 6: Send a message back
	originalFileSize := fileSize
	newFileSize, _ := newFile.Stat()
	response := fmt.Sprintf("Original file size: %d bytes, New file size: %d bytes", originalFileSize, newFileSize.Size())
	conn.Write([]byte(response + "\n"))

	// Step 7: Close the connection
	fmt.Printf("File uploaded and processed. Connection closed.\n")

	// close connection & terminate
	conn.Close()
}
