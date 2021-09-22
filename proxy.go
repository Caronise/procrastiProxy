// package proxy will verify if an IPaddr is blocked.
package proxy

import (
	"fmt"
	"net/http"
	"sync"
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
	fmt.Printf("Got request for %q\n", req.RequestURI)
	if req.RequestURI == "http://blocked.com/" {
		http.Error(w, "No soup for you", http.StatusForbidden)
	}
}

func ListenAsync(addr string) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		mux := http.NewServeMux()
		mux.HandleFunc("/", redirect)
		fmt.Printf("Listening on port %s\n", addr)
		err := http.ListenAndServe(addr, mux)
		if err != nil {
			fmt.Println("Error listening:", err)
		}
	}()
	return &wg
}