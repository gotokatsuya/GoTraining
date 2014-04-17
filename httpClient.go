package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Client struct {
	c       *http.Client
	baseUrl string
}

func (c Client) Request() []byte {
	res, err := c.c.Get(c.baseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return nil
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return nil
	}
	res.Body.Close()
	return body
}

func safeCounter() {
	count := make(chan int, 1)
	go func() {
		for i := 1; ; i++ { // sahred counter
			count <- i
		}
	}()
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		time.Sleep(100 * time.Millisecond) // do domething
		fmt.Fprintf(res, "Hello. You are #%d visitor.", <-count)
	})
}

func unsafeCounter() {
	var i = 0 // shared counter
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		i++
		time.Sleep(100 * time.Millisecond) // do domething
		fmt.Fprintf(res, "Hello. You are #%d visitor.", i)
	})
}

func runServer(safe bool) *Client {
	const addr = "127.0.0.1:3000"
	go func() {
		if safe {
			safeCounter()
		} else {
			unsafeCounter()
		}
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return &Client{new(http.Client), "http://" + addr}
}

func main() {
	c := runServer(true) // safe counter
	// c := runServer(false) // unsafe counter

	out := make(chan []byte)
	fin := make(chan bool)

	go func() {
		var w sync.WaitGroup
		for i := 0; i < 5; i++ {
			w.Add(1)
			go func() {
				out <- c.Request()
				w.Done()
			}()
		}
		w.Wait()
		fin <- true
	}()

LOOP:
	for {
		select {
		case res := <-out:
			fmt.Println(string(res))
		case <-fin:
			break LOOP
		}
	}
}
