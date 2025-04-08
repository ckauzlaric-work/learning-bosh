package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"proxy/handlers"
)

const (
	defaultName = "proxy"
)

func main() {
	port := getEnvVar("PORT", 8080, false)
	name, ok := os.LookupEnv("NAME")
	if !ok {
		log.Printf("warning: NAME env var not provided, using default name: %s", defaultName)
		name = defaultName
	}

	mux := http.NewServeMux()
	mux.Handle("/", &handlers.RequestInfoHandler{Name: name})
	mux.Handle("/proxy/", &handlers.ProxyHandler{})

	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), mux)
}

func getEnvVar(key string, defaultValue int, failIfDNE bool) int {
	var result int
	var err error

	v, ok := os.LookupEnv(key)
	if !ok && failIfDNE {
		log.Fatalf("invalid required env var %s", key)
	} else if !ok {
		return defaultValue
	}

	result, err = strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return result
}
