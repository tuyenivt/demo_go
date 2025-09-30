package app

import "fmt"

func ReadWriteChannel(name string) {
	ch := make(chan string)
	defer close(ch)
	go writeOnlyChannel(name, ch)
	readOnlyChannel(ch)
}

func readOnlyChannel(ch <-chan string) {
	name := <-ch
	fmt.Printf("Read Write Channel, %s\n", name)
}

func writeOnlyChannel(name string, ch chan<- string) {
	ch <- name
}
