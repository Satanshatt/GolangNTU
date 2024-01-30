package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

/* CONCURRENT SERVER
PA5.go handles clients’ file upload requests sequentially. That is, it will complete a
file upload before starting another one. If the upload time is long, the waiting time
for the clients farther back in line can be unbearably long. What’s nice about Go is
that it is easy to specify a code block for concurrent execution. That concurrent code
block is referred to as a Goroutine. As a result, processing of later file upload
requests can be started before the earlier file upload requests are completed.
Although the completion time as a whole (time to finish all file uploads) might be the
same, clients are generally happier to see their requested being started early. In
other words, with concurrency, we get better user experience, while the system
capacity remains unchanged.
You will turn the file uploading part of your PA5.go concurrent in PA6.go. This is the
specification of the server in PA5.go.

(1) Listens at <your port#> until there’s an upload request
(2) reads from the socket first the file size (just the number in a single line)
(3) reads from the socket one line at a time
(4) prepend the line count to each line and store the new line into a new
file: whatever.txt
(5) repeats (3) and (4) until all lines in the file is processed
(6) sends a message back that tells the client the original file and the new
file size
(7) closes the connections and goes back to (1)
*/

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		panic(err)
	}
}

func handleConnection(conn net.Conn, process int) {

	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	fileSizeStr := scanner.Text()

	fileSize, err := strconv.Atoi(fileSizeStr)
	check(err)

	newFile, err := os.Create("whatever.txt")
	check(err)
	defer newFile.Close()

	lineCount := 1
	tmpFileSize := fileSize

	for tmpFileSize > 0 {

		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		newLine := fmt.Sprintf("%d: %s\n", lineCount, line)
		//fmt.Printf("%s", newLine)
		lineCount++
		tmpFileSize -= len(line) + 1

		//fmt.Printf("%d: tmpfileSize: %d\n", process, tmpFileSize)
		_, err := newFile.WriteString(newLine)
		check(err)
	}

	// Step 6: Send a message back
	originalFileSize := fileSize
	newFileSize, _ := newFile.Stat()
	response := fmt.Sprintf("Process %d: Original file size: %d bytes New file size: %d bytes \n ", process, originalFileSize, newFileSize.Size())
	print(response)
	conn.Write([]byte(response + "\n"))

	time.Sleep(5 * time.Second)

	// Step 7: Close the connection
	newFile.Close()
	conn.Close()
	fmt.Printf("\n Process %d: File uploaded and processed. Connection closed.\n", process)
}

func main() {

	port := "127.0.0.1:12013"

	//starting server
	fmt.Println("Launching server...")
	listener, err := net.Listen("tcp", port)
	check(err)
	defer listener.Close()

	//for loop to keep connection alive
	i := 1
	for {
		conn, _ := listener.Accept()
		fmt.Printf("%d Listening on port %s...\n", i, port)

		go handleConnection(conn, i)
		i++
	}

}
