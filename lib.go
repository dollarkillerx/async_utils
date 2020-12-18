package async_utils

import (
	"fmt"
	"runtime"
)

func PrintStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	fmt.Printf("==> %s\n", string(buf[:n]))
}
