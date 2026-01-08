package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Khaym03/kumo/internal/pkg/types"
	"github.com/Khaym03/kumo/internal/ports"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

const (
	itemsSelector   = `#ausform > div .item .inner .col_3.clfix`
	anchorSelector  = `div:nth-child(1) a[href^="/exhibitor/"]`
	countrySelector = `div:nth-child(1) p`
	hallSelector    = `div:nth-child(3) p a`
)

const domain = "https://www.anuga.com"
const phase1Identifier = "[anuga-phase-1-collector]"
const phase2Identifier = "[anuga-phase-2-collector]"

type AnugaPhase1Collector struct{}

func (a *AnugaPhase1Collector) String() string {
	return phase1Identifier
}

func (a *AnugaPhase1Collector) ProcessPage(
	ctx context.Context,
	page *rod.Page,
	req *types.Request,
	queue ports.Enqueuer,
	fs ports.FileStorage,
) error {
	var err error

	// Check for context cancellation before starting
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Use Context to navigate
	if err = page.Context(ctx).Navigate(req.URL); err != nil {
		return err
	}

	// Use Context to handle cookies
	cookiesEl, err := page.Context(ctx).Timeout(time.Second * 30).Element(`//*[@id="onetrust-reject-all-handler"]`)
	if err == nil {
		log.Println("Aceptando cookies...")
		// Click and wait with context
		if err := cookiesEl.Click(proto.InputMouseButtonLeft, 1); err != nil {
			return err
		}
		if err = page.WaitLoad(); err != nil {
			return err
		}
	}

	// Pass context to the helper function
	details, err := a.getDetails(ctx, page)
	if err != nil {
		// Return the error so it can be handled
		return fmt.Errorf(
			"error al obtener detalles de la página %s: %v",
			req.URL,
			err,
		)
	}

	for _, item := range details {
		// Check for context cancellation inside the loop
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		jsonData, err := json.Marshal(&item)
		if err != nil {
			log.Printf("error al convertir a JSON, omitiendo: %v", err)
			continue
		}

		if err = fs.SaveJSON(item.URL, jsonData); err != nil {
			return fmt.Errorf("fallo al guardar el json: %w", err)
		}

		// Enqueue with the context
		queue.Dispatch(&types.Request{
			URL:       item.URL,
			Collector: phase2Identifier,
		})
	}

	return nil
}

func (a *AnugaPhase1Collector) getDetails(ctx context.Context, page *rod.Page) ([]Phase1Info, error) {
	// Use Context to wait for elements
	if _, err := page.Context(ctx).Timeout(time.Second * 20).Element("#ausform"); err != nil {
		return nil, fmt.Errorf("no se encontró el contenedor principal 'div#ausform': %w", err)
	}

	// Use Context to get elements
	items, err := page.Context(ctx).Elements(itemsSelector)
	if err != nil {
		return nil, fmt.Errorf("no se encontraron elementos para raspar: %w", err)
	}

	var output []Phase1Info

	for _, item := range items {
		// Check for context cancellation inside the loop
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Use regular methods and check for errors
		anchor, err := item.Element(anchorSelector)
		if err != nil {
			log.Printf("Advertencia: No se pudo encontrar el enlace del expositor, omitiendo este item: %v", err)
			continue
		}

		title, err := anchor.Attribute("title")
		if err != nil {
			return nil, err
		}

		href, err := anchor.Attribute("href")
		if err != nil {
			return nil, err
		}

		countryEl, err := item.Element(countrySelector)
		if err != nil {
			return nil, err
		}
		country, err := countryEl.Text()
		if err != nil {
			return nil, err
		}

		hallAnchor, err := item.Element(hallSelector)
		var hallAndStand string
		if err == nil && hallAnchor != nil {
			hallAndStand, err = hallAnchor.Text()
			if err != nil {
				return nil, err
			}
		}

		output = append(output, Phase1Info{
			CompanyName:  strings.TrimSpace(*title),
			URL:          domain + *href,
			Country:      strings.TrimSpace(country),
			HallAndStand: strings.TrimSpace(hallAndStand),
		})
	}

	return output, nil
}

type Phase1Info struct {
	CompanyName  string `json:"company_name"`
	Country      string `json:"country"`
	HallAndStand string `json:"hall_and_stand"`
	URL          string `json:"url"`
}
