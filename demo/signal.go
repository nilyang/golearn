package main

import (
	"os/signal"
	"syscall"
	"fmt"
	"os"
	"go/src/runtime"
)

func main() {
	demo()
}

// 系统信号demo
func demo() {
	var forever chan bool
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT)
	fmt.Println("Num of Cpus: ", runtime.NumCPU())
	go func() {
		for {
			for {
				<-s
				fmt.Println("[TEST_SIGINT] recieved! ")
				goto hello
			}
		}

		hello:
			fmt.Println("Game over")
	}()

	<-forever

}


