package ports

import "github.com/Khaym03/kumo/pkg/proxy"

type RemoteBrowserURLBuilder interface {
	WithProxy(proxy.Proxy)
	Build() (string, error)
}
