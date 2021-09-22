package proxy_test

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"proxy"

	"github.com/phayes/freeport"
)

var debugOutput = os.Stdout

// var debugOutput = io.Discard

func debug(msg string) {
	fmt.Fprintln(debugOutput, msg)
}

func TestNotBlocked(t *testing.T) {
	t.Parallel()
	var blocked bool
	domain := "example.com"
	blockList := proxy.NewBlockList()
	blockList.UnBlock(domain)

	url := "http://example.com/"
	blocked = blockList.Blocked(url)
	if blocked {
		t.Fatal("It shouldn't be blocked... but it was!")
	}
}

func TestBlocked(t *testing.T) {
	t.Parallel()
	var blocked bool
	domain := "example.com"
	blockList := proxy.NewBlockList()
	blockList.Block(domain)

	url := "http://example.com/"
	// Need to check if url contains the domain
	blocked = blockList.Blocked(url)
	if !blocked {
		t.Fatal("It should be blocked... but it wasn't.")
	}
}

func TestDomain(t *testing.T) {
	t.Parallel()
	url := "http://example.com/locations/"
	domain := "example.com"
	d := proxy.ExtractDomain(url)
	if d != domain {
		t.Fatal("Domain did not match, extracted domain.")
	}
}

func startAndWaitForServer(blockList proxy.BlockList) (url string, err error) {
	port, err := freeport.GetFreePort()
	if err != nil {
		return "", err
	}
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	wg := proxy.ListenAsync(addr, blockList)
	wg.Wait()
	// make a request to the proxy
	for i := 0; i < 3; i++ {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			debug("Server is up!")
			conn.Close()
			break
		}
		debug("Didn't connect!")
		time.Sleep(time.Millisecond * 30)
	}
	return "http://" + addr, nil
}

func proxifiedClientForServer(blockList proxy.BlockList) (*http.Client, error) {
	serverURL, err := startAndWaitForServer(blockList)
	if err != nil {
		return nil, err
	}
	proxyUrl, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}, nil
}

func TestRequestToUnblockedURLReturns200(t *testing.T) {
	t.Parallel()
	url := "http://example.com"
	blockList := proxy.NewBlockList()

	c, err := proxifiedClientForServer(blockList)
	if err != nil {
		t.Fatal(err)
	}

	response, err := c.Get(url)
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
	url := "blocked.com"
	blockList := proxy.NewBlockList()
	blockList.Block(url)

	c, err := proxifiedClientForServer(blockList)
	if err != nil {
		t.Fatal(err)
	}

	target := "http://blocked.com/locations"
	response, err := c.Get(target)
	if err != nil {
		t.Fatal("Error getting address", err)
	}

	defer response.Body.Close()

	fmt.Println("StatusCode is:", response.StatusCode)

	if response.StatusCode != http.StatusForbidden {
		t.Fatal("Status is not 403")
	}
}
