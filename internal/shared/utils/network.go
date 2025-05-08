package utils

import (
	"fmt"
	"net"
)

func IsPortInUse(port int) bool {
	address := fmt.Sprintf("0.0.0.0:%d", port)
	ln, err := net.Listen("tcp", address)

	if err != nil {
		return true
	}

	_ = ln.Close()

	return false
}
