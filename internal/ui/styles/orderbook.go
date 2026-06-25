package styles

func Orderbook() string {
	return `
/* ================= ORDERBOOK & TRADES ================= */
.ob-header {
    display: flex;
    justify-content: space-between;
    padding: 8px 12px;
    color: var(--text-muted);
    font-size: 11px;
}
.ob-row {
    display: flex;
    justify-content: space-between;
    padding: 2px 12px;
    position: relative;
    cursor: pointer;
    font-family: var(--font-mono);
    line-height: 1.5;
    font-size: 12px;
}
.ob-row:hover { background: var(--bg-panel); }
.ob-row span { z-index: 1; text-align: right; flex: 1; }
.ob-row span:first-child { text-align: left; }
.ob-row .price.ask { color: var(--color-sell); }
.ob-row .price.bid { color: var(--color-buy); }
.depth-bg {
    position: absolute;
    top: 0; right: 0; bottom: 0;
    opacity: 0.15;
    z-index: 0;
}
.ask .depth-bg { background: var(--color-sell); }
.bid .depth-bg { background: var(--color-buy); }

.ob-asks, .ob-bids { flex: 1; overflow: hidden; display: flex; flex-direction: column; }
.ob-asks { justify-content: flex-end; }
.ob-mid {
    padding: 8px 12px;
    font-size: 18px;
    font-weight: bold;
    color: var(--color-buy);
    font-family: var(--font-mono);
    border-top: 1px solid var(--border-color);
    border-bottom: 1px solid var(--border-color);
    background: var(--bg-panel);
    display: flex;
    align-items: center;
    gap: 10px;
}
.ob-mid .fiat { font-size: 13px; color: var(--text-muted); font-weight: normal; }

/* Custom Dropdown for Tick Size */
.custom-select-wrapper {
    position: relative;
    user-select: none;
    cursor: pointer;
    font-size: 11px;
}
.custom-select-display {
    display: flex;
    align-items: center;
    gap: 4px;
    color: var(--text-muted);
}
.custom-select-display:hover {
    color: var(--text-main);
}
.custom-select-options {
    position: absolute;
    top: 100%;
    right: 0;
    background: var(--bg-panel);
    border: 1px solid var(--border-color);
    border-radius: 4px;
    padding: 4px 0;
    margin-top: 4px;
    display: none;
    z-index: 100;
    box-shadow: 0 4px 12px rgba(0,0,0,0.5);
    min-width: 60px;
}
.custom-select-options.open {
    display: block;
}
.custom-select-option {
    padding: 6px 12px;
    color: var(--text-main);
    text-align: right;
}
.custom-select-option:hover, .custom-select-option.active {
    background: var(--bg-hover);
    color: var(--color-accent);
}
`
}
