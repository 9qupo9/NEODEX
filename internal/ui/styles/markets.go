package styles

func RenderMarketsStyles() string {
	return `
	.global-stats-banner {
		background: #0d0e12;
		border-bottom: 1px solid var(--border-color);
		padding: 12px 24px;
		font-size: 13px;
		color: var(--text-muted);
	}
	.global-stats-content {
		max-width: 1200px;
		margin: 0 auto;
		display: flex;
		gap: 32px;
	}
	.gs-item {
		display: flex;
		align-items: center;
		gap: 6px;
	}
	.gs-val {
		color: #fff;
		font-weight: 600;
	}

	.markets-main {
		display: flex;
		flex-direction: column;
		width: 100%;
		max-width: 1200px;
		margin: 0 auto;
		padding: 24px;
		box-sizing: border-box;
		gap: 32px;
	}

	.highlight-cards {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 24px;
	}

	.hl-card {
		background: #181a20;
		border-radius: 12px;
		padding: 20px;
		display: flex;
		flex-direction: column;
		gap: 16px;
		position: relative;
		overflow: hidden;
		transition: transform 0.2s ease, box-shadow 0.2s ease;
	}

	.hl-card:hover {
		transform: translateY(-2px);
		box-shadow: 0 8px 24px rgba(0,0,0,0.4);
	}

	.highlight-cards > .hl-card:nth-child(1) {
		border-top: 2px solid #f3ba2f;
		background: linear-gradient(180deg, rgba(243, 186, 47, 0.03) 0%, #181a20 100%);
	}
	.highlight-cards > .hl-card:nth-child(2) {
		border-top: 2px solid var(--color-buy);
		background: linear-gradient(180deg, rgba(14, 203, 129, 0.03) 0%, #181a20 100%);
	}
	.highlight-cards > .hl-card:nth-child(3) {
		border-top: 2px solid var(--primary-color);
		background: linear-gradient(180deg, rgba(52, 104, 208, 0.03) 0%, #181a20 100%);
	}
	.highlight-cards > .hl-card:nth-child(4) {
		border-top: 2px solid #0ECB81;
		background: linear-gradient(180deg, rgba(14, 203, 129, 0.03) 0%, #181a20 100%);
	}

	.hl-card-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.hl-title {
		font-size: 18px;
		font-weight: 600;
		color: #fff;
		display: flex;
		align-items: center;
	}

	.hl-more {
		color: var(--text-muted);
		font-size: 14px;
		text-decoration: none;
		transition: color 0.2s;
	}
	.hl-more:hover { color: var(--primary-color); }

	.hl-list {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.hl-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		cursor: pointer;
		padding: 8px 12px;
		border-radius: 8px;
		margin: 0 -12px;
		transition: background 0.2s, transform 0.2s;
	}
	.hl-item:hover { 
		background: rgba(255,255,255,0.03); 
		transform: translateX(4px);
	}

	.hl-item-left {
		display: flex;
		align-items: center;
		gap: 12px;
	}
	
	.hl-rank {
		color: #5e6673;
		font-size: 15px;
		width: 14px;
		font-weight: 600;
		text-align: center;
	}

	.hl-item:nth-child(1) .hl-rank { color: #f3ba2f; text-shadow: 0 0 8px rgba(243,186,47,0.4); }
	.hl-item:nth-child(2) .hl-rank { color: #E0E6ED; text-shadow: 0 0 8px rgba(224,230,237,0.4); }
	.hl-item:nth-child(3) .hl-rank { color: #D19159; text-shadow: 0 0 8px rgba(209,145,89,0.4); }
	.hl-item:nth-child(4) .hl-rank { color: #5e6673; }

	.hl-coin {
		color: #fff;
		font-weight: 600;
		font-size: 15px;
	}

	.hl-item-right {
		display: flex;
		flex-direction: column;
		align-items: flex-end;
		gap: 2px;
	}

	.hl-price {
		color: #fff;
		font-size: 14px;
		font-weight: 600;
	}

	.hl-change {
		font-size: 13px;
		font-weight: 600;
	}

	.sparkline-svg {
		width: 80px;
		height: 30px;
		filter: drop-shadow(0px 2px 4px rgba(0,0,0,0.3));
	}

	.center-align {
		text-align: center !important;
	}

	.markets-tabs {
		display: flex;
		align-items: center;
		gap: 12px;
		margin-bottom: 20px;
		border-bottom: 1px solid var(--border-color);
		padding-bottom: 16px;
	}

	.market-tab {
		background: transparent;
		border: 1px solid transparent;
		color: var(--text-muted);
		font-size: 15px;
		font-weight: 500;
		cursor: pointer;
		padding: 8px 16px;
		border-radius: 8px;
		transition: all 0.2s ease;
	}

	.market-tab:hover {
		color: #fff;
		background: rgba(255,255,255,0.02);
	}

	.market-tab.active {
		color: #fff;
		background: rgba(52, 104, 208, 0.1);
		border: 1px solid rgba(52, 104, 208, 0.3);
		box-shadow: 0 0 12px rgba(52, 104, 208, 0.15);
	}

	.markets-search {
		margin-left: auto;
		display: flex;
		align-items: center;
		background: #0d0e12;
		padding: 10px 16px;
		border-radius: 8px;
		border: 1px solid #2b3139;
		transition: all 0.3s ease;
	}

	.markets-search:focus-within {
		border-color: var(--primary-color);
		box-shadow: 0 0 0 2px rgba(52, 104, 208, 0.2);
	}

	.markets-search input {
		background: transparent;
		border: none;
		color: #fff;
		font-size: 14px;
		margin-left: 10px;
		outline: none;
		width: 220px;
	}

	.markets-search input::placeholder {
		color: var(--text-muted);
	}

	.table-container {
		background: #181a20;
		border-radius: 12px;
		overflow: hidden;
		box-shadow: 0 8px 24px rgba(0,0,0,0.2);
	}

	.markets-table {
		width: 100%;
		border-collapse: collapse;
	}

	.markets-table th {
		text-align: left;
		padding: 16px 20px;
		color: var(--text-muted);
		font-size: 13px;
		font-weight: 500;
		border-bottom: 1px solid var(--border-color);
		background: rgba(255,255,255,0.01);
	}

	.markets-table td {
		padding: 16px 20px;
		color: var(--text-main);
		font-size: 14px;
		border-bottom: 1px solid rgba(255,255,255,0.03);
	}

	.markets-table tbody tr {
		transition: all 0.2s ease;
		cursor: pointer;
	}

	.markets-table tbody tr:hover {
		background: rgba(255,255,255,0.03);
		transform: scale(1.002);
	}

	.right-align {
		text-align: right !important;
	}

	.coin-name-cell {
		display: flex;
		align-items: center;
		gap: 14px;
	}

	.coin-icon {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		background: #2b3139;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 12px;
		font-weight: 600;
		color: #fff;
		box-shadow: 0 2px 8px rgba(0,0,0,0.2);
	}

	.coin-symbol {
		font-weight: 600;
		font-size: 16px;
		color: #fff;
	}

	.coin-quote {
		color: var(--text-muted);
		font-size: 13px;
		margin-left: 2px;
	}

	.action-btn {
		background: var(--primary-color);
		border: none;
		color: #fff;
		padding: 8px 16px;
		border-radius: 6px;
		cursor: pointer;
		font-size: 13px;
		font-weight: 600;
		transition: all 0.2s ease;
	}

	.action-btn:hover {
		filter: brightness(1.15);
		box-shadow: 0 4px 12px rgba(52, 104, 208, 0.4);
		transform: translateY(-1px);
	}

	.star-btn {
		background: transparent;
		border: none;
		color: var(--text-muted);
		cursor: pointer;
		padding: 0 8px 0 0;
		font-size: 16px;
	}
	.star-btn.active {
		color: var(--primary-color);
	}

	@media (max-width: 768px) {
		.hide-mobile {
			display: none;
		}
		.markets-hero {
			flex-direction: column;
			gap: 16px;
			align-items: flex-start;
		}
		.markets-tabs {
			flex-wrap: wrap;
		}
		.markets-search {
			margin-left: 0;
			width: 100%;
		}
		.markets-search input {
			width: 100%;
		}
	}
	`
}
