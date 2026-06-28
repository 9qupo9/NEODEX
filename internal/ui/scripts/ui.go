package scripts

func UI() string {
	return `
// --- Toggles ---
sideBtns.forEach(btn => {
    btn.addEventListener('click', (e) => {
        sideBtns.forEach(b => b.classList.remove('active'));
        e.target.classList.add('active');
        currentSide = e.target.getAttribute('data-side');
        updateUI();
    });
});

if (tabBtns) {
	tabBtns.forEach(btn => {
		btn.addEventListener('click', (e) => {
			tabBtns.forEach(b => b.classList.remove('active'));
			e.target.classList.add('active');
			currentType = e.target.getAttribute('data-type');
			updateUI();
		});
	});
}

const isFuturesInit = window.location.pathname.startsWith('/futures');
let currentHistoryTab = isFuturesInit ? 'positions' : 'open'; 
if (historyTabBtns) {
    historyTabBtns.forEach(btn => {
        btn.addEventListener('click', (e) => {
            historyTabBtns.forEach(b => {
                b.classList.remove('active');
                b.style.color = 'var(--text-muted)';
                b.style.borderBottomColor = 'transparent';
            });
            e.target.classList.add('active');
            e.target.style.color = 'var(--color-accent)';
            e.target.style.borderBottomColor = 'var(--color-accent)';
            currentHistoryTab = e.target.getAttribute('data-tab');
            renderOrders(window.allUserOrders || []);
            if (typeof renderPositions === 'function') {
                renderPositions(window.allUserPositions || []);
            }
        });
    });
}

// Quote Currency Tabs — switch e.g. BTC/USDT → BTC/USDC
const quoteTabs = document.querySelectorAll('.quote-tab');
quoteTabs.forEach(tab => {
    tab.addEventListener('click', async () => {
        const newQuote = tab.getAttribute('data-quote');
        if (newQuote === currentQuote) return;

        // If symbols not loaded yet, fetch first
        if (allBinanceSymbols.length === 0) await fetchBinanceSymbols();

        // Find the pair for currentBase + newQuote
        const targetSymbol = currentBase + newQuote;
        const found = allBinanceSymbols.find(s => s.symbol === targetSymbol);
        if (!found) {
            // Try reverse: maybe base is not available in this quote
            // Find ANY pair with this quote and switch to first available
            const fallback = allBinanceSymbols.find(s => s.quoteAsset === newQuote);
            if (fallback) {
                switchSymbol(fallback);
            }
            return;
        }

        // Update active tab
        quoteTabs.forEach(t => t.classList.remove('active'));
        tab.classList.add('active');

        switchSymbol(found);
    });
});

function updateQuoteTabActive() {
    quoteTabs.forEach(t => {
        t.classList.toggle('active', t.getAttribute('data-quote') === currentQuote);
    });
}

const customTickSelect = document.getElementById('customTickSelect');
const customTickDisplay = document.getElementById('customTickDisplay');
const customTickOptions = document.getElementById('customTickOptions');
const customTickValue = document.getElementById('customTickValue');

if (customTickSelect && customTickDisplay && customTickOptions) {
    customTickDisplay.addEventListener('click', (e) => {
        e.stopPropagation();
        customTickOptions.classList.toggle('open');
    });

    // Close on click outside
    document.addEventListener('click', (e) => {
        if (!customTickSelect.contains(e.target)) {
            customTickOptions.classList.remove('open');
        }
    });

    const options = customTickOptions.querySelectorAll('.custom-select-option');
    options.forEach(opt => {
        opt.addEventListener('click', (e) => {
            // Update UI
            options.forEach(o => o.classList.remove('active'));
            opt.classList.add('active');
            customTickValue.innerText = opt.innerText;
            
            // Update logic
            currentTickSize = parseFloat(opt.getAttribute('data-val'));
            
            customTickOptions.classList.remove('open');
            
            // Re-render Orderbook to apply grouping immediately
            // Since we receive deltas, it might wait 100ms, but typically depth stream is fast.
        });
    });
}

const symbolSearchInput = document.getElementById('symbolSearchInput');
const symbolSearchResults = document.getElementById('symbolSearchResults');

if (symbolSearchInput && symbolSearchResults) {
    function renderSymbolResults(results) {
        symbolSearchResults.innerHTML = '';
        if (results.length === 0) {
            symbolSearchResults.classList.remove('open');
            return;
        }
        
        results.forEach(s => {
            const isNew = dynamicNewListings && dynamicNewListings.includes(s.symbol);
            const badge = isNew ? '<span style="font-size: 10px; margin-left: 4px; padding: 2px 4px; background: rgba(14, 203, 129, 0.2); color: #0ECB81; border-radius: 4px; font-weight: bold;">NEW</span>' : '';
            
            const div = document.createElement('div');
            div.className = 'symbol-result-item';
            div.innerHTML = '<span>' + s.baseAsset + '/' + s.quoteAsset + badge + '</span>';
            div.addEventListener('click', () => {
                symbolSearchInput.value = '';
                symbolSearchResults.classList.remove('open');
                switchSymbol(s);
            });
            symbolSearchResults.appendChild(div);
        });
        symbolSearchResults.classList.add('open');
    }

    symbolSearchInput.addEventListener('input', async (e) => {
        const query = e.target.value.toLowerCase().trim();
        if (!query) {
            symbolSearchResults.classList.remove('open');
            return;
        }
        
        // If symbols not loaded yet, fetch first
        if (allBinanceSymbols.length === 0) {
            await fetchBinanceSymbols();
        }
        
        const filtered = allBinanceSymbols.filter(s => 
            s.symbol.toLowerCase().includes(query) || 
            s.baseAsset.toLowerCase().includes(query)
        ).slice(0, 50);
        
        renderSymbolResults(filtered);
    });

    document.addEventListener('click', (e) => {
        if (!symbolSearchInput.contains(e.target) && !symbolSearchResults.contains(e.target)) {
            symbolSearchResults.classList.remove('open');
        }
    });
}

// ==================== MARKET SELECTOR MODAL ====================
const marketModal    = document.getElementById('marketModal');
const marketSelectorBtn = document.getElementById('marketSelectorBtn');
const marketModalClose  = document.getElementById('marketModalClose');
const marketModalList   = document.getElementById('marketModalList');
const modalSearchInput  = document.getElementById('modalSearchInput');
const marketTabs        = document.querySelectorAll('.market-tab');

let currentModalQuote = 'USDT';
let marketTickers = {};   // symbol -> { price, change }

function openMarketModal() {
    marketModal.style.display = 'flex';
    if (modalSearchInput) { modalSearchInput.value = ''; modalSearchInput.focus(); }
    
    // Sync modal tabs to current quote
    currentModalQuote = currentQuote;
    marketTabs.forEach(t => {
        t.classList.toggle('active', t.getAttribute('data-quote') === currentModalQuote);
    });
    
    fetchModalTickers(currentModalQuote);
}

function closeMarketModal() {
    marketModal.style.display = 'none';
}

if (marketSelectorBtn) marketSelectorBtn.addEventListener('click', openMarketModal);
if (marketModalClose)  marketModalClose.addEventListener('click', closeMarketModal);
if (marketModal) {
    marketModal.addEventListener('click', (e) => {
        if (e.target === marketModal) closeMarketModal();
    });
}

marketTabs.forEach(tab => {
    tab.addEventListener('click', () => {
        marketTabs.forEach(t => t.classList.remove('active'));
        tab.classList.add('active');
        currentModalQuote = tab.getAttribute('data-quote');
        if (modalSearchInput) modalSearchInput.value = '';
        fetchModalTickers(currentModalQuote);
    });
});

if (modalSearchInput) {
    modalSearchInput.addEventListener('input', () => {
        renderModalList(currentModalQuote, modalSearchInput.value.toLowerCase().trim());
    });
}

async function fetchModalTickers(quote) {
    if (!marketModalList) return;
    marketModalList.innerHTML = '<div class="market-modal-loading">Loading...</div>';
    try {
        const res = await fetch('https://api.binance.com/api/v3/ticker/24hr');
        const data = await res.json();
        marketTickers = {};
        data.forEach(t => { marketTickers[t.symbol] = t; });
        renderModalList(quote, modalSearchInput ? modalSearchInput.value.toLowerCase().trim() : '');
    } catch(e) {
        marketModalList.innerHTML = '<div class="market-modal-loading">Failed to load.</div>';
    }
}

function renderModalList(quote, query) {
    if (!marketModalList) return;
    const symbols = allBinanceSymbols.filter(s => {
        if (s.quoteAsset !== quote) return false;
        if (!query) return true;
        return s.baseAsset.toLowerCase().includes(query) || s.symbol.toLowerCase().includes(query);
    });

    if (symbols.length === 0) {
        marketModalList.innerHTML = '<div class="market-modal-loading">No results.</div>';
        return;
    }

    marketModalList.innerHTML = '';
    symbols.forEach(s => {
        const ticker = marketTickers[s.symbol] || {};
        const price   = ticker.lastPrice  ? parseFloat(ticker.lastPrice)  : 0;
        const change  = ticker.priceChangePercent ? parseFloat(ticker.priceChangePercent) : 0;
        const dec     = getPriceDecimals(price);
        const isActive = s.symbol.toLowerCase() === currentSymbol;
        const isNew = dynamicNewListings && dynamicNewListings.includes(s.symbol);
        const badge = isNew ? '<span style="font-size: 9px; margin-left: 4px; padding: 1px 3px; background: rgba(14, 203, 129, 0.2); color: #0ECB81; border-radius: 3px; vertical-align: top;">NEW</span>' : '';

        const div = document.createElement('div');
        div.className = 'market-item' + (isActive ? ' active-pair' : '');
        div.innerHTML =
            '<span><span class="pair-name">' + s.baseAsset + '</span><span class="pair-quote">/' + s.quoteAsset + '</span>' + badge + '</span>' +
            '<span class="pair-price">' + (price > 0 ? price.toFixed(dec) : '—') + '</span>' +
            '<span class="pair-change ' + (change >= 0 ? 'up' : 'down') + '">' +
                (change >= 0 ? '+' : '') + change.toFixed(2) + '%' +
            '</span>';

        div.addEventListener('click', () => {
            switchSymbol(s);
            closeMarketModal();
        });
        marketModalList.appendChild(div);
    });
}

// Central symbol switch function used by both modal and inline search
function switchSymbol(s) {
    const newSymbol = s.symbol.toLowerCase();
    if (newSymbol === currentSymbol) return;

    currentSymbol = newSymbol;
    currentBase   = s.baseAsset;
    currentQuote  = s.quoteAsset;

    // Sync quote tabs
    updateQuoteTabActive();

    // Update header button label
    const pairLabel = document.getElementById('currentPairLabel');
    if (pairLabel) pairLabel.innerText = s.baseAsset + '/' + s.quoteAsset;

    // Update OB header labels
    const obBaseLabel = document.getElementById('obBaseLabel');
    const obQuoteLabel = document.getElementById('obQuoteLabel');
    if (obBaseLabel) obBaseLabel.innerText = currentBase;
    if (obQuoteLabel) obQuoteLabel.innerText = currentQuote;

    // Reset total label
    const totalEl = document.getElementById('totalUSDT');
    if (totalEl) totalEl.innerText = '0.00 ' + currentQuote;

    // Reload chart
    if (window.reloadChart) window.reloadChart(currentSymbol);

    // Clear stale data
    if (asksContainer) asksContainer.innerHTML = '';
    if (bidsContainer) bidsContainer.innerHTML = '';
    if (marketTradesContainer) marketTradesContainer.innerHTML = '';

    fetchBalance();
    updateUI();
    connectWS();
}
// ================================================================

function updateSliderFill() {
    if(qtySlider && sliderFill) {
        sliderFill.style.width = qtySlider.value + '%';
    }
}

function updateUI() {
    if (!submitBtn) return;
    
    // Update dynamic input suffixes
    const priceSuffix = document.getElementById('priceSuffix');
    const stopSuffix  = document.getElementById('stopSuffix');
    const qtySuffix   = document.getElementById('qtySuffix');
    if (priceSuffix) priceSuffix.innerText = currentQuote;
    if (stopSuffix)  stopSuffix.innerText  = currentQuote;
    if (qtySuffix)   qtySuffix.innerText   = currentBase;
    
    // 1. Update button style, text, and slider color mode
    submitBtn.onclick = submitOrder;
    const isFutures = window.location.pathname.startsWith('/futures');
    if (currentSide === 'BUY') {
        submitBtn.className = 'btn btn-main btn-buy';
        submitBtn.innerText = isFutures ? ('Open Long ' + currentBase) : ('Buy ' + currentBase);
        if (availBalanceEl) availBalanceEl.innerText = availUSDT.toFixed(2) + ' ' + currentQuote;
        if (sliderContainer) {
            sliderContainer.classList.remove('sell-mode');
            sliderContainer.classList.add('buy-mode');
        }
    } else {
        submitBtn.className = 'btn btn-main btn-sell';
        submitBtn.innerText = isFutures ? ('Open Short ' + currentBase) : ('Sell ' + currentBase);
        if (availBalanceEl) availBalanceEl.innerText = availBase.toFixed(8) + ' ' + currentBase;
        if (sliderContainer) {
            sliderContainer.classList.remove('buy-mode');
            sliderContainer.classList.add('sell-mode');
        }
    }

    // 2. Update inputs visibility
    if (currentType === 'MARKET') {
        if(stopInputGroup) stopInputGroup.style.display = 'none';
        if(priceInput) {
            priceInput.disabled = true;
            priceInput.value = '';
            priceInput.placeholder = 'Market';
        }
    } else if (currentType === 'LIMIT') {
        if(stopInputGroup) stopInputGroup.style.display = 'none';
        if(priceInput) {
            priceInput.disabled = false;
            priceInput.placeholder = '0.00';
        }
    } else if (currentType === 'STOP_LIMIT' || currentType === 'TAKE_PROFIT') {
        priceInput.disabled = false;
        priceInputGroup.style.display = 'flex';
        stopInputGroup.style.display = 'flex';
    }
    
    calculateTotal();
}

function calculateTotal() {
    if (!priceInput || !qtyInput || !totalEl) return;
    
    let p = parseFloat(priceInput.value);
    let q = parseFloat(qtyInput.value);
    
    if (currentType === 'MARKET') {
        p = currentMarketPrice; // rough estimate for display
    }
    
    if (isNaN(p) || isNaN(q)) {
        totalEl.innerText = "0.00 USDT";
    } else {
        totalEl.innerText = (p * q).toFixed(2) + " USDT";
    }
}

// Range Slider Drag Logic
if (qtySlider) {
    qtySlider.addEventListener('input', (e) => {
        updateSliderFill();
        const percentage = e.target.value / 100;
        
        let p = parseFloat(priceInput.value);
        if (currentType === 'MARKET') p = currentMarketPrice;
        
        if (currentSide === 'BUY') {
            if (isNaN(p) || p <= 0) return; 
            const maxBuyQty = availUSDT / p;
            qtyInput.value = (maxBuyQty * percentage).toFixed(4);
        } else {
            qtyInput.value = (availBase * percentage).toFixed(4);
        }
        calculateTotal();
    });
}

// Sync slider when typing Amount
if (qtyInput) {
    qtyInput.addEventListener('input', () => {
        calculateTotal();
        
        let q = parseFloat(qtyInput.value);
        if (isNaN(q)) q = 0;
        
        let p = parseFloat(priceInput.value);
        if (currentType === 'MARKET') p = currentMarketPrice;
        
        let maxQty = 0;
        if (currentSide === 'BUY') {
            if (!isNaN(p) && p > 0) maxQty = availUSDT / p;
        } else {
            maxQty = availBase;
        }
        
        let pct = 0;
        if (maxQty > 0) {
            pct = (q / maxQty) * 100;
        }
        if (pct > 100) pct = 100;
        if (pct < 0) pct = 0;
        
        if (qtySlider) {
            qtySlider.value = pct;
            updateSliderFill();
        }
    });
}

if(priceInput) priceInput.addEventListener('input', () => {
    calculateTotal();
    qtyInput.dispatchEvent(new Event('input')); // re-sync slider when price changes
});

// For the visual 5 markers
function setSlider(val) {
    if(qtySlider) {
        qtySlider.value = val;
        qtySlider.dispatchEvent(new Event('input'));
    }
}
`
}
