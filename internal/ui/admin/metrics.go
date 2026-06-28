package admin

func GetMetricsCSS() string {
	return `
        /* Метрики */
        .metrics-container {
            grid-column: span 12;
            display: grid;
            grid-template-columns: repeat(6, 1fr);
            gap: 16px;
        }

        .metric-box {
            background-color: transparent;
            border-right: 1px solid var(--border-color);
            padding-right: 16px;
        }
        
        .metric-box:last-child {
            border-right: none;
        }

        .metric-label {
            font-size: 12px;
            color: var(--text-muted);
            margin-bottom: 8px;
        }

        .metric-val {
            font-family: 'JetBrains Mono', monospace;
            font-size: 24px;
            color: var(--text-main);
        }
    `
}

func RenderMetrics() string {
	return `
            <div class="panel metrics-container">
                <div class="metric-box">
                    <div class="metric-label">VOLUME (24H)</div>
                    <div class="metric-val" id="metric-volume">0.00</div>
                </div>
                <div class="metric-box">
                    <div class="metric-label">REVENUE</div>
                    <div class="metric-val" id="metric-revenue">0.00</div>
                </div>
                <div class="metric-box">
                    <div class="metric-label">ORDERBOOK</div>
                    <div class="metric-val" id="metric-orders">0</div>
                </div>
                <div class="metric-box">
                    <div class="metric-label">REGISTRATIONS</div>
                    <div class="metric-val" id="metric-users">0</div>
                </div>
                <div class="metric-box">
                    <div class="metric-label">WS CLIENTS</div>
                    <div class="metric-val" id="metric-ws">0</div>
                </div>
                <div class="metric-box">
                    <div class="metric-label">TCP BOTS</div>
                    <div class="metric-val" id="metric-bots">0</div>
                </div>
            </div>`
}
