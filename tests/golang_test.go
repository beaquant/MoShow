package test

import (
	"fmt"
	"runtime"
	"testing"
	"time"
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

func TestSelect(t *testing.T) {
	ticker := time.NewTicker(time.Second)

	fmt.Println("start")
	for i := 0; ; i++ {
		select {
		case <-ticker.C:
			if i < 5 {
				fmt.Println("continue")
				continue
			}

			fmt.Println("break")
			break

			if i > 10 {
				fmt.Println("return")
				return
			}
		}
	}
}
