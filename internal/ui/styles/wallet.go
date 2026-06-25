package styles

func RenderWalletStyles() string {
	return `
	.wallet-main {
		padding: 24px;
		max-width: 1400px;
		margin: 0 auto;
		width: 100%;
		box-sizing: border-box;
	}

	.wallet-dashboard {
		display: grid;
		grid-template-columns: 350px 1fr;
		gap: 24px;
	}

	@media (max-width: 900px) {
		.wallet-dashboard {
			grid-template-columns: 1fr;
		}
	}

	.wallet-card {
		background: #181a20;
		border: 1px solid rgba(255, 255, 255, 0.05);
		border-radius: 12px;
		padding: 32px;
		margin-bottom: 24px;
		box-shadow: 0 4px 24px rgba(0,0,0,0.2);
		transition: transform 0.2s ease;
	}

	.wallet-card h2, .wallet-card h3 {
		margin: 0 0 16px 0;
		color: #EAECEF;
		font-weight: 600;
	}

	.wallet-card h2 {
		font-size: 20px;
	}
	
	.wallet-card h3 {
		font-size: 16px;
	}

	.balance-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 16px;
	}

	.balance-header h2 {
		margin-bottom: 0 !important;
	}

	.eye-icon {
		color: var(--text-muted);
		cursor: pointer;
		transition: color 0.2s;
	}
	.eye-icon:hover {
		color: #EAECEF;
	}

	.balance-amount {
		font-size: 36px;
		font-weight: 700;
		color: #fff;
		margin-bottom: 8px;
		font-family: monospace;
		letter-spacing: -1px;
	}

	.balance-pnl {
		font-size: 14px;
		margin-bottom: 24px;
		display: flex;
		gap: 8px;
	}

	.pnl-label {
		color: var(--text-muted);
	}

	.pnl-value {
		font-weight: 600;
	}
	.pnl-value.up { color: #0ECB81; }
	.pnl-value.down { color: #F6465D; }

	.balance-breakdown {
		display: flex;
		background: rgba(255, 255, 255, 0.02);
		border-radius: 8px;
		padding: 16px;
		gap: 24px;
	}

	.breakdown-item {
		flex: 1;
	}

	.bd-label {
		color: var(--text-muted);
		font-size: 13px;
		margin-bottom: 6px;
	}

	.bd-value {
		color: #EAECEF;
		font-size: 15px;
		font-weight: 600;
		font-family: monospace;
	}

	.wallet-actions {
		display: flex;
		gap: 16px;
	}

	.w-btn {
		flex: 1;
		padding: 12px;
		font-weight: 600;
		font-size: 14px;
		border-radius: 8px;
		cursor: pointer;
		text-align: center;
		transition: all 0.2s ease;
		border: none;
	}

	.w-btn-primary {
		background: var(--color-accent);
		color: #111;
	}
	.w-btn-primary:hover {
		background: #F4C443;
		transform: translateY(-1px);
	}

	.w-btn-secondary {
		background: rgba(255, 255, 255, 0.1);
		color: #EAECEF;
	}
	.w-btn-secondary:hover {
		background: rgba(255, 255, 255, 0.15);
		transform: translateY(-1px);
	}

	.w-btn-outline {
		background: transparent;
		color: #EAECEF;
		border: 1px solid rgba(255, 255, 255, 0.2);
	}
	.w-btn-outline:hover {
		border-color: rgba(255, 255, 255, 0.4);
		background: rgba(255, 255, 255, 0.05);
		transform: translateY(-1px);
	}

	/* Pie Chart */
	.pie-chart-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 24px;
	}

	.pie-chart {
		width: 220px;
		height: 220px;
		border-radius: 50%;
		background: conic-gradient(#3468D0 0% 30%, #F3BA2F 30% 70%, #0ECB81 70% 100%);
		position: relative;
		box-shadow: inset 0 0 10px rgba(0,0,0,0.5);
	}
	
	.pie-chart::after {
		content: "";
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		width: 170px;
		height: 170px;
		background: #181a20;
		border-radius: 50%;
		box-shadow: 0 4px 12px rgba(0,0,0,0.5);
	}

	.pie-legend {
		display: flex;
		flex-direction: column;
		width: 100%;
		gap: 16px;
		margin-top: 12px;
	}

	.legend-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		font-size: 14px;
	}

	.legend-left {
		display: flex;
		align-items: center;
		gap: 8px;
		color: #EAECEF;
	}

	.legend-color {
		width: 12px;
		height: 12px;
		border-radius: 50%;
	}

	.legend-value {
		font-weight: 600;
		color: #fff;
	}
	.legend-pct {
		color: var(--text-muted);
		width: 45px;
		text-align: right;
	}

	/* Assets Table */
	.assets-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 16px;
	}

	.assets-search input {
		background: var(--bg-hover);
		border: 1px solid var(--border-color);
		color: #fff;
		padding: 8px 12px;
		border-radius: 6px;
		width: 200px;
		font-size: 14px;
	}
	.assets-search input:focus {
		border-color: var(--color-accent);
		outline: none;
	}

	.assets-table {
		width: 100%;
		border-collapse: collapse;
		font-size: 14px;
	}

	.assets-table th {
		text-align: left;
		padding: 16px 12px;
		color: var(--text-muted);
		font-weight: 500;
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
	}

	.assets-table td {
		padding: 20px 12px;
		border-bottom: 1px solid rgba(255, 255, 255, 0.02);
		color: #EAECEF;
		vertical-align: middle;
		transition: background 0.2s ease;
	}
	
	.assets-table tr:hover td {
		background: rgba(255,255,255,0.03);
	}

	.right-align {
		text-align: right !important;
	}

	.asset-name-cell {
		display: flex;
		align-items: center;
		gap: 12px;
	}
	
	.asset-icon {
		width: 32px;
		height: 32px;
		border-radius: 50%;
	}
	
	.asset-symbol {
		font-weight: 600;
		color: #fff;
		font-size: 15px;
	}
	
	.asset-fullname {
		color: var(--text-muted);
		font-size: 13px;
		margin-top: 2px;
	}

	.action-links {
		display: flex;
		gap: 16px;
		justify-content: flex-end;
	}
	
	.action-link {
		color: #EAECEF;
		text-decoration: none;
		font-weight: 600;
		font-size: 13px;
		padding: 6px 12px;
		background: rgba(255, 255, 255, 0.05);
		border-radius: 4px;
		transition: background 0.2s;
	}
	.action-link:hover {
		background: rgba(255, 255, 255, 0.1);
		color: var(--color-accent);
	}

	@media (max-width: 768px) {
		.hide-mobile { display: none !important; }
	}
	`
}
