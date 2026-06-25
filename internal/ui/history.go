package ui

func RenderHistory() string {
	return `
	<div class="history-panel" style="flex: 1; border-top: 1px solid var(--border-color); display: flex; flex-direction: column;">
		<div class="history-tabs" style="display: flex; border-bottom: 1px solid var(--border-color);">
			<div class="tab history-tab-btn active" data-tab="open" style="flex: 1; text-align: center; padding: 10px 0; font-size: 11px; cursor: pointer; color: var(--color-accent); border-bottom: 2px solid var(--color-accent);">Open Orders(0)</div>
			<div class="tab history-tab-btn" data-tab="history" style="flex: 1; text-align: center; padding: 10px 0; font-size: 11px; cursor: pointer; color: var(--text-muted); border-bottom: 2px solid transparent;">History</div>
		</div>
		
		<div class="table-container" style="flex: 1; overflow-y: auto;">
			<table id="historyTable" style="display: none;">
				<thead>
					<tr>
						<th>Side</th>
						<th>Price</th>
						<th>Amount</th>
					</tr>
				</thead>
				<tbody id="historyTableBody">
					<!-- Data will be injected here via JS -->
				</tbody>
			</table>
			<div id="historyEmptyState" class="empty-state" style="display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 30px; color: var(--text-muted); font-size: 12px;">
				<svg width="30" height="30" viewBox="0 0 24 24" fill="var(--border-color)"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"/></svg>
				<div style="margin-top: 10px;">No Open Orders</div>
			</div>
		</div>
	</div>
	`
}
