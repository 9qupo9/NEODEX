package ui

import "dex/internal/ui/styles"

func RenderCSS() string {
	return "<style>\n" +
		styles.Vars() +
		styles.Layout() +
		styles.Header() +
		styles.Orderbook() +
		styles.Chart() +
		styles.Orderform() +
		styles.History() +
		"\n</style>"
}
