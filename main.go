package main

import "github.com/Khaym03/kumo/core"

func main() {
	kumo := core.NewKumo()

	defer kumo.Shutdown()

	kumo.Run()
}
