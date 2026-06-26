package ui

func RenderOrderForm() string {
	return `
	<div class="order-form-panel">
		<div class="order-tabs" id="orderTypeTabs">
			<div class="tab active" data-type="MARKET">Market</div>
		</div>

		<!-- Quote Currency Tabs (like Binance) -->
		<div class="quote-tabs" id="quoteTabs">
			<button class="quote-tab active" data-quote="USDT">USDT</button>
			<button class="quote-tab" data-quote="USDC">USDC</button>
			<button class="quote-tab" data-quote="BTC">BTC</button>
			<button class="quote-tab" data-quote="ETH">ETH</button>
		</div>
		
		<div class="side-toggle-container" style="padding: 12px 12px 0 12px;">
			<div class="side-toggle" id="sideToggle">
				<div class="side-btn active buy" data-side="BUY">Buy</div>
				<div class="side-btn sell" data-side="SELL">Sell</div>
			</div>
		</div>
		
		<div class="order-forms-wrapper">
			<div class="balance-row">
				<span>Avail:</span>
				<span id="availBalance">0.00 USDT</span>
			</div>
			
			<div class="input-group" id="stopInputGroup" style="display: none;">
				<span class="input-label">Stop</span>
				<input type="number" id="stopInput" placeholder="0.00">
				<span class="input-suffix" id="stopSuffix">USDT</span>
			</div>
			
			<div class="input-group" id="priceInputGroup">
				<span class="input-label">Price</span>
				<input type="number" id="priceInput" placeholder="Market" disabled>
				<span class="input-suffix" id="priceSuffix">USDT</span>
			</div>
			
			<div class="input-group">
				<span class="input-label">Amount</span>
				<input type="number" id="qtyInput" placeholder="0.00000">
				<span class="input-suffix" id="qtySuffix">BTC</span>
			</div>
			
			<div class="slider-container buy-mode">
				<div class="slider-track-wrap">
					<div class="slider-bg"></div>
					<div class="slider-fill" id="sliderFill"></div>
					<div class="slider-marks">
						<div class="slider-mark" onclick="setSlider(0)"></div>
						<div class="slider-mark" onclick="setSlider(25)"></div>
						<div class="slider-mark" onclick="setSlider(50)"></div>
						<div class="slider-mark" onclick="setSlider(75)"></div>
						<div class="slider-mark" onclick="setSlider(100)"></div>
					</div>
					<input type="range" min="0" max="100" value="0" class="slider-range" id="qtySlider">
				</div>
				<div class="slider-labels">
					<span onclick="setSlider(0)">0%</span>
					<span onclick="setSlider(25)">25%</span>
					<span onclick="setSlider(50)">50%</span>
					<span onclick="setSlider(75)">75%</span>
					<span onclick="setSlider(100)">100%</span>
				</div>
			</div>
			
			<div class="balance-row" style="margin-top: 4px;">
				<span>Total:</span>
				<span id="totalUSDT">0.00 USDT</span>
			</div>
			
			<div class="action-row" style="margin-top: 10px;">
				<button class="btn btn-main btn-buy" id="submitOrderBtn">Buy BTC</button>
			</div>
		</div>
		
		<!-- Account Balance Block -->
		<div class="account-balance-block">
			<div class="acct-bal-header">
				<svg width="13" height="13" viewBox="0 0 24 24" fill="currentColor"><path d="M21 18v1c0 1.1-.9 2-2 2H5c-1.11 0-2-.9-2-2V5c0-1.1.89-2 2-2h14c1.1 0 2 .9 2 2v1h-9c-1.11 0-2 .9-2 2v8c0 1.1.89 2 2 2h9zm-9-2h10V8H12v8zm4-2.5c-.83 0-1.5-.67-1.5-1.5s.67-1.5 1.5-1.5 1.5.67 1.5 1.5-.67 1.5-1.5 1.5z"/></svg>
				<span>My Balance</span>
			</div>
			<div class="acct-bal-row">
				<span class="acct-bal-coin" id="acctBaseCoin">BTC</span>
				<span class="acct-bal-amount" id="acctBaseAmount">0.00000000</span>
			</div>
			<div class="acct-bal-row">
				<span class="acct-bal-coin" id="acctQuoteCoin">USDT</span>
				<span class="acct-bal-amount" id="acctQuoteAmount">0.00</span>
			</div>
			<div class="acct-bal-row acct-bal-total-row">
				<span class="acct-bal-coin">Est. Value</span>
				<span class="acct-bal-amount" id="acctTotalUSD">≈ $0.00</span>
			</div>
		</div>
	</div>
	`
}
