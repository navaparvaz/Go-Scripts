//go:build !linux

package main

import (
	"fmt"
	"os"
)

func main() { fmt.Fprintln(os.Stderr, "disk-monitor requires Linux"); os.Exit(1) }
