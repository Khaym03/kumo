package core

import (
	"context"
	"log"

	"github.com/Khaym03/kumo/collectors"
)

type KumoHTTP struct {
	ctx context.Context
	*collectors.CollectorRegistry
}

func NewKumoHTTP() *KumoHTTP {
	return &KumoHTTP{
		ctx:               context.Background(),
		CollectorRegistry: collectors.NewCollectorRegistry(),
	}
}

func (k *KumoHTTP) Run() {
	for _, c := range k.Collectors {
		err := c.Collect(k.ctx)
		if err != nil {
			log.Println(err)
		}
	}
}

func (k KumoHTTP) Shutdown() {
}
