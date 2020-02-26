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
	works := make(chan struct{})
	wg := new(sync.WaitGroup)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		url := scanner.Text()
		wg.Add(1)
		go func() {
			for range works {
				count := strings.Count(doWork(url), "Go")
				total += count
				fmt.Println("Count for", url, ":", count)
			}
			wg.Done()
		}()
		works <- struct{}{}
	}
	close(works)
	wg.Wait()
	fmt.Println("Total:", total)
}
