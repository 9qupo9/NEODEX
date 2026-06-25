package scripts

func RenderMarketsScripts() string {
	return `
	let allMarkets = [];
	let currentFilter = 'all'; // 'all', 'gainers', 'losers', 'meme', 'ai', 'defi', 'l1'

	// Hardcoded categories for demonstration
	const CATEGORIES = {
		meme: ['DOGEUSDT','SHIBUSDT','PEPEUSDT','WIFUSDT','FLOKIUSDT','BONKUSDT','BOMEUSDT','MEMEUSDT'],
		ai: ['FETUSDT','AGIXUSDT','RNDRUSDT','OCEANUSDT','TAOUSDT','WLDUSDT','NFPUSDT','AIUSDT'],
		defi: ['UNIUSDT','AAVEUSDT','MKRUSDT','CRVUSDT','CAKEUSDT','LDOUSDT','COMPUSDT','SNXUSDT'],
		l1: ['BTCUSDT','ETHUSDT','SOLUSDT','ADAUSDT','AVAXUSDT','DOTUSDT','SUIUSDT','APTUSDT','INJUSDT']
	};

	const tbody = document.getElementById('marketsTableBody');
	const searchInput = document.getElementById('marketsSearch');
	const tabs = document.querySelectorAll('.market-tab');

	let dynamicNewListings = [];

	async function fetchMarkets() {
		try {
			// Fetch new listings first
			const newRes = await fetch('/api/v1/new-listings');
			dynamicNewListings = await newRes.json();

			const res = await fetch('https://api.binance.com/api/v3/ticker/24hr');
			const data = await res.json();
			
			allMarkets = data.filter(s => s.symbol.endsWith('USDT') || s.symbol.endsWith('USDC') || s.symbol.endsWith('BTC'));
			
			updateGlobalStats();
			updateHighlightCards();
			renderTable();
		} catch (err) {
			console.error(err);
			tbody.innerHTML = '<tr><td colspan="8" style="text-align:center; color: var(--color-sell);">Error loading market data</td></tr>';
		}
	}

	function updateGlobalStats() {
		let totalVol = 0;
		let totalCap = 0;
		let upCount = 0;
		let downCount = 0;
		let btcCap = 0;

		allMarkets.forEach(m => {
			if (m.symbol.endsWith('USDT')) {
				const vol = parseFloat(m.quoteVolume);
				totalVol += vol;
				
				// Mock Market Cap based on 24h Volume (e.g., Vol * 15 for realistic looking numbers)
				let cap = vol * 15;
				if (m.symbol === 'BTCUSDT') { cap = 1.2e12; btcCap = cap; } // Hardcode BTC cap for realism
				else if (m.symbol === 'ETHUSDT') { cap = 350e9; }
				
				totalCap += cap;
			}
			if (parseFloat(m.priceChangePercent) > 0) upCount++;
			else if (parseFloat(m.priceChangePercent) < 0) downCount++;
		});

		const volStr = totalVol > 1e9 ? (totalVol/1e9).toFixed(2) + 'B' : (totalVol/1e6).toFixed(2) + 'M';
		const capStr = totalCap > 1e12 ? (totalCap/1e12).toFixed(2) + 'T' : (totalCap/1e9).toFixed(2) + 'B';
		const btcDom = btcCap > 0 && totalCap > 0 ? ((btcCap / totalCap) * 100).toFixed(1) + '%' : '--';

		document.getElementById('gsMarketCap').innerText = '$' + capStr;
		document.getElementById('gsVol').innerText = '$' + volStr;
		document.getElementById('gsBtcDom').innerText = btcDom;

		const trendEl = document.getElementById('gsTrend');
		if (upCount > downCount) {
			trendEl.innerText = 'Bullish';
			trendEl.style.color = 'var(--color-buy)';
		} else {
			trendEl.innerText = 'Bearish';
			trendEl.style.color = 'var(--color-sell)';
		}
	}

	function generateSparkline(changePercent) {
		const isUp = changePercent >= 0;
		const color = isUp ? 'var(--color-buy)' : 'var(--color-sell)';
		const points = isUp ? '0,25 20,20 40,25 60,10 80,5' : '0,5 20,10 40,5 60,20 80,25';
		return ` + "`" + `
			<svg class="sparkline-svg" viewBox="0 0 80 30" preserveAspectRatio="none">
				<polyline fill="none" stroke="${color}" stroke-width="2" points="${points}" stroke-linecap="round" stroke-linejoin="round"/>
			</svg>
		` + "`" + `;
	}

	function getCoinIcon(base) {
		const url = "/api/v1/icon?s=" + base;
		return ` + "`" + `<img src="${url}" onerror="this.onerror=null; this.outerHTML='<div class=\\'coin-icon\\'>${base.substring(0,2)}</div>'" style="width:32px; height:32px; border-radius:50%;">` + "`" + `;
	}

	function updateHighlightCards() {
		let gainers = [...allMarkets].sort((a,b) => parseFloat(b.priceChangePercent) - parseFloat(a.priceChangePercent)).slice(0, 3);
		let hot = [...allMarkets].filter(m => parseFloat(m.priceChangePercent) > 0).sort((a,b) => parseFloat(b.quoteVolume) - parseFloat(a.quoteVolume)).slice(0, 3);
		let volume = [...allMarkets].sort((a,b) => parseFloat(b.quoteVolume) - parseFloat(a.quoteVolume)).slice(0, 3);
		
		let newListings = [];
		if (dynamicNewListings && dynamicNewListings.length > 0) {
			newListings = allMarkets.filter(m => dynamicNewListings.includes(m.symbol)).slice(0, 3);
		} else {
			const recentCryptoSymbols = ['NOTUSDT', 'TONUSDT', 'ZKUSDT', 'ZROUSDT', 'IOUSDT', 'WUSDT', 'ENAUSDT'];
			newListings = allMarkets.filter(m => recentCryptoSymbols.includes(m.symbol)).slice(0, 3);
		}

		function renderCardList(elementId, items) {
			const el = document.getElementById(elementId);
			if (!el) return;
			let html = '';
			items.forEach((m, idx) => {
				const price = parseFloat(m.lastPrice);
				const change = parseFloat(m.priceChangePercent);
				const colorClass = change >= 0 ? 'text-buy' : 'text-sell';
				const sign = change > 0 ? '+' : '';
				const base = m.symbol.replace(/(USDT|USDC|BTC|ETH)$/, '');
				html += ` + "`" + `
					<div class="hl-item" onclick="goToSpot('${m.symbol}')">
						<div class="hl-item-left">
							<span class="hl-rank">${idx + 1}</span>
							${getCoinIcon(base)}
							<span class="hl-coin">${base}</span>
						</div>
						<div class="hl-item-right">
							<span class="hl-price">${formatPrice(price)}</span>
							<span class="hl-change ${colorClass}">${sign}${change.toFixed(2)}%</span>
						</div>
					</div>
				` + "`" + `;
			});
			el.innerHTML = html;
		}

		renderCardList('hlTopGainers', gainers);
		renderCardList('hlHotCoins', hot);
		renderCardList('hlTopVolume', volume);
		renderCardList('hlNewListings', newListings);
	}

	function renderTable() {
		const searchVal = searchInput.value.toLowerCase();
		
		let filtered = allMarkets.filter(m => m.symbol.toLowerCase().includes(searchVal));

		if (currentFilter === 'gainers') {
			filtered.sort((a,b) => parseFloat(b.priceChangePercent) - parseFloat(a.priceChangePercent));
			filtered = filtered.slice(0, 50);
		} else if (currentFilter === 'losers') {
			filtered.sort((a,b) => parseFloat(a.priceChangePercent) - parseFloat(b.priceChangePercent));
			filtered = filtered.slice(0, 50);
		} else if (CATEGORIES[currentFilter]) {
			const allowed = CATEGORIES[currentFilter];
			filtered = filtered.filter(m => allowed.includes(m.symbol));
			filtered.sort((a,b) => parseFloat(b.quoteVolume) - parseFloat(a.quoteVolume));
		} else {
			filtered.sort((a,b) => parseFloat(b.quoteVolume) - parseFloat(a.quoteVolume));
			filtered = filtered.slice(0, 200);
		}

		if (filtered.length === 0) {
			tbody.innerHTML = '<tr><td colspan="8" style="text-align:center;">No coins found</td></tr>';
			return;
		}

		let html = '';
		filtered.forEach(m => {
			const price = parseFloat(m.lastPrice);
			const change = parseFloat(m.priceChangePercent);
			const colorClass = change >= 0 ? 'text-buy' : 'text-sell';
			const sign = change > 0 ? '+' : '';
			
			const baseAsset = m.symbol.replace(/(USDT|USDC|BTC|ETH)$/, '');
			const quoteAsset = m.symbol.replace(baseAsset, '');
			
			const vol = parseFloat(m.quoteVolume);
			const volStr = vol > 1e6 ? (vol/1e6).toFixed(2) + 'M' : vol.toFixed(2);

			// Mock Market Cap for visual completeness
			let cap = vol * 15;
			if (m.symbol === 'BTCUSDT') cap = 1.2e12;
			else if (m.symbol === 'ETHUSDT') cap = 350e9;
			const capStr = cap > 1e9 ? '$' + (cap/1e9).toFixed(2) + 'B' : '$' + (cap/1e6).toFixed(2) + 'M';

			html += ` + "`" + `
				<tr onclick="goToSpot('${m.symbol}')">
					<td>
						<div class="coin-name-cell">
							${getCoinIcon(baseAsset)}
							<div>
								<span class="coin-symbol">${baseAsset}</span><span class="coin-quote">/${quoteAsset}</span>
							</div>
						</div>
					</td>
					<td class="right-align font-mono">${formatPrice(price)}</td>
					<td class="right-align font-mono ${colorClass}">${sign}${change.toFixed(2)}%</td>
					<td class="right-align font-mono hide-mobile">${formatPrice(parseFloat(m.highPrice))}</td>
					<td class="right-align font-mono hide-mobile">${capStr}</td>
					<td class="right-align font-mono hide-mobile">${volStr}</td>
					<td class="center-align hide-mobile">
						${generateSparkline(change)}
					</td>
					<td class="right-align hide-mobile">
						<button class="action-btn">Trade</button>
					</td>
				</tr>
			` + "`" + `;
		});
		
		tbody.innerHTML = html;
	}

	function formatPrice(p) {
		if (p < 0.001) return p.toFixed(6);
		if (p < 1) return p.toFixed(4);
		return p.toFixed(2);
	}

	window.goToSpot = function(symbol) {
		window.location.href = '/spot?symbol=' + symbol;
	};

	tabs.forEach(tab => {
		tab.addEventListener('click', (e) => {
			tabs.forEach(t => t.classList.remove('active'));
			e.target.classList.add('active');
			currentFilter = e.target.getAttribute('data-filter');
			renderTable();
		});
	});

	searchInput.addEventListener('input', () => {
		renderTable();
	});

	fetchMarkets();
	setInterval(fetchMarkets, 10000);
	`
}
