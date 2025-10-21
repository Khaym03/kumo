package pagepool

import (
	"strings"

	fingerprint "github.com/Khaym03/kumo/internal/adapters/fingerprint"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	log "github.com/sirupsen/logrus"
)

func SetupHeaders(c fingerprint.CombineFingerprint) PageSetupFunc {
	log.Info(c.Headers)

	return func (page *rod.Page) error {
		
		if len(c.Headers) == 0 {
			return nil
		}

		var arr []string

		for k, v := range c.Headers {
			switch {
			case strings.EqualFold(k, "User-Agent"):
				userAgentParams := &proto.NetworkSetUserAgentOverride{
					UserAgent: v,
				}
				if err := page.SetUserAgent(userAgentParams); err != nil {
					log.Errorf("headless: could not set user agent: %v", err)
				}
			default:
				arr = append(arr, k, v)
			}
		}

		if len(arr) > 0 {
			_, err := page.SetExtraHeaders(arr)
			if err != nil {
				log.Errorf("headless: could not set extra headers: %v", err)
			}
		}

		return nil
	}
}

