package main

import (
	"fmt"
	"net/http"
)

const port string = ":8888"

func redirect(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Accessed homepage...")
	// I should redirect here?
	http.Redirect(w, req, "example.com", 403)
}

func main() {
	http.HandleFunc("/", redirect)

	fmt.Printf("Listening on port %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error listening:", err)
	}
}
