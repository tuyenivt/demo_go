package app

import (
	"fmt"
	"sync"
)

func SimpleSyncPool() {
	// sync.Pool short term only, GC collect anytime, don't use on important resources (file descriptor, DB connection, ...)
	bufferPool := sync.Pool{
		New: func() interface{} {
			return make([]byte, 1024)
		},
	}

	buf1 := bufferPool.Get().([]byte)
	copy(buf1, []byte("Hello, sync.Pool!"))
	fmt.Println("buf1 content:", string(buf1[:18]))

	bufferPool.Put(buf1)

	buf2 := bufferPool.Get().([]byte)
	fmt.Println("buf2 content after reuse:", string(buf2[:18]))

	bufferPool.Put(buf2)
}
