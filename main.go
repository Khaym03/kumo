package main

import (
	"time"
)

func main() {
	k := NewKumo()
	defer k.Shutdown()
	time.Sleep(5 * time.Second)
}
