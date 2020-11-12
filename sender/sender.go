package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	// log every 10000th request so we know something is still happening
	// and we have a rough idea about rate of errors
	logEvery = 10000
)

var (
	sink = flag.String("sink", "", "URL to send requests to.")
	// Use 32kB per request, as that is a reasonable size of a cloudevent
	length = flag.Int("length", 32 * 1024, "Expected length of the requests' data")
	parallel = flag.Int("parallel", 32, "Number of parallel senders")

	end = false

	c *http.Client
)

func rootHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("GET /end to stop"))
}

func endHandler(w http.ResponseWriter, req *http.Request) {
	end = true

	log.Printf("sender ending...")

	w.WriteHeader(200)
	w.Write([]byte("stopping!"))
}

func send(sink string) {
	body := make([]byte, *length)

	for i := 0; i < *length; i++ {
		body[i] = 42
	}

	r := bytes.NewReader(body)

	resp, err := c.Post(sink, "text/plain", r)
	if err != nil {
		log.Printf("error on POST: %v", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("status %d not OK!", resp.StatusCode)
		return
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading body: %v", err)
		return
	}

	if len(body) != *length {
		log.Printf("unexpected body length: %d != %d", len(body), *length)
		return
	}

	for i, c := range body {
		if c != 42 {
			log.Printf("unexpected byte at index %d (%d != 42)", i, c)
			return
		}
	}
}

func main() {
	flag.Parse()

	var base = http.DefaultTransport.(*http.Transport).Clone()

	// Same params that mtbroker-filter uses, see
	// https://github.com/knative/eventing/blob/release-0.17/pkg/mtbroker/filter/filter_handler.go#L50-L57
	base.MaxIdleConns = 1000
	base.MaxIdleConnsPerHost = 100
	c = &http.Client{
		Transport: base,
	}

	for i := 0; i < *parallel; i++ {
		i := i
		go func() {
			// Initial sleep to spread parallel senders
			time.Sleep(time.Duration(i) * time.Millisecond)

			var count uint64
			url := *sink
			for end != true {
				// Just log every Nth request
				for j := 0; j < logEvery; j++ {
					send(url)

					count++

					// Limit max. throughput by a little delay between requests
					time.Sleep(10 * time.Millisecond)
				}
				log.Printf("client #%d, requests sent: %d", i, count)
			}
		}()
	}


	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/end", endHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
