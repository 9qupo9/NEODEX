package ui

func RenderHeader(activeMenu string) string {
	marketsClass := ""
	spotClass := ""
	walletClass := ""
	switch activeMenu {
	case "markets":
		marketsClass = ` class="active"`
	case "spot":
		spotClass = ` class="active"`
	case "wallet":
		walletClass = ` class="active"`
	}

	return `
	<header class="header">
		<div class="header-left">
			<a href="/" class="logo" style="text-decoration:none; color:inherit;">NEO<span style="color:var(--text-main)">DEX</span></a>
			<nav class="nav-links">
				<a href="/markets"` + marketsClass + `>Markets</a>
				<a href="/spot"` + spotClass + `>Spot</a>
				<a href="#">Futures</a>
				<a href="/wallet"` + walletClass + `>Wallet</a>
				
				<button class="market-selector-btn" id="marketSelectorBtn">
					<svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor" style="color:var(--text-muted)"><path d="M15.5 14h-.79l-.28-.27A6.471 6.471 0 0 0 16 9.5 6.5 6.5 0 1 0 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/></svg>
					<span id="currentPairLabel">BTC/USDT</span>
					<svg width="12" height="12" viewBox="0 0 24 24" fill="currentColor"><path d="M7 10l5 5 5-5z"/></svg>
				</button>
			</nav>
		</div>
		<div class="header-right">
			<div class="lang-selector">
				<div class="lang-current">
					<svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/></svg>
					<span>English</span>
					<svg width="10" height="10" viewBox="0 0 24 24" fill="currentColor"><path d="M7 10l5 5 5-5z"/></svg>
				</div>
				<div class="lang-dropdown">
					<div class="lang-option active">English</div>
					<div class="lang-option">Русский</div>
					<div class="lang-option">中文</div>
				</div>
			</div>
		</div>
	</header>
	`
}
