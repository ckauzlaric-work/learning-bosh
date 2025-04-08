package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

type RequestInfoHandler struct {
	Name string
}

func (h *RequestInfoHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log.Printf("Endpoint (%s) received a request.", h.Name)

	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("Request Info: \n\n%s\n", string(requestDump))

	fmt.Fprintf(resp, "(%s) Request Info: \n\n%s\n", h.Name, string(requestDump))

}
