package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func handler(s string) {
	body := strings.NewReader(s)
	resp, err := http.Post("http://localhost:8080/helloWorld", "text/plain", body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

}

func main() {

	// resp, err := http.Get("http://localhost:8080/login")
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	// fmt.Println("Response status:", resp.Status)

	// scanner := bufio.NewScanner(resp.Body)
	// for i := 0; scanner.Scan() && i < 5; i++ {
	// 	fmt.Println(scanner.Text())
	// }
	// if err := scanner.Err(); err != nil {
	// 	panic(err)
	// }

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("You typed:", line)
		handler(line)
	}

}
