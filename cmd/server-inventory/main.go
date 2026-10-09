package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func readPairs(path, separator string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	values := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		key, value, found := strings.Cut(line, separator)
		if found {
			values[strings.TrimSpace(key)] = strings.Trim(strings.TrimSpace(value), `"`)
		}
	}
	return values, nil
}

func run() int {
	if runtime.GOOS != "linux" {
		fmt.Fprintln(os.Stderr, "Requires Linux /proc and /etc/os-release")
		return 1
	}
	host, err := os.Hostname()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	system, err := readPairs("/etc/os-release", "=")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	cpu, err := readPairs("/proc/cpuinfo", ":")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	memory, err := readPairs("/proc/meminfo", ":")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fields := strings.Fields(memory["MemTotal"])
	if len(fields) != 2 || fields[1] != "kB" {
		fmt.Fprintln(os.Stderr, "Missing or invalid MemTotal")
		return 1
	}
	kib, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	model := cpu["model name"]
	if model == "" {
		model = cpu["Hardware"]
	}
	if model == "" {
		model = "unknown"
	}
	fmt.Printf("Hostname: %s\nOS: %s (%s)\nArchitecture: %s\nLogical CPUs: %d\nCPU model: %s\nRAM: %.2f GiB\n", host, system["PRETTY_NAME"], runtime.GOOS, runtime.GOARCH, runtime.NumCPU(), model, float64(kib)/(1024*1024))
	return 0
}
func main() { os.Exit(run()) }
