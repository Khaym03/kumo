package pagepool

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
)

func SetupStealth(page *rod.Page) error {
	_, err := page.EvalOnNewDocument(stealth.JS)
	return err
}
