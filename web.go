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
	fs := http.FileServer(http.Dir(filePath))
	http.Handle(urlPath, http.StripPrefix(urlPath, fs))

	// TODO: make sure it handles URLs missing ".html" at the end
	// https://stackoverflow.com/questions/57281010/remove-the-html-extension-from-every-file-in-a-simple-http-server
	//
	// fs := http.FileServer(http.Dir("."))
	// http.HandleFunc(urlPath, func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Printf("Processing: %s\r\n", r.URL.Path)
	// 	if strings.HasSuffix(r.URL.Path, "/") {
	// 		r.URL.Path += "index.html"
	// 		fmt.Printf("Adding index.html\r\n")
	// 	} else if ext := path.Ext(r.URL.Path); ext == "" {
	// 		r.URL.Path += ".html"
	// 		fmt.Printf("Adding .html\r\n")
	// 	}
	// 	fmt.Printf("Serving: %s\r\n", r.URL.Path)
	// 	fs.ServeHTTP(w, r)
	// })

	// Start the web server
	webAddress := fmt.Sprintf("localhost:%d", port)
	log.Printf("Listening for requests at: http://%s", webAddress)
	log.Printf("Serving files from: %s", filePath)
	err := http.ListenAndServe(webAddress, nil)
	if err != nil {
		log.Fatal("Failed to start the web server: ", err)
	}
}
