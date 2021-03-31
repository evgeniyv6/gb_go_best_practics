package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("hello from beat 1")
	sigintChan:= make(chan os.Signal)
	signal.Notify( sigintChan, syscall.SIGINT)

	sigtermChan:= make(chan os.Signal)
	signal.Notify( sigtermChan, syscall.SIGTERM)

	for {
		select {
		case <-sigintChan:
			log.Println("got sigint")

		case <-sigtermChan:
			log.Println("got sigtrem")

		default:
			log.Println("still running")
			time.Sleep(2* time.Second)
		}
	}

	//sig:= <-sigChan
	//fmt.Printf(" ==> got signal %s\n", sig.String())
	//fmt.Print("end")
}
