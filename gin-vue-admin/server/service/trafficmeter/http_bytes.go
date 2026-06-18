package trafficmeter

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func estimateRequestBytes(req *http.Request) uint64 {
	if req == nil {
		return 0
	}
	uri := "/"
	if req.URL != nil {
		uri = req.URL.RequestURI()
	}
	proto := req.Proto
	if proto == "" {
		proto = "HTTP/1.1"
	}
	total := len(req.Method) + len(" ") + len(uri) + len(" ") + len(proto) + len("\r\n")
	if req.Host != "" {
		total += len("Host: ") + len(req.Host) + len("\r\n")
	}
	total += estimateHTTPHeaderBytes(req.Header)
	total += len("\r\n")
	if req.ContentLength > 0 {
		total += int(req.ContentLength)
	}
	return uint64(total)
}

func estimateResponseBytes(writer gin.ResponseWriter) uint64 {
	if writer == nil {
		return 0
	}
	statusCode := writer.Status()
	total := len("HTTP/1.1 ") + len(strconv.Itoa(statusCode)) + len(" ") + len(http.StatusText(statusCode)) + len("\r\n")
	header := writer.Header()
	total += estimateHTTPHeaderBytes(header)
	if _, ok := header["Content-Length"]; !ok {
		size := writer.Size()
		if size > 0 {
			total += len("Content-Length: ") + len(strconv.Itoa(size)) + len("\r\n")
		}
	}
	total += len("\r\n")
	if size := writer.Size(); size > 0 {
		total += size
	}
	return uint64(total)
}

func estimateHTTPHeaderBytes(header http.Header) int {
	total := 0
	for key, values := range header {
		for _, value := range values {
			total += len(key) + len(": ") + len(value) + len("\r\n")
		}
	}
	return total
}
