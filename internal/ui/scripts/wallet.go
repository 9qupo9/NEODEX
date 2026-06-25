package scripts

func RenderWalletScripts() string {
	return `
	const USER_ID = "test_user_1";
	
	let userBalances = {};
	let userOrders = [];
	let marketPrices = {};
	let totalUSD = 0;

	const PALETTE = ['#F3BA2F', '#3468D0', '#0ECB81', '#E0E6ED', '#D19159', '#8a2be2'];

	async function initWallet() {
		try {
			// Fetch prices
			const priceRes = await fetch('https://api.binance.com/api/v3/ticker/price');
			const priceData = await priceRes.json();
			priceData.forEach(p => { marketPrices[p.symbol] = parseFloat(p.price); });
			marketPrices['USDTUSDT'] = 1;
			marketPrices['USDCUSDT'] = 1;
			marketPrices['FDUSDUSDT'] = 1;

			// Fetch balance
			const balRes = await fetch("/api/v1/balance?accountId=" + USER_ID);
			const balData = await balRes.json();
			
			// Format balances
			for (const [coin, amtStr] of Object.entries(balData)) {
				const amt = parseFloat(amtStr);
				if (amt > 0) {
					userBalances[coin] = {
						available: amt,
						inOrder: 0,
						total: amt,
						valueUsd: 0
					};
				}
			}

			// Fetch orders to calculate "In Order"
			const ordRes = await fetch("/api/v1/orders?accountId=" + USER_ID);
			const ordData = await ordRes.json();
			if (ordData) {
				ordData.forEach(o => {
					if (o.status === "NEW" || o.status === "PARTIALLY_FILLED") {
						const qty = parseFloat(o.qty) - (parseFloat(o.executedQty) || 0);
						const price = parseFloat(o.price);
						if (o.side === "BUY") {
							// For buy, quote asset is locked (qty * price)
							const quote = o.pair.quoteAsset;
							const locked = qty * price;
							if (userBalances[quote]) {
								userBalances[quote].inOrder += locked;
								userBalances[quote].total += locked;
							} else {
								userBalances[quote] = { available: 0, inOrder: locked, total: locked, valueUsd: 0 };
							}
						} else {
							// For sell, base asset is locked
							const base = o.pair.baseAsset;
							if (userBalances[base]) {
								userBalances[base].inOrder += qty;
								userBalances[base].total += qty;
							} else {
								userBalances[base] = { available: 0, inOrder: qty, total: qty, valueUsd: 0 };
							}
						}
					}
				});
			}

			calculateValues();
			renderWallet();
		} catch (e) {
			console.error("Wallet initialization error:", e);
		}
	}

	function getUsdPrice(coin) {
		if (coin === 'USDT' || coin === 'USDC' || coin === 'FDUSD') return 1;
		if (marketPrices[coin + 'USDT']) return marketPrices[coin + 'USDT'];
		if (marketPrices[coin + 'USDC']) return marketPrices[coin + 'USDC'];
		// If can't find direct, try via BTC
		if (marketPrices[coin + 'BTC'] && marketPrices['BTCUSDT']) {
			return marketPrices[coin + 'BTC'] * marketPrices['BTCUSDT'];
		}
		return 0;
	}

	function calculateValues() {
		totalUSD = 0;
		let availUsdTotal = 0;
		let lockedUsdTotal = 0;

		for (const coin in userBalances) {
			const b = userBalances[coin];
			const price = getUsdPrice(coin);
			b.valueUsd = b.total * price;
			
			totalUSD += b.valueUsd;
			availUsdTotal += b.available * price;
			lockedUsdTotal += b.inOrder * price;
		}

		// Update breakdown DOM
		const availEl = document.getElementById('availUsdTotal');
		const lockedEl = document.getElementById('lockedUsdTotal');
		if (availEl) availEl.innerText = '$' + availUsdTotal.toLocaleString('en-US', {minimumFractionDigits: 2, maximumFractionDigits: 2});
		if (lockedEl) lockedEl.innerText = '$' + lockedUsdTotal.toLocaleString('en-US', {minimumFractionDigits: 2, maximumFractionDigits: 2});
		
		// Fake PnL calculation based on current total
		const pnlPct = 2.45; // Simulated daily PnL
		const pnlVal = (totalUSD * pnlPct) / 100;
		const pnlEl = document.getElementById('todayPnl');
		if (pnlEl) {
			pnlEl.innerText = '+$' + pnlVal.toLocaleString('en-US', {minimumFractionDigits: 2, maximumFractionDigits: 2}) + ' (+' + pnlPct + '%)';
		}
	}

	function getCoinIcon(base) {
		const url = "/api/v1/icon?s=" + base;
		return ` + "`" + `<img src="${url}" onerror="this.onerror=null; this.outerHTML='<div class=\\'asset-icon\\' style=\\'background: #333; color: #fff; display: flex; align-items: center; justify-content: center; font-size: 10px;\\'>${base.substring(0,2)}</div>'" class="asset-icon">` + "`" + `;
	}

	function renderWallet() {
		document.getElementById('totalUsdBalance').innerText = '≈ $' + totalUSD.toLocaleString('en-US', {minimumFractionDigits: 2, maximumFractionDigits: 2});

		// Sort assets by USD value
		const sortedAssets = Object.keys(userBalances).map(coin => ({
			coin,
			...userBalances[coin]
		})).sort((a, b) => b.valueUsd - a.valueUsd);

		const tbody = document.getElementById('assetsTableBody');
		if (sortedAssets.length === 0) {
			tbody.innerHTML = '<tr><td colspan="6" style="text-align: center; padding: 40px;">No assets found.</td></tr>';
		} else {
			let html = '';
			sortedAssets.forEach(a => {
				html += ` + "`" + `
					<tr>
						<td>
							<div class="asset-name-cell">
								${getCoinIcon(a.coin)}
								<div>
									<div class="asset-symbol">${a.coin}</div>
									<div class="asset-fullname">${a.coin} Token</div>
								</div>
							</div>
						</td>
						<td class="right-align font-mono">${a.total.toLocaleString('en-US', {maximumFractionDigits: 8})}</td>
						<td class="right-align font-mono hide-mobile">${a.available.toLocaleString('en-US', {maximumFractionDigits: 8})}</td>
						<td class="right-align font-mono hide-mobile">${a.inOrder.toLocaleString('en-US', {maximumFractionDigits: 8})}</td>
						<td class="right-align font-mono">$${a.valueUsd.toLocaleString('en-US', {minimumFractionDigits: 2, maximumFractionDigits: 2})}</td>
						<td class="right-align hide-mobile">
							<div class="action-links">
								<a href="/spot?symbol=${a.coin === 'USDT' ? 'BTCUSDT' : a.coin + 'USDT'}" class="action-link">Trade</a>
							</div>
						</td>
					</tr>
				` + "`" + `;
			});
			tbody.innerHTML = html;
		}

		renderChart(sortedAssets);
	}

	function renderChart(assets) {
		const legendContainer = document.getElementById('assetPieLegend');
		const pieChart = document.getElementById('assetPieChart');
		
		if (totalUSD === 0 || assets.length === 0) {
			pieChart.style.background = '#333';
			legendContainer.innerHTML = '<div>No assets</div>';
			return;
		}

		let conicStops = [];
		let currentDeg = 0;
		let legendHtml = '';
		
		// Group small assets into "Others"
		let othersValue = 0;
		let chartData = [];
		
		assets.forEach((a, i) => {
			const pct = a.valueUsd / totalUSD;
			if (pct < 0.01 && assets.length > 5 && i >= 4) { // less than 1% and not top 4
				othersValue += a.valueUsd;
			} else {
				chartData.push({ name: a.coin, val: a.valueUsd, pct });
			}
		});

		if (othersValue > 0) {
			chartData.push({ name: 'Others', val: othersValue, pct: othersValue / totalUSD });
		}

		chartData.forEach((d, i) => {
			const color = PALETTE[i % PALETTE.length];
			const degrees = d.pct * 360;
			
			conicStops.push(color + ' ' + currentDeg + 'deg ' + (currentDeg + degrees) + 'deg');
			currentDeg += degrees;

			legendHtml += ` + "`" + `
				<div class="legend-item">
					<div class="legend-left">
						<div class="legend-color" style="background: ${color}"></div>
						<div class="legend-name">${d.name}</div>
					</div>
					<div class="legend-right">
						<span class="legend-value">$${d.val.toLocaleString('en-US', {minimumFractionDigits: 2, maximumFractionDigits: 2})}</span>
						<span class="legend-pct">${(d.pct * 100).toFixed(1)}%</span>
					</div>
				</div>
			` + "`" + `;
		});

		pieChart.style.background = 'conic-gradient(' + conicStops.join(', ') + ')';
		legendContainer.innerHTML = legendHtml;
	}

	document.getElementById('assetSearchInput')?.addEventListener('input', (e) => {
		const term = e.target.value.toLowerCase();
		const rows = document.querySelectorAll('#assetsTableBody tr');
		rows.forEach(row => {
			if (row.innerText.toLowerCase().includes(term)) {
				row.style.display = '';
			} else {
				row.style.display = 'none';
			}
		});
	});

	initWallet();
	`
}
