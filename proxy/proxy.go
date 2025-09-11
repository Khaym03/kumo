package proxy

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Proxy represents a single proxy server with its credentials.
type Proxy struct {
	Host     string
	Port     string
	User     string
	Password string
}

func (p Proxy) Address() string {
	return fmt.Sprintf("%s:%s", p.Host, p.Port)
}

func (p Proxy) RequireAuth() bool {
	return p.User != "" && p.Password != ""
}

const (
	proxyURLScheme = "http"
)

// NewClient creates a new Client configured to use the given proxy.
func NewClient(p Proxy) *http.Client {
	proxyURL := &url.URL{
		Scheme: proxyURLScheme,
		User:   url.UserPassword(p.User, p.Password),
		Host:   fmt.Sprintf("%s:%s", p.Host, p.Port),
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	return httpClient
}

// webshareProxyProvider is a concrete implementation of ProxiesDownloader.
//
// When the server's rate limit is exceeded (HTTP 429), it's necessary to wait
// for 60 seconds before making a new request to prevent errors.
type webshareProxyProvider struct{}

func NewWebshareProxyProvider() *webshareProxyProvider {
	return &webshareProxyProvider{}
}

func (wpp *webshareProxyProvider) Download() ([]Proxy, error) {
	ProxyListAPIURL := os.Getenv("PROXY_LIST_API_URL")
	if ProxyListAPIURL == "" {
		return nil, errors.New("PROXY_LIST_API_URL environment variable is not set")
	}

	resp, err := http.Get(ProxyListAPIURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download proxy list: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	proxies := wpp.parse(body)

	return proxies, nil
}

func (wpp *webshareProxyProvider) parse(body []byte) []Proxy {
	var proxies []Proxy

	lines := strings.SplitSeq(strings.TrimSpace(string(body)), "\n")
	for line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) != 4 {
			log.Printf("Invalid proxy format, skipping: %s", line)
			continue
		}
		p := Proxy{
			Host:     parts[0],
			Port:     parts[1],
			User:     parts[2],
			Password: parts[3],
		}
		proxies = append(proxies, p)
	}

	return proxies
}
