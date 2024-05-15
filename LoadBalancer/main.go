package main

import (
	"fmt"
	"net/http"
	"time"
)

func worker(workerId int, data chan int) {
	for x := range data {
		fmt.Printf("Worker %d received %d", workerId, x)
		time.Sleep(time.Second)
	}
}

func main() {
	canal := make(chan int)
	qtdWorkers := 50
	for i := 0; i < qtdWorkers; i++ {
		go worker(i, canal)
	}

	for i := 0; i < 100000; i++ {
		canal <- i
	}
	http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil)
}
