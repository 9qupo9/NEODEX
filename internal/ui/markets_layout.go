package ui

import (
	"dex/internal/ui/scripts"
	"dex/internal/ui/styles"
)

func RenderMarketsPage() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>NEODEX - Markets</title>
    <style>
        ` + styles.Vars() + `
        ` + styles.Header() + `
        ` + styles.RenderMarketsStyles() + `
    </style>
</head>
<body>
    <div class="app-container">
        ` + RenderHeader("markets") + `
        
		<!-- Global Stats Banner -->
		<div class="global-stats-banner">
			<div class="global-stats-content">
				<div class="gs-item">Global Market Cap: <span id="gsMarketCap" class="gs-val">--</span></div>
				<div class="gs-item">24h Vol: <span id="gsVol" class="gs-val">--</span></div>
				<div class="gs-item">BTC Dominance: <span id="gsBtcDom" class="gs-val">--</span></div>
				<div class="gs-item">Market Trend: <span id="gsTrend" class="gs-val">--</span></div>
			</div>
		</div>

        <main class="markets-main">
            <!-- Highlight Cards Section -->
            <section class="highlight-cards">
                <div class="hl-card">
                    <div class="hl-card-header">
                        <span class="hl-title">
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="var(--color-sell)" style="margin-right:8px; vertical-align:text-bottom;"><path d="M11.99 2C8.54 4.54 5.5 8.16 5.5 12.38c0 3.08 1.48 5.86 3.78 7.57-1.1-1.46-1.57-3.32-1.2-5.1.04-.21.28-.27.4-.1.6 1.05 1.54 1.84 2.65 2.22.42.14.86-.14.93-.58.17-1.12.11-2.28-.24-3.4-.11-.35-.25-.7-.41-1.03-.43-.88-.98-1.7-1.63-2.42-.14-.15-.05-.4.15-.4 2.8.03 5.48 1.63 6.84 4.1.1.18.36.14.41-.06.27-1.07.24-2.21-.13-3.26-.52-1.48-1.48-2.78-2.73-3.72-1.3-.98-2.88-1.58-4.57-1.74-.24-.02-.33-.31-.18-.46C10.74 3.03 11.38 2.5 11.99 2z"/></svg>
                            Hot Coins
                        </span>
                    </div>
                    <div class="hl-list" id="hlHotCoins">
                        <div class="hl-loading">Loading...</div>
                    </div>
                </div>
                <div class="hl-card">
                    <div class="hl-card-header">
                        <span class="hl-title">
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="var(--color-buy)" style="margin-right:8px; vertical-align:text-bottom;"><path d="M16 6l2.29 2.29-4.88 4.88-4-4L2 16.59 3.41 18l6-6 4 4 6.3-6.29L22 12V6z"/></svg>
                            Top Gainers
                        </span>
                    </div>
                    <div class="hl-list" id="hlTopGainers">
                        <div class="hl-loading">Loading...</div>
                    </div>
                </div>
                <div class="hl-card">
                    <div class="hl-card-header">
                        <span class="hl-title">
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="var(--primary-color)" style="margin-right:8px; vertical-align:text-bottom;"><path d="M3.5 18.49l6-6.01 4 4L22 6.92l-1.41-1.41-7.09 7.97-4-4L2 16.99z"/></svg>
                            Top Volume
                        </span>
                    </div>
                    <div class="hl-list" id="hlTopVolume">
                        <div class="hl-loading">Loading...</div>
                    </div>
                </div>
				<div class="hl-card">
					<div class="hl-card-header">
						<span class="hl-title" style="color: #0ECB81;">
							<svg width="16" height="16" viewBox="0 0 24 24" fill="#0ECB81" style="margin-right:8px; vertical-align:text-bottom;"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm5 11h-4v4h-2v-4H7v-2h4V7h2v4h4v2z"/></svg>
							New Listings
						</span>
					</div>
					<div class="hl-list" id="hlNewListings">
						<div class="hl-loading">Loading...</div>
					</div>
				</div>
            </section>

            <!-- Markets Table Section -->
            <section class="markets-table-section">
                <div class="markets-tabs">
                    <button class="market-tab active" data-filter="all">All Cryptos</button>
                    <button class="market-tab" data-filter="gainers">Top Gainers</button>
                    <button class="market-tab" data-filter="losers">Top Losers</button>
					<button class="market-tab" data-filter="meme">Meme</button>
					<button class="market-tab" data-filter="ai">AI</button>
					<button class="market-tab" data-filter="defi">DeFi</button>
					<button class="market-tab" data-filter="l1">Layer 1 / 2</button>
                    <div class="markets-search">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="var(--text-muted)"><path d="M15.5 14h-.79l-.28-.27A6.471 6.471 0 0 0 16 9.5 6.5 6.5 0 1 0 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/></svg>
                        <input type="text" id="marketsSearch" placeholder="Search coin Name">
                    </div>
                </div>

                <div class="table-container">
                    <table class="markets-table">
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th class="right-align">Price</th>
                                <th class="right-align">24h Change</th>
                                <th class="right-align hide-mobile">24h High</th>
								<th class="right-align hide-mobile">Market Cap</th>
                                <th class="right-align hide-mobile">24h Volume</th>
                                <th class="center-align hide-mobile">Trend</th>
                                <th class="right-align hide-mobile">Action</th>
                            </tr>
                        </thead>
                        <tbody id="marketsTableBody">
                            <tr><td colspan="8" style="text-align: center; padding: 40px;">Loading markets data...</td></tr>
                        </tbody>
                    </table>
                </div>
            </section>
        </main>
    </div>

    <script>
        ` + scripts.RenderMarketsScripts() + `
    </script>
</body>
</html>`
}
