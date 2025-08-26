package ports

import "github.com/Khaym03/kumo/proxy"

type RemoteBrowserURLBuilder interface {
	WithProxy(proxy.Proxy)
	Build() (string, error)
}
