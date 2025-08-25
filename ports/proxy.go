package ports

import "github.com/Khaym03/kumo/proxy"

type ProxiesDownloader interface {
	Download() ([]proxy.Proxy, error)
}
