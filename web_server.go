// Copyright 2020 Kuei-chun Chen. All rights reserved.

package gox

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// StartWebServer starts a web server
func StartWebServer(port int) error {
	var err error
	var hostname string
	if hostname, err = os.Hostname(); err != nil {
		hostname = "127.0.0.1"
	}
	log.Printf("HTTP server ready, URL: http://%s:%d/\n", hostname, port)
	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(addr, nil)
}
