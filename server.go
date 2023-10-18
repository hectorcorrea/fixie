package main

import (
	"fmt"
	"log"
	"net/http"
)

func server(port int) {
	filePath := "."
	urlPath := "/"

	// Setup the request handler
	// https://stackoverflow.com/questions/57281010/remove-the-html-extension-from-every-file-in-a-simple-http-server
	fs := http.FileServer(http.Dir("."))
	http.HandleFunc(urlPath, func(w http.ResponseWriter, r *http.Request) {
		filename := "." + r.URL.Path
		fmt.Printf("Processing: %s\r\n", filename)
		if fileExist(filename) {
			fs.ServeHTTP(w, r)
			return
		}

		if fileExist(filename + ".html") {
			fmt.Printf("Processing: %s (defaulted to .html)\r\n", filename+".html")
			r.URL.Path += ".html"
			fs.ServeHTTP(w, r)
			return
		}

		fmt.Printf("File does not exist: %s\r\n", filename)
		fs.ServeHTTP(w, r)
	})

	// Start the web server
	webAddress := fmt.Sprintf("localhost:%d", port)
	fmt.Printf("\r\nLoading the local web server\r\n")
	fmt.Printf("Listening for requests at: http://%s\r\n", webAddress)
	fmt.Printf("Serving files from: %s\r\n", filePath)
	err := http.ListenAndServe(webAddress, nil)
	if err != nil {
		log.Fatal("Failed to start the web server: ", err)
	}
}
