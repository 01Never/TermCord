package main

import (
	"fmt"
	"net/http"
)

func helloWorldhandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World! One day I will replace Discord!!")
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
