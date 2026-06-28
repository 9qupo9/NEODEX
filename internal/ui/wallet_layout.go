package ui

import (
	"dex/internal/ui/scripts"
	"dex/internal/ui/styles"
)

func RenderWalletPage() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>NEODEX - Wallet</title>
    <style>
        ` + styles.Vars() + `
        ` + styles.Header() + `
        ` + styles.RenderWalletStyles() + `
    </style>
</head>
<body>
    <div class="app-container">
        ` + RenderHeader("wallet") + `
        
        <main class="wallet-main">
			<div class="wallet-dashboard">
				
				<!-- Left Column: Total Balance & Charts -->
				<div class="wallet-left">
					<div class="wallet-card balance-card">
						<div class="balance-header">
							<h2>Estimated Balance</h2>
							<svg class="eye-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
						</div>
						<div class="balance-amount" id="totalUsdBalance">≈ $0.00</div>
						<div class="balance-pnl">
							<span class="pnl-label">Today's PnL:</span>
							<span class="pnl-value up" id="todayPnl">+$0.00 (+0.00%)</span>
						</div>
						
						<div class="balance-breakdown">
							<div class="breakdown-item">
								<div class="bd-label">Available Balance</div>
								<div class="bd-value" id="availUsdTotal">$0.00</div>
							</div>
							<div class="breakdown-item">
								<div class="bd-label">Locked in Orders</div>
								<div class="bd-value" id="lockedUsdTotal">$0.00</div>
							</div>
						</div>
					</div>

					<div class="wallet-card chart-card">
						<h3>Asset Allocation</h3>
						<div class="pie-chart-container">
							<div class="pie-chart" id="assetPieChart"></div>
							<div class="pie-legend" id="assetPieLegend">
								<!-- Dynamically populated -->
							</div>
						</div>
					</div>
				</div>

				<!-- Right Column: Asset List -->
				<div class="wallet-right">
					<div class="wallet-card assets-card">
						<div class="assets-header">
							<h3>Your Assets</h3>
							<div class="assets-search">
								<input type="text" id="assetSearchInput" placeholder="Search coin">
							</div>
						</div>
						
						<div class="table-container">
							<table class="assets-table">
								<thead>
									<tr>
										<th>Asset</th>
										<th class="right-align">Total</th>
										<th class="right-align hide-mobile">Available</th>
										<th class="right-align hide-mobile">In Order</th>
										<th class="right-align">Value (USD)</th>
										<th class="right-align hide-mobile">Action</th>
									</tr>
								</thead>
								<tbody id="assetsTableBody">
									<tr><td colspan="6" style="text-align: center; padding: 40px;">Loading assets...</td></tr>
								</tbody>
							</table>
						</div>
					</div>
					
					<!-- Staking / Earn Section -->
					<div class="wallet-card assets-card" style="margin-top: 20px;">
						<div class="assets-header">
							<h3>DeFi Earn (Staking)</h3>
						</div>
						
						<div class="table-container">
							<table class="assets-table">
								<thead>
									<tr>
										<th>Asset</th>
										<th class="right-align">Amount Staked</th>
										<th class="right-align hide-mobile">Est. APY</th>
										<th class="right-align hide-mobile">Current Reward</th>
										<th class="right-align">Action</th>
									</tr>
								</thead>
								<tbody id="stakingTableBody">
									<tr><td colspan="5" style="text-align: center; padding: 40px;">No active stakes</td></tr>
								</tbody>
							</table>
						</div>
						
						<div class="stake-action-box" style="padding: 15px; border-top: 1px solid var(--border-color); display: flex; gap: 10px; align-items: center; margin-top: 10px;">
							<select id="stakeAssetSelect" class="btn" style="background: var(--bg-hover); color: var(--text-main); border: 1px solid var(--border-color);">
								<option value="USDT">USDT (10% APY)</option>
								<option value="BTC">BTC (5% APY)</option>
								<option value="ETH">ETH (5% APY)</option>
							</select>
							<input type="number" id="stakeAmountInput" placeholder="Amount" style="background: var(--bg-hover); border: 1px solid var(--border-color); color: var(--text-main); padding: 8px; border-radius: 4px; width: 120px;">
							<select id="stakeDaysSelect" class="btn" style="background: var(--bg-hover); color: var(--text-main); border: 1px solid var(--border-color);">
								<option value="30">30 Days</option>
								<option value="90">90 Days</option>
								<option value="365">1 Year</option>
							</select>
							<button class="btn btn-main btn-buy" onclick="stakeFunds()">Stake Now</button>
						</div>
					</div>
				</div>

			</div>
        </main>
    </div>

    <script>
        ` + scripts.RenderWalletScripts() + `
    </script>
</body>
</html>`
}
