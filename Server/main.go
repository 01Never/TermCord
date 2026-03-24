package main

import (
	"fmt"
	"io"
	"net/http"
)

func helloWorldhandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World! One day I will replace Discord!!")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	fmt.Println(string(body))
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "TermCord is a fully terminal-based communication platform designed for speed, simplicity, privacy and control.")
	case "POST":
		fmt.Fprintf(w, "Brewing some coffe...and a great terminal-based chat!")
	}
}

func main() {
	http.HandleFunc("/helloWorld", helloWorldhandler)
	http.HandleFunc("/about", aboutHandler)
	http.ListenAndServe(":8080", nil)
}
