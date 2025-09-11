package ports

import "github.com/Khaym03/kumo/internal/pkg/proxy"

type ProxiesDownloader interface {
	Download() ([]proxy.Proxy, error)
}
