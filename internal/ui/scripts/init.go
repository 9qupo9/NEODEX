package scripts

func Init() string {
	return `
// Initial calls
const urlParams = new URLSearchParams(window.location.search);
const initSymbol = urlParams.get('symbol');

fetchBinanceSymbols().then(() => {
    if (initSymbol && allBinanceSymbols.length > 0) {
        const found = allBinanceSymbols.find(s => s.symbol.toUpperCase() === initSymbol.toUpperCase());
        if (found) {
            switchSymbol(found);
        }
    }
});

fetchBalance();
fetchOrders();
connectWS();

setInterval(fetchBalance, 5000);
setInterval(fetchOrders, 5000);
`
}
