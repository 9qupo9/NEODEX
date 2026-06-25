package ui

import "dex/internal/ui/scripts"

func RenderJS() string {
	return "<script>\n" +
		scripts.State() +
		scripts.UI() +
		scripts.API() +
		scripts.WS() +
		scripts.Render() +
		scripts.Init() +
		"\n</script>"
}
