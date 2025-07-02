package utils

import (
	"crypto/tls"
	log "log/slog"
	"net/http"
	"time"
)

// HTTPSClient implements HTTPClient interface
type HTTPSClient struct {
	client *http.Client
}

func NewHTTPSClient(skipTLSVerify bool, timeout time.Duration) *HTTPSClient {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipTLSVerify,
		},
	}

	return &HTTPSClient{
		client: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
	}
}

func (h *HTTPSClient) Do(req *http.Request) (*http.Response, error) {
	log.Debug("Executing HTTPS request", "url", req.URL.String(), "method", req.Method)
	return h.client.Do(req)
}
