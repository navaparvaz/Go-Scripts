package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type container struct {
	Name  string
	State struct {
		Running bool
		Status  string
		Health  *struct{ Status string }
	}
}

func inspect(ctx context.Context, names []string) ([]container, error) {
	args := append([]string{"container", "inspect", "--"}, names...)
	output, err := exec.CommandContext(ctx, "docker", args...).Output()
	if err != nil {
		return nil, fmt.Errorf("docker inspect: %w", err)
	}
	var items []container
	err = json.Unmarshal(output, &items)
	return items, err
}

func run() int {
	timeout := flag.Duration("timeout", 10*time.Second, "Command timeout")
	flag.Parse()
	if *timeout <= 0 {
		fmt.Fprintln(os.Stderr, "Timeout must be positive")
		return 2
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	names := flag.Args()
	// With no names supplied, discover both running and stopped containers.
	if len(names) == 0 {
		output, err := exec.CommandContext(ctx, "docker", "container", "ls", "--all", "--quiet").Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		names = strings.Fields(string(output))
	}
	if len(names) == 0 {
		fmt.Println("No containers found")
		return 0
	}
	items, err := inspect(ctx, names)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if len(items) != len(names) {
		fmt.Fprintln(os.Stderr, "Docker returned an incomplete container list")
		return 1
	}
	failures := 0
	for _, item := range items {
		health := "not configured"
		okay := item.State.Running
		if item.State.Health != nil {
			health = item.State.Health.Status
			okay = okay && health == "healthy"
		}
		fmt.Printf("%s running=%t state=%s health=%s\n", strings.TrimPrefix(item.Name, "/"), item.State.Running, item.State.Status, health)
		if !okay {
			failures++
		}
	}
	if failures > 0 {
		return 1
	}
	return 0
}
func main() { os.Exit(run()) }
