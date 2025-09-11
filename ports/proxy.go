package ports

import "github.com/Khaym03/kumo/pkg/proxy"

type ProxiesDownloader interface {
	Download() ([]proxy.Proxy, error)
}
