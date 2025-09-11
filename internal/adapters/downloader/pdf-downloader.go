package downloader

import (
	"log"
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type pageBasedPDFDownloader struct {
	page *rod.Page
}

func NewPageBasedPDFDownloader(p *rod.Page) *pageBasedPDFDownloader {
	return &pageBasedPDFDownloader{page: p}
}

func (d *pageBasedPDFDownloader) Download(url, savePath string) error {
	router := d.page.HijackRequests()
	done := make(chan error, 1)

	router.MustAdd(url, func(ctx *rod.Hijack) {
		ctx.ContinueRequest(&proto.FetchContinueRequest{})
		ctx.MustLoadResponse()
		res := ctx.Response.Body()
		err := os.WriteFile(savePath, []byte(res), 0644)
		if err != nil {
			done <- err
			return
		}
		log.Printf("Successfully downloaded %s", savePath)
		done <- nil
	})

	go router.Run()

	if err := d.page.Navigate(url); err != nil {
		return err
	}

	router.Remove(url)

	return <-done
}
