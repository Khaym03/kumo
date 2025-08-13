package main

import "github.com/Khaym03/kumo/core"

func main() {
	var kumo core.Kumo = core.NewKumoHTTP()

	defer kumo.Shutdown()

	kumo.Run()
}
