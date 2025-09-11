package ports

import "github.com/Khaym03/kumo/internal/pkg/proxy"

type RemoteBrowserURLBuilder interface {
	WithProxy(proxy.Proxy)
	Build() (string, error)
}
