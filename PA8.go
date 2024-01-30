package main

import (
	"fmt"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {

	// Extracting the path from the request
	filePath := "." + r.URL.Path

	// Checking if the file exists
	if _, err := os.Stat(filePath); err == nil {
		// If the file exists, serve it
		http.ServeFile(w, r, filePath) //replies to req. with content of file
	} else {
		// If the file doesn't exist, DO NOT return 404! We return File Not Found
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Sorry Polly..., not here. \n")
	}
}

func main() {
	port := "127.0.0.1:12013"

	// Setting up the HTTP server
	http.HandleFunc("/", handler)
	fmt.Printf("Listening on port %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
