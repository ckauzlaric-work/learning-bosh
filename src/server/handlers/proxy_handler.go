package handlers

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

type ProxyHandler struct {
}

var httpClient = &http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: true,
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 0,
		}).Dial,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

func (h *ProxyHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	protocol := req.URL.Query().Get("protocol")
	if protocol == "" {
		protocol = "http"
	}
	destination := fmt.Sprintf("%s://%s", protocol, strings.TrimPrefix(req.URL.Path, "/proxy/"))
	log.Printf("received proxy request for: %s\n", destination)

	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("Request Info: \n%s\n", string(requestDump))

	getResp, err := httpClient.Get(destination)
	if err != nil {
		fmt.Fprintf(os.Stderr, "request failed: %s", err)
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(fmt.Sprintf("request failed: %s", err)))
		return
	}

	log.Printf("made request for destination: %s\n", destination)

	if getResp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("request returned non-ok status code: %d\n", getResp.StatusCode)
		fmt.Fprint(os.Stderr, msg)
		resp.WriteHeader(getResp.StatusCode)
		resp.Write([]byte(msg))
		return
	}

	log.Println("status code was okay, reading body...")

	defer getResp.Body.Close()

	readBytes, err := io.ReadAll(getResp.Body)
	if err != nil {
		msg := fmt.Sprintf("read body failed: %s", err)
		fmt.Fprint(os.Stderr, msg)
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(msg))
		return
	}

	log.Println("successfully read body")

	resp.Write(readBytes)
}
