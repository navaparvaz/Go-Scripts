package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type result struct {
	URL     string
	Healthy bool
	Detail  string
	Elapsed time.Duration
}

func probe(client *http.Client, address string) result {
	start := time.Now()
	response, err := client.Get(address)
	if err != nil {
		return result{address, false, err.Error(), time.Since(start)}
	}
	defer response.Body.Close()
	return result{address, response.StatusCode >= 200 && response.StatusCode < 300, response.Status, time.Since(start)}
}

func run() int {
	timeout := flag.Duration("timeout", 5*time.Second, "Timeout per service")
	flag.Parse()
	addresses := flag.Args()
	if len(addresses) == 0 || len(addresses) > 50 || *timeout <= 0 {
		fmt.Fprintln(os.Stderr, "Use 1 to 50 URL arguments and a positive timeout")
		return 2
	}
	client := &http.Client{Timeout: *timeout}
	results := make(chan result, len(addresses))
	// Every goroutine sends exactly one result back to the main goroutine.
	for _, address := range addresses {
		go func(url string) { results <- probe(client, url) }(address)
	}
	failed := 0
	for range addresses {
		item := <-results
		state := "OK"
		if !item.Healthy {
			state = "FAIL"
			failed++
		}
		fmt.Printf("%-4s %s | %s | %s\n", state, strings.TrimSpace(item.URL), item.Detail, item.Elapsed.Round(time.Millisecond))
	}
	fmt.Printf("%d/%d services healthy\n", len(addresses)-failed, len(addresses))
	if failed > 0 {
		return 1
	}
	return 0
}
func main() { os.Exit(run()) }
