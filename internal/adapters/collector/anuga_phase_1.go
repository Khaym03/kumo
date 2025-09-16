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
	page.MustNavigate(req.URL)
	err := page.Timeout(30 * time.Second).WaitLoad()
	if err != nil {
		return err
	}
	if _, err := page.Timeout(time.Second * 30).Element(`//*[@id="onetrust-reject-all-handler"]`); err == nil {
		log.Println("Aceptando cookies...")
		page.MustElement(`//*[@id="onetrust-reject-all-handler"]`).MustClick()
		// Después de la navegación y el clic, espera a que el contenido se cargue
		page.MustWaitLoad()
	}

	details, err := a.getDetails(page)
	if err != nil {
		return fmt.Errorf(
			"error al obtener detalles de la página %s: %v",
			req.URL,
			err,
		)
	}

	for _, item := range details {
		jsonData, err := json.Marshal(&item)
		if err != nil {
			log.Printf("error al convertir a JSON, omitiendo: %v", err)
			return err
		}

		if err = fs.SaveJSON(item.URL, jsonData); err != nil {
			return fmt.Errorf("fallo al guardar el json: %w", err)
		}

		queue.Enqueue(&types.Request{
			URL: item.URL,
		})
	}

	return nil
}

func (a *AnugaPhase1Collector) getDetails(page *rod.Page) ([]Phase1Info, error) {
	if _, err := page.Timeout(time.Second * 20).Element("#ausform"); err != nil {
		return nil, fmt.Errorf("no se encontró el contenedor principal 'div#ausform': %w", err)
	}

	items, err := page.Elements(itemsSelector)
	if err != nil {
		return nil, fmt.Errorf("no se encontraron elementos para raspar: %w", err)
	}

	var output []Phase1Info

	for _, item := range items {
		var exhibitor Phase1Info

		anchor, err := item.Element(anchorSelector)
		if err != nil {
			log.Printf("Advertencia: No se pudo encontrar el enlace del expositor, omitiendo este item: %v", err)
			continue
		}

		exhibitor.CompanyName = strings.TrimSpace(*anchor.MustAttribute("title"))
		exhibitor.URL = domain + *anchor.MustAttribute("href")
		exhibitor.Country = strings.TrimSpace(item.MustElement(countrySelector).MustText())

		hallAnchor, err := item.Element(hallSelector)
		if err == nil && hallAnchor != nil {
			exhibitor.HallAndStand = strings.TrimSpace(hallAnchor.MustText())
		} else {
			exhibitor.HallAndStand = ""
		}

		output = append(output, exhibitor)
	}

	return output, nil
}

type Phase1Info struct {
	CompanyName  string `json:"company_name"`
	Country      string `json:"country"`
	HallAndStand string `json:"hall_and_stand"`
	URL          string `json:"url"`
}
