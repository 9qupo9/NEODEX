package scripts

func API() string {
	return `
async function fetchBalance() {
    try {
        if (!isWalletConnected || !walletAddress) {
            availUSDT = 0;
            availBase = 0;
        } else {
            const res = await fetch("/api/v1/balance?accountId=" + walletAddress);
            const data = await res.json();
            availUSDT = parseFloat(data[currentQuote]) || 0;
            availBase = parseFloat(data[currentBase]) || 0;
        }
        
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
        if (!isWalletConnected || !walletAddress) {
            window.allUserOrders = [];
        } else {
            const res = await fetch("/api/v1/orders?accountId=" + walletAddress);
            const data = await res.json();
            window.allUserOrders = data || [];
        }
        renderOrders(window.allUserOrders);
    } catch(e) {
        console.error("Failed to fetch orders", e);
    }
}

async function fetchPositions() {
    try {
        if (!isWalletConnected || !walletAddress) {
            window.allUserPositions = [];
        } else {
            const res = await fetch("/api/v1/futures/positions?accountId=" + walletAddress);
            const data = await res.json();
            window.allUserPositions = data || [];
        }
        if (typeof renderPositions === 'function') {
            renderPositions(window.allUserPositions);
        }
    } catch(e) {
        console.error("Failed to fetch positions", e);
    }
}

async function submitOrder() {
    if (!isWalletConnected || !walletAddress) {
        alert("Please connect your wallet first.");
        return;
    }

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

    const orderData = {
        accountId: walletAddress,
        base: currentBase,
        quote: currentQuote,
        side: currentSide,
        type: currentType,
        price: (currentType === 'MARKET' ? 0 : p).toString(),
        qty: q.toString(),
        timestamp: Date.now()
    };
    
    // Futures specific fields
    const isFutures = window.location.pathname.startsWith('/futures');
    if (isFutures) {
        orderData.leverage = parseInt(document.getElementById('leverageSelect').value) || 10;
        orderData.marginMode = document.getElementById('marginModeSelect').value || 'CROSS';
    }
    
    if (isWalletConnected && typeof window.ethereum !== 'undefined') {
        const domain = {
            name: "NEODEX",
            version: "1",
            chainId: 1
        };
        const types = {
            Order: [
                { name: "accountId", type: "string" },
                { name: "base", type: "string" },
                { name: "quote", type: "string" },
                { name: "side", type: "string" },
                { name: "type", type: "string" },
                { name: "price", type: "string" },
                { name: "qty", type: "string" },
                { name: "timestamp", type: "uint256" }
            ]
        };
        
        try {
            const signature = await window.ethereum.request({
                method: 'eth_signTypedData_v4',
                params: [walletAddress, JSON.stringify({
                    types: {
                        EIP712Domain: [
                            { name: "name", type: "string" },
                            { name: "version", type: "string" },
                            { name: "chainId", type: "uint256" }
                        ],
                        Order: types.Order
                    },
                    primaryType: "Order",
                    domain: domain,
                    message: orderData
                })]
            });
            orderData.signature = signature;
        } catch (e) {
            console.error("Ордер отклонен пользователем:", e);
            alert("Отменено: вы не подписали транзакцию.");
            return;
        }
    }
    
    try {
        const endpoint = isFutures ? "/api/v1/futures/order" : "/api/v1/order";
        const res = await fetch(endpoint, {
            method: "POST",
            headers: {"Content-Type":"application/json"},
            body: JSON.stringify(orderData)
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
