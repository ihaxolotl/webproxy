package proxy

import (
	"bufio"
	"bytes"
	"net/http"
	"strings"

	"github.com/ihaxolotl/webproxy/internal/buffer"
)

type filter struct {
	Find    string
	Replace string
}

// parseProxyRequest parses an HTTP request crafted for a proxy and creates a new request
// that can be processed by the target web server.
func parseProxyRequest(src *buffer.Buffer, req *http.Request) (*buffer.Buffer, error) {
	var (
		dst  []byte
		n    int
		diff int
	)

	// HACK: Replace the proxy headers in the request.
	filters := []filter{
		{Find: "Proxy-Connection:", Replace: "Connection:"},
		{Find: req.URL.Scheme + "://" + req.URL.Host, Replace: ""},
	}

	dst = make([]byte, buffer.ReadBufferSize)
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

	return buffer.NewBufferFrom(dst, n), nil
}

// readRequest parses an http.Request object from a byte slice.
func readRequest(buf *buffer.Buffer) *http.Request {
	// HACK: Parse the the request to get the hostname.
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(buf.Buffer())))
	if err != nil {
		panic(err)
	}

	return req
}

// readResponse parses an http.Response object from a byte slice.
func readResponse(req *http.Request, buf *buffer.Buffer) *http.Response {
	res, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(buf.Buffer())), req)
	if err != nil {
		panic(err)
	}

	return res
}
