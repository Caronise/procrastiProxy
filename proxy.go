// package proxy will verify if an IPaddr is blocked.
package proxy

type BlockList []string

// func Blocked will check if ipAddr is blocked by ranging over BlockList.
func Blocked(ipAddr string, BlockList BlockList) bool {
	for _, listItem := range BlockList {
		if ipAddr == listItem {
			return true
		}
	}
	return false
}
