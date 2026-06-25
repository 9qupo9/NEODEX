package scripts

func WS() string {
	return `
// Websocket logic
let wsIntentionalClose = false;

function connectWS() {
    wsIntentionalClose = false;
    if (ws) {
        wsIntentionalClose = true;
        ws.close();
    }
    
    const s = currentSymbol;
    ws = new WebSocket("wss://stream.binance.com:9443/stream?streams=" + s + "@depth20@100ms/" + s + "@trade/" + s + "@ticker");
    
    ws.onopen = () => {
        wsIntentionalClose = false;
    };

    ws.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        if (!msg.stream || !msg.data) return;
        
        // Guard: ignore messages from old symbol if symbol changed mid-flight
        if (!msg.stream.startsWith(currentSymbol)) return;
        
        const stream = msg.stream;
        const data = msg.data;
        
        if (stream.endsWith("@depth20@100ms")) {
            renderOrderbook(data.asks, data.bids);
        } else if (stream.endsWith("@ticker")) {
            currentMarketPrice = parseFloat(data.c);
            const priceChange = parseFloat(data.P);
            const tickerEl = document.getElementById('tickerPrice');
            const decimals = getPriceDecimals(currentMarketPrice);
            if (tickerEl) {
                tickerEl.innerText = currentMarketPrice.toFixed(decimals);
                tickerEl.style.color = priceChange >= 0 ? "var(--color-buy)" : "var(--color-sell)";
            }
            const obMid = document.querySelector('.ob-mid');
            if (obMid) {
                obMid.innerHTML = currentMarketPrice.toFixed(decimals) + ' <span class="fiat">$' + currentMarketPrice.toFixed(decimals) + '</span>';
                obMid.style.color = priceChange >= 0 ? "var(--color-buy)" : "var(--color-sell)";
            }
        } else if (stream.endsWith("@trade")) {
            renderMarketTrades([data]);
        }
    };
    
    ws.onclose = () => {
        if (!wsIntentionalClose) {
            console.log("Binance WS closed unexpectedly, reconnecting in 2s...");
            setTimeout(() => {
                if (!wsIntentionalClose) connectWS();
            }, 2000);
        }
    };

    ws.onerror = (err) => {
        console.error("WS error:", err);
    };
}

// Helper: determine how many decimal places to show based on price magnitude
function getPriceDecimals(price) {
    if (price >= 1000) return 2;
    if (price >= 100)  return 2;
    if (price >= 10)   return 3;
    if (price >= 1)    return 4;
    if (price >= 0.1)  return 5;
    if (price >= 0.01) return 6;
    if (price >= 0.001) return 7;
    return 8;
}
`
}
