package admin

func GetChartCSS() string {
	return `
        /* График */
        .chart-panel {
            grid-column: span 8;
            height: 400px;
        }

        .chart-wrapper {
            flex: 1;
            width: 100%;
        }
    `
}

func RenderChart() string {
	return `
            <div class="panel chart-panel">
                <div class="panel-title">MATCHING LATENCY</div>
                <div id="metric-latency" style="font-family: 'JetBrains Mono'; font-size: 14px; margin-bottom: 12px; color: var(--success-color);">0.0000 ms</div>
                <div id="chart" class="chart-wrapper"></div>
            </div>`
}
