package external

import (
	"crypto/tls"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TrimURL removes the protocol version and trailing slashes.
func TrimURL(url string) string {
	url = strings.TrimSuffix(url, "/")
	parts := strings.Split(url, "://")
	if len(parts) < 2 {
		return url
	}
	return strings.Join(parts[1:], "://")
}

// BaseURL prepends `protocol + "://"` or `protocol + "s://"` to the url depending on TLS support.
func BaseURL(protocol string, tls bool, trimmedURL string, a ...any) string {
	trimmedURL = fmt.Sprintf(trimmedURL, a...)
	if tls {
		return protocol + "s://" + trimmedURL
	} else {
		return protocol + "://" + trimmedURL
	}
}

// IsTLS verifies the TLS certificate of a trimmed URL.
func IsTLS(trimmedURL string) bool {
	url, err := url.Parse("https://" + trimmedURL)
	if err != nil {
		return false
	}
	host := url.Host
	if url.Port() == "" {
		host = host + ":443"
	}

	conn, err := tls.Dial("tcp", url.Host, &tls.Config{})
	if err != nil {
		return false
	}
	defer conn.Close()

	err = conn.VerifyHostname(url.Hostname())
	if err != nil {
		return false
	}

	expiry := conn.ConnectionState().PeerCertificates[0].NotAfter
	if time.Now().After(expiry) {
		return false
	}

	return true
}

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
