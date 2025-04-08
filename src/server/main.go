package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"

	"server/handlers"
)

const (
	defaultName    = "proxy"
	ConfigFilename = "server.yml"
)

type Config struct {
	Port int `yaml:"port"`
}

func GetConfig() (*Config, error) {
	c := &Config{}
	yamlFile, err := os.ReadFile(ConfigFilename)
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c, nil
}

func main() {
	c, err := GetConfig()
	if err != nil {
		log.Printf("warning: failed to fetch config %#v", err)
		panic(err)
	}

	name, ok := os.LookupEnv("NAME")
	if !ok {
		log.Printf("warning: NAME env var not provided, using default name: %s", defaultName)
		name = defaultName
	}

	mux := http.NewServeMux()
	mux.Handle("/", &handlers.RequestInfoHandler{Name: name})
	mux.Handle("/proxy/", &handlers.ProxyHandler{})

	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", c.Port), mux)
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
