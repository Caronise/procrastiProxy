package proxy_test

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"proxy"
)

func TestNotBlocked(t *testing.T) {
	t.Parallel()
	var blocked bool
	url := "example.com"
	blockList := proxy.NewBlockList()
	blockList.UnBlock(url)

	blocked = blockList.Blocked(url)
	if blocked {
		t.Fatal("It shouldn't be blocked... but it was!")
	}
}

func TestBlocked(t *testing.T) {
	t.Parallel()
	var blocked bool
	url := "example.com"
	blockList := proxy.NewBlockList()
	blockList.Block(url)

	blocked = blockList.Blocked(url)
	if !blocked {
		t.Fatal("It should be blocked... but it wasn't.")
	}
}

func TestRequestToUnblockedURLReturns200(t *testing.T) {
	t.Parallel()

	port := ":8001"

	go proxy.Listener(port)

	// make a request to the proxy
	for i := 0; i < 3; i++ {
		conn, err := net.Dial("tcp", "127.0.0.1"+port)
		if err == nil {
			fmt.Println("Server is up!")
			conn.Close()
			break
		}
		fmt.Println("Didn't connect!")
		time.Sleep(time.Millisecond * 30)
	}

	response, err := http.Get("http://127.0.0.1" + port)
	if err != nil {
		t.Fatal("Error getting address", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatal("Status is not 200")
	}

}

func TestRequestToBlockURLReturns403(t *testing.T) {
	t.Parallel()

	port := ":8002"

	go proxy.Listener(port)

	// make a request to the proxy
	for i := 0; i < 3; i++ {
		conn, err := net.Dial("tcp", "127.0.0.1"+port)
		if err == nil {
			fmt.Println("Server is up!")
			conn.Close()
			break
		}
		fmt.Println("Didn't connect!")
		time.Sleep(time.Millisecond * 30)
	}

	response, err := http.Get("http://127.0.0.1" + port)
	if err != nil {
		t.Fatal("Error getting address", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusForbidden {
		t.Fatal("Status is not 403")
	}
}
