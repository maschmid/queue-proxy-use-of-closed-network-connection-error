package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	length = flag.Int("length", 32 * 1024, "Expected length of the requests' data")
)

func rootHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "This receiver only accepts POST", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// return 418 here, to distinguish between other causes of HTTP 500..."
		http.Error(w, fmt.Sprintf("error reading body: %v", err), http.StatusTeapot)
		return
	}

	if len(body) != *length {
		http.Error(w, fmt.Sprintf("unexpected body length: %d != %d", len(body), *length), http.StatusBadRequest)
		return
	}

	for i, c := range body {
		if c != 42 {
			http.Error(w, fmt.Sprintf("unexpected byte at index %d (%d != 42)", i, c), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(200)
	w.Write(body)
}

func main() {

	flag.Parse()

	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
