package composer

import "github.com/go-rod/rod"

func (ac *AppComposer) ComposeRemoteBrowser() (*rod.Browser, error) {
	return ac.BrowserFactory.Get(
		ac.BrowserFactory.RemoteBrowserCreator,
	)
}
