package ui

func RenderChart() string {
	return `
	<div class="chart-container" style="position: relative;">
		<!-- We embed a real TradingView Advanced Chart widget -->
		<div id="tvchart" style="height: 100%; width: 100%;"></div>
		
		<!-- HACK: Hiding TradingView watermark logo -->
		<div style="position: absolute; bottom: 0; left: 0; width: 60px; height: 35px; background: #0d0e12; z-index: 100; pointer-events: none;"></div>
		<!-- Block clicks on the symbol title without covering the top toolbar (timeframes) -->
		<div style="position: absolute; top: 42px; left: 10px; width: 300px; height: 32px; background: transparent; z-index: 100; pointer-events: auto; cursor: default;"></div>
		
		<script type="text/javascript" src="https://s3.tradingview.com/tv.js"></script>
		<script type="text/javascript">
        window.reloadChart = function(symbol) {
            document.getElementById('tvchart').innerHTML = '';
            new TradingView.widget({
              "autosize": true,
              "symbol": "BINANCE:" + symbol.toUpperCase(),
              "interval": "D",
              "timezone": "Etc/UTC",
              "theme": "dark",
              "style": "1",
              "locale": "en",
              "enable_publishing": false,
              "backgroundColor": "#0d0e12",
              "gridColor": "#262831",
              "hide_top_toolbar": false,
              "allow_symbol_change": false, // Enforce our own search
              "hide_legend": false,
              "save_image": false,
              "container_id": "tvchart"
            });
        };
        // Initial load
        window.reloadChart("BTCUSDT");
		</script>
	</div>
	`
}
