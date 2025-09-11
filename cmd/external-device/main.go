package main

import (
	"flag"
	"net"
	"net/http"

	"github.com/Khaym03/kumo/internal/adapters/remote"
	"github.com/go-rod/rod/lib/utils"
)

func main() {
	flag.Parse()

	m := remote.NewMultiManager()

	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		utils.E(err)
	}

	srv := &http.Server{Handler: m}
	utils.E(srv.Serve(l))
}
