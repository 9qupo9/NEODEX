package scripts

func Render() string {
	return `
function renderOrderbook(asks, bids) {
	if (!asksContainer || !bidsContainer) return;
	
    asksContainer.innerHTML = '';
    bidsContainer.innerHTML = '';
    
    // ---------------------------------------------------------
    // GROUPING LOGIC BY currentTickSize
    // ---------------------------------------------------------
    function groupData(data, tickSize, isAsk) {
        let grouped = new Map();
        (data || []).forEach(item => {
            const p = parseFloat(Array.isArray(item) ? item[0] : item.price);
            const q = parseFloat(Array.isArray(item) ? item[1] : (item.volume || item.qty || 0));
            if (isNaN(p) || isNaN(q)) return;
            
            const factor = 1 / tickSize;
            let bucket = isAsk ? Math.ceil(p * factor) / factor : Math.floor(p * factor) / factor;
            
            let decimals = 0;
            if (tickSize < 1) decimals = tickSize < 0.1 ? 2 : 1;
            const key = bucket.toFixed(decimals);
            
            grouped.set(key, (grouped.get(key) || 0) + q);
        });
        
        let arr = [];
        grouped.forEach((q, p) => arr.push({ price: p, qty: q }));
        arr.sort((a, b) => parseFloat(a.price) - parseFloat(b.price));
        if (!isAsk) arr.reverse();
        return arr;
    }

    const groupedAsks = groupData(asks, currentTickSize, true);
    const groupedBids = groupData(bids, currentTickSize, false);

    let maxTotal = 0;
    
    // We want to calculate cumulative from lowest -> highest
    let askCum = 0;
    const askRows = groupedAsks.map(a => {
        const price = a.price;
        const v = parseFloat(a.qty);
        askCum += v;
        maxTotal = Math.max(maxTotal, askCum);
        return { p: price, q: v, t: askCum };
    }).reverse(); // Then reverse to display highest price on top

    let bidCum = 0;
    const bidRows = groupedBids.map(b => {
        const price = b.price;
        const v = parseFloat(b.qty);
        bidCum += v;
        maxTotal = Math.max(maxTotal, bidCum);
        return { p: price, q: v, t: bidCum };
    });

    let displayLimit = 15; // Show 15 per side

    askRows.slice(-displayLimit).forEach(a => {
        const price = parseFloat(a.p);
        const dec = getPriceDecimals(price);
        const p = price.toFixed(dec);
        const q = parseFloat(a.q).toFixed(6);
        const t = a.t.toFixed(4);
        const depth = maxTotal > 0 ? (a.t / maxTotal) * 100 : 0;
        
        const row = document.createElement('div');
        row.className = 'ob-row ask';
        row.onclick = function() { document.getElementById('priceInput').value = p; };
        row.innerHTML = '<div class="depth-bg" style="width:' + depth + '%"></div>' +
            '<span class="price ask">' + p + '</span>' +
            '<span>' + q + '</span>' +
            '<span>' + t + '</span>';
        asksContainer.appendChild(row);
    });

    bidRows.slice(0, displayLimit).forEach(b => {
        const price = parseFloat(b.p);
        const dec = getPriceDecimals(price);
        const p = price.toFixed(dec);
        const q = parseFloat(b.q).toFixed(6);
        const t = b.t.toFixed(4);
        const depth = maxTotal > 0 ? (b.t / maxTotal) * 100 : 0;
        
        const row = document.createElement('div');
        row.className = 'ob-row bid';
        row.onclick = function() { document.getElementById('priceInput').value = p; };
        row.innerHTML = '<div class="depth-bg" style="width:' + depth + '%"></div>' +
            '<span class="price bid">' + p + '</span>' +
            '<span>' + q + '</span>' +
            '<span>' + t + '</span>';
        bidsContainer.appendChild(row);
    });
}

function renderMarketTrades(trades) {
	if (!marketTradesContainer) return;
	
    if (trades.length === 0 && marketTradesContainer.children.length === 0) {
        if (marketTradesTable) marketTradesTable.style.display = 'none';
        if (marketTradesEmptyState) marketTradesEmptyState.style.display = 'flex';
        return;
    }

    if (marketTradesTable) marketTradesTable.style.display = 'table';
    if (marketTradesEmptyState) marketTradesEmptyState.style.display = 'none';

	trades.forEach(t => {
        // Handle both local format and Binance format (m=isMaker, p=price, q=qty, T=timestamp)
        const price = t.p || t.price;
        const qty = t.q || t.qty;
        
        let side = 'BUY';
        if (t.m !== undefined) {
            side = t.m ? 'SELL' : 'BUY'; // if buyer is maker, it means it's a SELL trade
        } else if (t.side) {
            side = t.side;
        }

        const timestamp = t.T || t.timestamp || Date.now();
		const p = parseFloat(price).toFixed(2);
        const q = parseFloat(qty).toFixed(4);
        const color = side === 'BUY' ? 'var(--color-buy)' : 'var(--color-sell)';
        const timeStr = new Date(timestamp).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit', second:'2-digit'});
        
        const tr = document.createElement('tr');
        tr.innerHTML = ` + "`" + `
			<td style="color:${color}">${p}</td>
			<td style="text-align:right; font-family:var(--font-mono)">${q}</td>
			<td style="text-align:right; color:var(--text-muted)">${timeStr}</td>
		` + "`" + `;
        
        marketTradesContainer.prepend(tr);
	});
    
    // Keep only top 30 rows
    while (marketTradesContainer.children.length > 30) {
        marketTradesContainer.removeChild(marketTradesContainer.lastChild);
    }
}

function renderOrders(orders) {
    if (!historyTableBody) return;

    // Filter based on active tab
    const filtered = orders.filter(o => {
        if (currentHistoryTab === 'open') {
            return o.status === 'NEW' || o.status === 'PARTIALLY_FILLED';
        } else {
            return o.status === 'FILLED' || o.status === 'CANCELED' || o.status === 'REJECTED';
        }
    });

    if (filtered.length === 0) {
        if (historyTable) historyTable.style.display = 'none';
        if (historyEmptyState) historyEmptyState.style.display = 'flex';
        // Update tab count
        if (currentHistoryTab === 'open' && historyTabBtns[0]) historyTabBtns[0].innerText = 'Open Orders(0)';
        return;
    }

    if (historyTable) historyTable.style.display = 'table';
    if (historyEmptyState) historyEmptyState.style.display = 'none';

    if (currentHistoryTab === 'open' && historyTabBtns[0]) {
        historyTabBtns[0].innerText = 'Open Orders(' + filtered.length + ')';
    }

    historyTableBody.innerHTML = '';
    filtered.slice().reverse().forEach(o => {
        const p = o.type === 'MARKET' ? 'Market' : parseFloat(o.price).toFixed(2);
        const q = parseFloat(o.qty).toFixed(4);
        const color = o.side === 'BUY' ? 'var(--color-buy)' : 'var(--color-sell)';
        
        historyTableBody.innerHTML += ` + "`" + `
			<tr>
				<td style="color:${color}; font-weight: 500;">${o.side}</td>
				<td style="font-family:var(--font-mono)">${p}</td>
				<td style="font-family:var(--font-mono)">${q}</td>
			</tr>
		` + "`" + `;
    });
}

function renderPositions(positions) {
    const isFutures = window.location.pathname.startsWith('/futures');
    if (!isFutures) return;
    
    const posTable = document.getElementById('positionsTable');
    const posTableBody = document.getElementById('positionsTableBody');
    const emptyState = document.getElementById('historyEmptyState');
    const emptyText = document.getElementById('historyEmptyText');
    const posTabBtn = document.querySelector('.history-tab-btn[data-tab="positions"]');

    if (currentHistoryTab !== 'positions') {
        if (posTable) posTable.style.display = 'none';
        // Only update badge count if needed, don't change visibility
        if (posTabBtn && positions) posTabBtn.innerText = 'Positions(' + positions.length + ')';
        return;
    }

    if (!positions || positions.length === 0) {
        if (posTable) posTable.style.display = 'none';
        if (emptyState) emptyState.style.display = 'flex';
        if (emptyText) emptyText.innerText = 'No Open Positions';
        if (posTabBtn) posTabBtn.innerText = 'Positions(0)';
        return;
    }

    if (posTable) posTable.style.display = 'table';
    if (emptyState) emptyState.style.display = 'none';
    if (posTabBtn) posTabBtn.innerText = 'Positions(' + positions.length + ')';

    if (!posTableBody) return;
    posTableBody.innerHTML = '';
    
    positions.forEach(p => {
        const sideColor = p.side === 'BUY' ? 'var(--color-buy)' : 'var(--color-sell)';
        const sideText = p.side === 'BUY' ? 'Long ' + p.leverage + 'x' : 'Short ' + p.leverage + 'x';
        
        // Mock PnL logic for demo
        const entry = parseFloat(p.entryPrice);
        const size = parseFloat(p.size);
        const markPrice = currentMarketPrice > 0 ? currentMarketPrice : entry; // Fallback
        
        let pnl = 0;
        if (entry > 0 && size > 0) {
            if (p.side === 'BUY') {
                pnl = (markPrice - entry) * size;
            } else {
                pnl = (entry - markPrice) * size;
            }
        }
        
        const pnlColor = pnl >= 0 ? 'var(--color-buy)' : 'var(--color-sell)';
        const pnlStr = (pnl > 0 ? '+' : '') + pnl.toFixed(2);
        const liqPrice = parseFloat(p.liquidationPrice).toFixed(2);

        posTableBody.innerHTML += ` + "`" + `
			<tr>
				<td><strong>${p.pair.BaseAsset}/${p.pair.QuoteAsset}</strong><br><span style="font-size:10px;color:var(--text-muted)">${p.marginMode}</span></td>
				<td style="color:${sideColor}; font-weight: 500;">${sideText}</td>
				<td style="font-family:var(--font-mono)">${size.toFixed(4)}</td>
				<td style="font-family:var(--font-mono)">${entry.toFixed(2)}</td>
				<td style="font-family:var(--font-mono); color:var(--color-accent)">${liqPrice > 0 ? liqPrice : '--'}</td>
                <td style="font-family:var(--font-mono); color:${pnlColor}; text-align:right;">${pnlStr}</td>
			</tr>
		` + "`" + `;
    });
}
`
}
