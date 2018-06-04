package test

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

type mutexStruct struct {
	mutex sync.Mutex
}

func TestRuntimeCall(t *testing.T) {
	for skip := 0; ; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		fmt.Printf("skip = %v, pc = %v, file = %v, line = %v\n", skip, pc, file, line)
	}
}

func TestMutexCall(t *testing.T) {
	mis := &mutexStruct{}

	mis.mutex.Lock()
	defer mis.mutex.Unlock()

	// mis.mutex.Lock()
	t.Log("do something")
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
