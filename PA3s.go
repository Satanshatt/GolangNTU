package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e.Error())
		panic(e)
	}
}

func main() {

	var inputFilename string
	//140.112.42.221:12000
	//conn, errc := net.Dial("tcp", "127.0.0.1:12013") //dial to connect to client, return socket handle & error msg
	conn, errc := net.Dial("tcp", "127.0.0.1:12013") //dial to connect to client, return socket handle & error msg
	check(errc)
	defer conn.Close() //defer closing until program finished

	//asks for input file
	fmt.Print("Input filename: ")      //ask for name of file
	_, err := fmt.Scan(&inputFilename) //read for input of file
	check(err)

	//opening the file to read content
	inputFile, err := os.Open(inputFilename)
	check(err)
	defer inputFile.Close()

	//get data of file
	fileInformation, err := os.Stat(inputFilename)
	fileSize := strconv.FormatInt(fileInformation.Size(), 10)

	var string_buffer bytes.Buffer

	_, errc = conn.Write([]byte(fileSize + "\r\n"))
	fmt.Print("Sending size of file: ", fileSize)
	check(errc)

	//scanning content of input file,
	fileScanner := bufio.NewScanner(inputFile)

	//should next send the content of the file to the server
	for fileScanner.Scan() {
		string_buffer.WriteString(fileScanner.Text() + "\r\n")
	}

	//fmt.Printf(string_buffer.String())

	// writing the array to conn to server
	conn.Write(string_buffer.Bytes())

	//scans for message from server, then prints it
	serverScanner := bufio.NewScanner(conn)

	if serverScanner.Scan() {
		fmt.Printf("Server replies: %s\r\n", serverScanner.Text())
	}

	//close connection & terminate
	conn.Close()

}
