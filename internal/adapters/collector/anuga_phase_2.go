package adapters

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Khaym03/kumo/internal/pkg/types"
	"github.com/Khaym03/kumo/internal/ports"
	"github.com/go-rod/rod"
	log "github.com/sirupsen/logrus"
)

const containerSelector = `//*[@id="top"]/div[5]/div[2]/section`

type AnugaPhase2Collector struct{}

func (a AnugaPhase2Collector) String() string {
	return phase2Identifier
}

func (a *AnugaPhase2Collector) ProcessPage(
	ctx context.Context,
	page *rod.Page,
	req *types.Request,
	queue ports.Enqueuer,
	fs ports.FileStorage,
) error {
	infoInBytes, err := fs.GetJSON(req.URL)
	if err != nil {
		return err
	}

	var info Phase1Info

	if err = json.Unmarshal(infoInBytes, &info); err != nil {
		return err
	}

	targetURL := info.URL

	const maxRetries = 1
	for retries := 0; retries < maxRetries; retries++ {
		// Create a new context with a specific timeout for navigation.
		navCtx, cancel := context.WithTimeout(ctx, time.Second*30)
		defer cancel()

		log.Infof("Navegando a %s", targetURL)
		// 1. Navigate with context.
		err := page.Navigate(targetURL)
		if err != nil {
			log.Warnf("Error en la navegación (intento %d/%d): %v", retries+1, maxRetries, err)
			time.Sleep(2 * time.Second)
			continue
		}

		title := page.MustInfo().URL

		// 2. Wait for the main content container to be visible. This is more reliable than WaitLoad.
		// The selector is the same one you use for extraction.

		_, err = page.Context(navCtx).Search(containerSelector)
		if err == nil {
			// ¡Éxito! Sal del bucle de reintentos.
			break
		}

		log.Warnf("WaitUntilExists ha fallado (intento %d/%d) currentURl: %s: %v", retries+1, maxRetries, title, err)
		time.Sleep(2 * time.Second) // Wait before retrying
	}

	html, err := page.MustSearch(containerSelector).HTML()
	if err != nil {
		return err
	}

	if err := fs.SaveHTML(info.URL, []byte(html)); err != nil {
		return err
	}

	log.Infof("Saved HTML for %s", targetURL)

	return nil
}
