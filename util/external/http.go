package external

import (
	"mime"
	"net/http"
	"strings"
)

// HasContentType returns true if the 'content-type' header includes mimetype.
func HasContentType(h http.Header, mimetype string) bool {
	contentType := h.Get("content-type")
	if contentType == "" {
		return mimetype == "application/octet-stream"
	}

	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}
		if t == mimetype {
			return true
		}
	}
	return false
}
