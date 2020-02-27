package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

func doWork(url string) string {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

func main() {
	var total int
	var goroutineCounter int64
	mutex := new(sync.Mutex)
	works := make(chan struct{}, 4)
	wg := new(sync.WaitGroup)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		url := scanner.Text()
		fmt.Println("Start. Active", atomic.AddInt64(&goroutineCounter, 1))
		works <- struct{}{}
		wg.Add(1)
		go func() {
			count := strings.Count(doWork(url), "Go")
			mutex.Lock()
			total += count
			mutex.Unlock()
			fmt.Println("Count for", url, ":", count)
			<-works
			wg.Done()
			fmt.Println("Stop. Active", atomic.AddInt64(&goroutineCounter, -1))
		}()
	}
	close(works)
	wg.Wait()
	fmt.Println("Total:", total)
}
