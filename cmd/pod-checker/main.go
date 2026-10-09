package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type pod struct {
	Metadata struct {
		Name      string
		Namespace string
	}
	Status struct {
		Phase      string
		Conditions []struct {
			Type   string
			Status string
		}
		ContainerStatuses []struct{ RestartCount int }
	}
}
type podPage struct {
	Items    []pod
	Metadata struct{ Continue string }
}

func listPods(client *http.Client, endpoint string) ([]pod, error) {
	var all []pod
	address, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	// Kubernetes may split a large list into pages.
	for {
		response, err := client.Get(address.String())
		if err != nil {
			return nil, err
		}
		if response.StatusCode != http.StatusOK {
			response.Body.Close()
			return nil, fmt.Errorf("Kubernetes returned %s", response.Status)
		}
		var page podPage
		err = json.NewDecoder(response.Body).Decode(&page)
		response.Body.Close()
		if err != nil {
			return nil, err
		}
		all = append(all, page.Items...)
		if page.Metadata.Continue == "" {
			return all, nil
		}
		query := address.Query()
		query.Set("continue", page.Metadata.Continue)
		address.RawQuery = query.Encode()
	}
}

func ready(item pod) bool {
	if item.Status.Phase == "Succeeded" {
		return true
	}
	if item.Status.Phase != "Running" {
		return false
	}
	for _, condition := range item.Status.Conditions {
		if condition.Type == "Ready" {
			return condition.Status == "True"
		}
	}
	return false
}

func run() int {
	api := flag.String("api", "http://127.0.0.1:8001", "kubectl proxy URL")
	namespace := flag.String("namespace", "default", "Namespace to inspect")
	maxRestarts := flag.Int("max-restarts", 5, "Maximum total restarts per pod")
	flag.Parse()
	base, err := url.Parse(*api)
	if err != nil || base.Host == "" || (base.Scheme != "http" && base.Scheme != "https") || *namespace == "" || strings.ContainsAny(*namespace, "/ ?#") || *maxRestarts < 0 {
		fmt.Fprintln(os.Stderr, "Invalid API URL, namespace or restart limit")
		return 2
	}
	endpoint := strings.TrimRight(*api, "/") + "/api/v1/namespaces/" + url.PathEscape(*namespace) + "/pods?limit=100"
	items, err := listPods(&http.Client{Timeout: 10 * time.Second}, endpoint)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	failed := false
	for _, item := range items {
		restarts := 0
		for _, state := range item.Status.ContainerStatuses {
			restarts += state.RestartCount
		}
		okay := ready(item) && restarts <= *maxRestarts
		fmt.Printf("%s/%s phase=%s restarts=%d healthy=%t\n", item.Metadata.Namespace, item.Metadata.Name, item.Status.Phase, restarts, okay)
		if !okay {
			failed = true
		}
	}
	fmt.Printf("Checked %d pods\n", len(items))
	if failed {
		return 1
	}
	return 0
}
func main() { os.Exit(run()) }
