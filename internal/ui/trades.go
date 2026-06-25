package ui

func RenderTrades() string {
	return `
	<div style="flex: 1; display: flex; flex-direction: column; max-height: 50%; border-top: 1px solid var(--border-color);">
		<div style="padding: 10px 15px; font-size: 13px; font-weight: 600; border-bottom: 1px solid var(--border-color); color: var(--text-main);">
			Market Trades
		</div>
		<div class="table-container" style="flex: 1; overflow-y: auto;">
			<table id="marketTradesTable" style="display: none;">
				<thead>
					<tr>
						<th>Price(USDT)</th>
						<th style="text-align: right">Amount(BTC)</th>
						<th style="text-align: right">Time</th>
					</tr>
				</thead>
				<tbody id="marketTradesContainer">
					<!-- Real trades will be injected here via JS -->
				</tbody>
			</table>
			
			<div id="marketTradesEmptyState" class="empty-state" style="display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 30px; color: var(--text-muted); font-size: 12px;">
				<svg width="30" height="30" viewBox="0 0 24 24" fill="var(--border-color)"><path d="M16 6l2.29 2.29-4.88 4.88-4-4L2 16.59 3.41 18l6-6 4 4 6.3-6.29L22 12V6z"/></svg>
				<div style="margin-top: 10px;">No Recent Trades</div>
			</div>
		</div>
	</div>
	`
}
