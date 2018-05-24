package test

import (
	"fmt"
	"runtime"
	"testing"
)

func TestRuntimeCall(t *testing.T) {
	for skip := 0; ; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		fmt.Printf("skip = %v, pc = %v, file = %v, line = %v\n", skip, pc, file, line)
	}
}
