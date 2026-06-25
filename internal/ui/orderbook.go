package ui

func RenderOrderbook() string {
	return `
		<div class="ob-header" style="border-bottom: 1px solid var(--border-color); padding-bottom: 10px; padding-top: 15px; justify-content: flex-end; display: flex; padding-right: 15px;">
			<div class="custom-select-wrapper" id="customTickSelect">
                <div class="custom-select-display" id="customTickDisplay">
                    <span id="customTickValue">0.01</span>
                    <span style="font-size: 8px;">▼</span>
                </div>
                <div class="custom-select-options" id="customTickOptions">
                    <div class="custom-select-option" data-val="0.000001">0.000001</div>
                    <div class="custom-select-option" data-val="0.00001">0.00001</div>
                    <div class="custom-select-option" data-val="0.0001">0.0001</div>
                    <div class="custom-select-option" data-val="0.001">0.001</div>
                    <div class="custom-select-option active" data-val="0.01">0.01</div>
                    <div class="custom-select-option" data-val="0.1">0.1</div>
                    <div class="custom-select-option" data-val="1">1</div>
                    <div class="custom-select-option" data-val="10">10</div>
                    <div class="custom-select-option" data-val="50">50</div>
                    <div class="custom-select-option" data-val="100">100</div>
                </div>
			</div>
		</div>
		
		<div class="ob-header">
			<span>Price(<span id="obQuoteLabel">USDT</span>)</span>
			<span>Amount(<span id="obBaseLabel">BTC</span>)</span>
			<span>Total</span>
		</div>
		
		<div class="ob-asks" id="asksContainer">
			<!-- Populated by JS -->
		</div>
		
		<div class="ob-mid">
			<span id="midPrice">64,230.00</span>
			<span class="fiat">≈ $64,230.00</span>
		</div>
		
		<div class="ob-bids" id="bidsContainer">
			<!-- Populated by JS -->
		</div>
	`
}
