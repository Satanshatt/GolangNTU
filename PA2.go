package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	var inputFilename, outputFilename string

	fmt.Print("Input filename: ")
	_, err := fmt.Scan(&inputFilename)
	check(err)

	fmt.Print("Output filename: ")
	_, err = fmt.Scan(&outputFilename)
	check(err)

	inputFile, err := os.Open(inputFilename)
	check(err)
	defer inputFile.Close()

	outputFile, err := os.Create(outputFilename)
	check(err)
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	scanner := bufio.NewScanner(inputFile) //creating a buffer I/O, scanner object

	lineCount := 1

	for scanner.Scan() { //scanning from f, looping until end of file

		line := scanner.Text()
		//var lineNumber int = 1
		//len, _ := writer.WriteString(scanner.Text())
		//fmt.Println(len)
		lineWithCount := strconv.Itoa(lineCount) + ": " + line + "\n"
		_, err := writer.WriteString(lineWithCount)
		check(err)
		lineCount++
	}

	err = writer.Flush() //enforce the string temporarily stored on the system memory to the file on the disk
	check(err)

}
