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
