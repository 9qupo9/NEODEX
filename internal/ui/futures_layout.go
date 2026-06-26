package ui

func RenderFuturesPage() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>NEODEX - Futures</title>
    ` + RenderCSS() + `
</head>
<body>
    <div class="app-container">
        ` + RenderHeader("futures") + `
        
        <main class="main-layout">
            <div class="panel left-panel">
                ` + RenderOrderbook() + `
                ` + RenderFuturesHistory() + `
            </div>
            
            <div class="panel chart-panel">
                ` + RenderChart() + `
            </div>
            
            <div class="panel right-panel">
                ` + RenderFuturesOrderForm() + `
            </div>
        </main>
    </div>
    
    <!-- Market Selector Modal -->
    <div id="marketModal" class="market-modal-overlay" style="display:none;">
        <div class="market-modal">
            <div class="market-modal-header">
                <span class="market-modal-title">Select Market</span>
                <button class="market-modal-close" id="marketModalClose">✕</button>
            </div>
            <div class="market-modal-search">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor" style="color:var(--text-muted)"><path d="M15.5 14h-.79l-.28-.27A6.471 6.471 0 0 0 16 9.5 6.5 6.5 0 1 0 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/></svg>
                <input type="text" id="modalSearchInput" placeholder="Search coin..." autocomplete="off">
            </div>
            <div class="market-modal-tabs">
                <button class="market-tab active" data-quote="USDT">USDT</button>
                <button class="market-tab" data-quote="USDC">USDC</button>
            </div>
            <div class="market-modal-cols">
                <span>Pair</span>
                <span>Price</span>
                <span>24h Change</span>
            </div>
            <div class="market-modal-list" id="marketModalList">
                <div class="market-modal-loading">Loading markets...</div>
            </div>
        </div>
    </div>
    
    ` + RenderJS() + `
</body>
</html>`
}
