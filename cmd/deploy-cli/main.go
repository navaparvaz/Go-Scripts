package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type deployment struct {
	Application string `json:"application"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
}

func send(client *http.Client, address string, payload []byte) error {
	request, err := http.NewRequest(http.MethodPost, address, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	if token := os.Getenv("DEPLOY_TOKEN"); token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("deployment endpoint returned %s", response.Status)
	}
	_, err = io.Copy(io.Discard, io.LimitReader(response.Body, 4096))
	return err
}

func run() int {
	if len(os.Args) < 2 || os.Args[1] != "deploy" {
		fmt.Fprintln(os.Stderr, "Usage: deploy-cli deploy -app NAME -version TAG [-env staging] [-execute -endpoint URL]")
		return 2
	}
	flags := flag.NewFlagSet("deploy", flag.ContinueOnError)
	app := flags.String("app", "", "Application name")
	version := flags.String("version", "", "Version to deploy")
	environment := flags.String("env", "staging", "Deployment environment")
	endpoint := flags.String("endpoint", "", "Deployment webhook URL")
	execute := flags.Bool("execute", false, "Send the request instead of previewing it")
	if err := flags.Parse(os.Args[2:]); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if *app == "" || *version == "" || *environment == "" || len(flags.Args()) > 0 {
		fmt.Fprintln(os.Stderr, "Application, version and environment are required")
		return 2
	}
	payload, err := json.MarshalIndent(deployment{*app, *version, *environment}, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if !*execute {
		fmt.Println("Preview only:")
		fmt.Println(string(payload))
		return 0
	}
	parsed, err := url.Parse(*endpoint)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		fmt.Fprintln(os.Stderr, "Provide an HTTP(S) endpoint")
		return 2
	}
	// Refuse redirects so a deployment payload or bearer token is not forwarded.
	client := &http.Client{Timeout: 15 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	if err := send(client, *endpoint, payload); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Println("Deployment request accepted; the endpoint controls rollout.")
	return 0
}
func main() { os.Exit(run()) }
