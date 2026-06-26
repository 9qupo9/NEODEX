package ui

import (
	"dex/internal/ui/scripts"
	"dex/internal/ui/walletconnect"
)

func RenderJS() string {
	wc := walletconnect.NewWalletConnectService()
	return "<script>\n" +
		scripts.State() +
		scripts.UI() +
		scripts.API() +
		scripts.WS() +
		scripts.Render() +
		wc.GetScript() +
		scripts.Init() +
		"\n</script>"
}
