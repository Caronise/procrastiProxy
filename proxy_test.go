package proxy_test

import (
	"testing"

	"proxy"
)

/* TODOS:
- Add package and function descriptions
- Should refactor blockList into a common object
- What should I see if an addr is blocked? What should happen?
	-> it should not allow the connection to go through and print a message.
	-> check the http header and modifying it so it won't go through?
*/

func TestNotBlocked(t *testing.T) {

	var blocked bool
	ipAddr := "example.com"
	blockList := []string{"one.com", "two.com", "three.com"}

	blocked = proxy.Blocked(ipAddr, blockList)
	if blocked {
		t.Fatal("It shouldn't be blocked... but it was!")
	}

}

func TestBlocked(t *testing.T) {

	var blocked bool
	ipAddr := "example.com"
	blockList := []string{"example.com", "two.com", "three.com"}

	blocked = proxy.Blocked(ipAddr, blockList)
	if !blocked {
		t.Fatal("It should be blocked... but it wasn't.")
	}
}
