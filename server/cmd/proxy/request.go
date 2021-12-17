package main

import (
	"bufio"
	"bytes"
	"net/http"
	"strings"
)

type filter struct {
	Find    string
	Replace string
}

// parseProxyRequest parses an HTTP request crafted for a proxy and creates a new request
// that can be processed by the target web server.
func parseProxyRequest(src *Buffer, filters []filter) (*Buffer, error) {
	var (
		dst  []byte
		n    int
		diff int
	)

	dst = make([]byte, ReadBufferSize)
	n = src.Size()

	copy(dst, src.Buffer())

	for _, v := range filters {
		if strings.Contains(string(dst), v.Find) {
			dst = []byte(strings.Replace(string(dst), v.Find, v.Replace, 1))

			diff = len(v.Find) - len(v.Replace)
			if n >= diff {
				n -= diff
			}
		}
	}

	return NewBufferFrom(dst, n), nil
}

// readRequest is a hack for parsing the hostname and URL object from a
// byte slice.
func readRequest(buf []byte, n int) *http.Request {
	// HACK: Parse the the request to get the hostname.
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(buf[:n])))
	if err != nil {
		panic(err)
	}

	return req
}
