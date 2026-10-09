package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type event struct {
	Time    string `json:"time"`
	URL     string `json:"url"`
	Healthy bool   `json:"healthy"`
	Detail  string `json:"detail"`
}

func observe(ctx context.Context, client *http.Client, address string) event {
	item := event{Time: time.Now().Format(time.RFC3339), URL: address}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, address, nil)
	if err != nil {
		item.Detail = err.Error()
		return item
	}
	response, err := client.Do(request)
	if err != nil {
		item.Detail = err.Error()
		return item
	}
	defer response.Body.Close()
	item.Healthy = response.StatusCode >= 200 && response.StatusCode < 300
	item.Detail = response.Status
	return item
}

func monitor(ctx context.Context, addresses []string, interval time.Duration, encoder *json.Encoder) error {
	client := &http.Client{Timeout: 5 * time.Second}
	previous := make(map[string]bool)
	for {
		for _, address := range addresses {
			if ctx.Err() != nil {
				return nil
			}
			item := observe(ctx, client, address)
			if ctx.Err() != nil {
				return nil
			}
			last, seen := previous[address]
			// Emit the first result and subsequent state changes only.
			if !seen || last != item.Healthy {
				if err := encoder.Encode(item); err != nil {
					return err
				}
			}
			previous[address] = item.Healthy
		}
		// Wait after a complete sweep, so slow requests do not create a busy loop.
		timer := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil
		case <-timer.C:
		}
	}
}

func run() int {
	interval := flag.Duration("interval", 30*time.Second, "Pause between complete sweeps")
	flag.Parse()
	if *interval <= 0 || len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "Provide URL arguments and a positive interval")
		return 2
	}
	for _, address := range flag.Args() {
		parsed, err := url.Parse(address)
		if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
			fmt.Fprintln(os.Stderr, "Invalid URL:", address)
			return 2
		}
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := monitor(ctx, flag.Args(), *interval, json.NewEncoder(os.Stdout)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
func main() { os.Exit(run()) }
