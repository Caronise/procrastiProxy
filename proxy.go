// package proxy will verify if an IPaddr is blocked.
package proxy

import (
	"fmt"
	"net/http"
)

type BlockList map[string]bool

func NewBlockList() BlockList {
	return BlockList{}
}

func (b BlockList) Block(url string) {
	b[url] = true
}

func (b BlockList) UnBlock(url string) {
	b[url] = false
}

// func Blocked will check if ipAddr is blocked by ranging over BlockList.
func (b BlockList) Blocked(url string) bool {
	return b[url]
}

func redirect(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Accessed homepage...")
	//http.Redirect(w, req, "example.com", 403)
}

func Listener(port string) {
	http.HandleFunc("/", redirect)

	fmt.Printf("Listening on port %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error listening:", err)
	}

}
