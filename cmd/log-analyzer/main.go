package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

type summary struct {
	Lines   int
	Matches int
}

func analyze(path string, pattern *regexp.Regexp) (summary, error) {
	counts := summary{}
	file, err := os.Open(path)
	if err != nil {
		return counts, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for scanner.Scan() {
		counts.Lines++
		if pattern.MatchString(scanner.Text()) {
			counts.Matches++
			fmt.Printf("%s:%d: %s\n", path, counts.Lines, scanner.Text())
		}
	}
	return counts, scanner.Err()
}

func run() int {
	expression := flag.String("pattern", `(?i)\b(error|fatal|panic)\b`, "Regular expression to find")
	flag.Parse()
	if len(flag.Args()) == 0 || *expression == "" {
		fmt.Fprintln(os.Stderr, "Provide log file arguments and a nonempty pattern")
		return 2
	}
	pattern, err := regexp.Compile(*expression)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	total := summary{}
	failed := false
	for _, path := range flag.Args() {
		counts, err := analyze(path, pattern)
		if err != nil {
			fmt.Fprintln(os.Stderr, path, err)
			failed = true
		}
		total.Lines += counts.Lines
		total.Matches += counts.Matches
		fmt.Printf("Summary %s: %d matches / %d lines\n", path, counts.Matches, counts.Lines)
	}
	fmt.Printf("Total: %d matches / %d lines\n", total.Matches, total.Lines)
	if failed || total.Matches > 0 {
		return 1
	}
	return 0
}
func main() { os.Exit(run()) }
