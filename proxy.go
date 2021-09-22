// package proxy will verify if an IPaddr is blocked.
package proxy

import (
	"fmt"
	"net/http"
	"net/url"
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
	u := ExtractDomain(url)
	return b[u]
}

func ExtractDomain(u string) string {
	data, err := url.Parse(u)
	if err != nil {
		return ""
	}
	return data.Host
}

func ListenAsync(addr string, blockList BlockList) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			fmt.Println("blockList is:", blockList)
			fmt.Println("req is:", req)
			fmt.Println("req.RequestURI is:", req.RequestURI)
			if blockList.Blocked(req.RequestURI) {
				fmt.Println("Checking if blocked...")
				http.Error(w, "No soup for you", http.StatusForbidden)
			}
		})
		fmt.Printf("Listening on port %s\n", addr)
		err := http.ListenAndServe(addr, mux)
		if err != nil {
			fmt.Println("Error listening:", err)
		}
	}()
	return &wg
}
