package main

import (
	"github.com/Khaym03/kumo/collectors"
)

func main() {
	k := NewKumo()
	defer k.Shutdown()

	k.RegisterCollector(collectors.NewCategoriesCollector(k.scheduler))

	k.Run()
}
