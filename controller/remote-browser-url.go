package controller

import (
	"fmt"
	"net/url"

	"github.com/Khaym03/kumo/pkg/proxy"
	"github.com/google/uuid"
)

type remoteBrowserURLBuilder struct {
	remoteHost string
	proxy      proxy.Proxy
}

// e.g 127.0.0.1:1234
func NewWSURLBuilder(remoteHost string) *remoteBrowserURLBuilder {
	return &remoteBrowserURLBuilder{
		remoteHost: remoteHost,
	}
}

func (b *remoteBrowserURLBuilder) WithProxy(p proxy.Proxy) *remoteBrowserURLBuilder {
	b.proxy = p
	return b
}

func (b *remoteBrowserURLBuilder) Build() (string, error) {
	values := url.Values{}
	values.Set("id", "instance-"+uuid.New().String())

	if b.proxy.Host != "" && b.proxy.Port != "" {
		values.Set("proxyHost", b.proxy.Host)
		values.Set("proxyPort", b.proxy.Port)
		if b.proxy.User != "" {
			values.Set("proxyUser", b.proxy.User)
		}
		if b.proxy.Password != "" {
			values.Set("proxyPass", b.proxy.Password)
		}
	}

	wsURL := fmt.Sprintf("ws://%s/?%s", b.remoteHost, values.Encode())
	return wsURL, nil
}
