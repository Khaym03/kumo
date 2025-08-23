package controller

import "github.com/Khaym03/kumo/ports"

type collectorsRegister struct {
	collectors []ports.Collector
}

func NewCollectorRegistry() ports.CollectorRegistry {
	return &collectorsRegister{}
}

func (c *collectorsRegister) RegisterCollector(newc ports.Collector) {
	c.collectors = append(c.collectors, newc)
}

func (c *collectorsRegister) Collectors() []ports.Collector {
	return c.collectors
}
