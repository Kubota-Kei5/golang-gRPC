package interceptor

import "fmt"

func Logger(format string, args ...interface{}) {
	fmt.Printf("[server] "+format+"\n", args...)
}