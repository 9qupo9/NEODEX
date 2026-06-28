package admin

func GetLogsCSS() string {
	return `
        /* Логи */
        .logs-panel {
            grid-column: span 4;
            height: 400px;
        }

        .logs-container {
            flex: 1;
            overflow-y: auto;
            font-family: 'JetBrains Mono', monospace;
            font-size: 11px;
            line-height: 1.5;
            background-color: #000;
            padding: 12px;
            border: 1px solid var(--border-color);
            border-radius: 2px;
        }

        .log-time { color: #555; margin-right: 8px; }
        .log-info { color: #888; }
        .log-success { color: var(--success-color); }
        .log-warn { color: var(--warning-color); }
        .log-error { color: var(--danger-color); }
    `
}

func RenderLogs() string {
	return `
            <div class="panel logs-panel">
                <div class="panel-title">EVENT LOG</div>
                <div class="logs-container" id="sysLogs">
                    <div><span class="log-time"></span> <span class="log-info">[SYSTEM]</span> Waiting for data...</div>
                </div>
            </div>`
}
