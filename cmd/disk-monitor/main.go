//go:build linux

package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"syscall"
)

func availableGiB(path string) (float64, error) {
	var stats syscall.Statfs_t
	if err := syscall.Statfs(path, &stats); err != nil {
		return 0, err
	}
	// Bavail excludes blocks reserved for the administrator.
	bytes := float64(stats.Bavail) * float64(stats.Bsize)
	return bytes / (1024 * 1024 * 1024), nil
}

func run() int {
	minimum := flag.Float64("min-free", 2, "Minimum available GiB")
	flag.Parse()
	if *minimum < 0 || math.IsNaN(*minimum) || math.IsInf(*minimum, 0) {
		fmt.Fprintln(os.Stderr, "Minimum must be a finite nonnegative number")
		return 2
	}
	paths := flag.Args()
	if len(paths) == 0 {
		paths = []string{"/"}
	}
	failed := false
	for _, path := range paths {
		free, err := availableGiB(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, path, err)
			failed = true
			continue
		}
		fmt.Printf("%s: %.2f GiB available\n", path, free)
		if free < *minimum {
			fmt.Printf("ALERT %s: below %.2f GiB\n", path, *minimum)
			failed = true
		}
	}
	if failed {
		return 1
	}
	return 0
}
func main() { os.Exit(run()) }
