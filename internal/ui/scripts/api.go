package scripts

func API() string {
	return `
async function fetchBalance() {
    try {
        const res = await fetch("/api/v1/balance?accountId=" + USER_ID);
        const data = await res.json();
        
        availUSDT = parseFloat(data[currentQuote]) || 0;
        availBase = parseFloat(data[currentBase]) || 0;
        
        // Update account balance block
        const acctBaseCoin   = document.getElementById('acctBaseCoin');
        const acctQuoteCoin  = document.getElementById('acctQuoteCoin');
        const acctBaseAmount = document.getElementById('acctBaseAmount');
        const acctQuoteAmount= document.getElementById('acctQuoteAmount');
        const acctTotalUSD   = document.getElementById('acctTotalUSD');
        if (acctBaseCoin)   acctBaseCoin.innerText   = currentBase;
        if (acctQuoteCoin)  acctQuoteCoin.innerText  = currentQuote;
        if (acctBaseAmount) acctBaseAmount.innerText  = availBase.toFixed(8);
        if (acctQuoteAmount)acctQuoteAmount.innerText = availUSDT.toFixed(2);
        if (acctTotalUSD)   acctTotalUSD.innerText    = '≈ $' + (availUSDT + availBase * currentMarketPrice).toFixed(2);
        
        updateUI();
    } catch(e) {
        console.error("Failed to fetch balance", e);
    }
}

async function fetchBinanceSymbols() {
    try {
        const newRes = await fetch('/api/v1/new-listings');
        dynamicNewListings = await newRes.json();

        const res = await fetch("https://api.binance.com/api/v3/exchangeInfo");
        const data = await res.json();
        if (data && data.symbols) {
            allBinanceSymbols = data.symbols.filter(s => 
                s.status === 'TRADING' && 
                (s.quoteAsset === 'USDT' || s.quoteAsset === 'USDC' || s.quoteAsset === 'FDUSD' || s.quoteAsset === 'BTC' || s.quoteAsset === 'ETH' || s.quoteAsset === 'BNB')
            );
        }
    } catch(e) {
        console.error("Failed to fetch binance symbols", e);
    }
}

async function fetchOrders() {
    try {
        const res = await fetch("/api/v1/orders?accountId=" + USER_ID);
        const data = await res.json();
        window.allUserOrders = data || [];
        renderOrders(window.allUserOrders);
    } catch(e) {
        console.error("Failed to fetch orders", e);
    }
}


async function submitOrder() {
    let p = parseFloat(priceInput.value);
    let q = parseFloat(qtyInput.value);
    let stop = parseFloat(stopInput.value);
    
    if (currentType !== 'MARKET' && (isNaN(p) || p <= 0)) {
        alert("Enter a valid price");
        return;
    }
    if (isNaN(q) || q <= 0) {
        alert("Enter a valid amount");
        return;
    }
    if (currentType === 'STOP_LIMIT' && (isNaN(stop) || stop <= 0)) {
        alert("Enter a valid stop price");
        return;
    }

    const payload = {
        accountId: USER_ID,
        base: currentBase,
        quote: currentQuote,
        side: currentSide,
        type: currentType,
        price: (currentType === 'MARKET' ? 0 : p).toString(),
        qty: q.toString()
    };
    
    try {
        const res = await fetch("/api/v1/order", {
            method: "POST",
            headers: {"Content-Type":"application/json"},
            body: JSON.stringify(payload)
        });
        const result = await res.json();
        if (result.success) {
            alert("Order Placed Successfully!");
            qtyInput.value = '';
            fetchBalance();
        } else {
            alert("Order failed: " + JSON.stringify(result));
        }
    } catch(e) {
        console.error(e);
        alert("Error placing order");
    }
}
`
}
