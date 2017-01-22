package main

import (
	"fmt"
	"time"
)


var chan1 chan int
var chanLength int = 18
var interval time.Duration = 1500 * time.Millisecond
func main() {
	fmt.Println("vim-go")
	chan1 = make(chan int, chanLength)
	go func() {
		for i := 0; i < chanLength; i++ {
			if i > 0 && i % 3 == 0 {
				fmt.Println("Reset chan1...")
				chan1 = make(chan int, chanLength)
			}
			fmt.Printf("Send element %d...\n", i)
			chan1 <- i
			time.Sleep(interval)
		}
		fmt.Println("Close chan1...")
		close(chan1)
	}()
	receive(chan1)
//	time.Sleep(time.Second)
}

func receive(chan2 chan int) {
	fmt.Println("Receive element from chan1...")
	timer := time.After(30 * time.Second)
Loop:
	for {
		select {
			case e, ok := <-getChan():
				if !ok {
					fmt.Println("--Closed chan1.")
					break Loop
				}
				fmt.Printf("Receive a element: %d\n", e)
				time.Sleep(interval)
			case <-timer:
				fmt.Println("Timeout!")
				break Loop
			default:
//				fmt.Println("------haha-------")
		}
	}
	fmt.Println("--End")
}

func getChan() chan int {
	return chan1
}
