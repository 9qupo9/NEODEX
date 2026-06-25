package styles

func Header() string {
	return `
/* ================= HEADER ================= */
.header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 60px;
    background: var(--bg-panel);
    border-bottom: 1px solid var(--border-color);
    padding: 0 20px;
}
.header-left, .header-right { display: flex; align-items: center; gap: 20px; }
.logo { 
    font-size: 22px; 
    font-weight: 700; 
    color: var(--color-accent); 
    letter-spacing: 1px;
}
.nav-links { display: flex; align-items: center; gap: 16px; overflow: visible; }
.nav-links a {
    color: var(--text-muted);
    text-decoration: none;
    font-weight: 500;
    transition: var(--transition);
}
.nav-links a:hover, .nav-links a.active { color: var(--text-main); }
.header-btn {
    background: transparent;
    color: var(--color-accent);
    border: 1px solid var(--color-accent);
    padding: 6px 16px;
    border-radius: 2px;
    font-weight: 600;
    cursor: pointer;
    transition: var(--transition);
}
.header-btn:hover { 
    background: var(--color-accent); 
    color: #fff; 
}

/* Language Selector */
.lang-selector {
    position: relative;
    color: var(--text-muted);
    font-size: 13px;
    font-weight: 500;
}
.lang-current {
    display: flex;
    align-items: center;
    gap: 6px;
    cursor: pointer;
    padding: 6px 0;
    transition: var(--transition);
}
.lang-current:hover { color: var(--text-main); }
.lang-dropdown {
    position: absolute;
    top: 100%;
    right: 0;
    background: var(--bg-panel);
    border: 1px solid var(--border-color);
    border-radius: 4px;
    width: 120px;
    display: none;
    flex-direction: column;
    z-index: 100;
    box-shadow: 0 4px 12px rgba(0,0,0,0.5);
}
.lang-selector:hover .lang-dropdown { display: flex; }
.lang-option {
    padding: 10px 15px;
    cursor: pointer;
    transition: var(--transition);
}
.lang-option:hover, .lang-option.active {
    background: var(--bg-hover);
    color: var(--color-accent);
}

/* Symbol Search */
.symbol-search-container {
    position: relative;
    flex-shrink: 0;
}
#symbolSearchInput {
    background: var(--bg-hover);
    border: 1px solid var(--border-color);
    color: var(--text-main);
    padding: 7px 12px;
    border-radius: 6px;
    font-family: var(--font-ui);
    font-size: 13px;
    width: 180px;
    outline: none;
    transition: var(--transition);
}
#symbolSearchInput:focus {
    border-color: var(--color-accent);
}
.symbol-search-results {
    position: absolute;
    top: 100%;
    left: 0;
    width: 100%;
    background: var(--bg-panel);
    border: 1px solid var(--border-color);
    border-top: none;
    border-radius: 0 0 4px 4px;
    max-height: 300px;
    overflow-y: auto;
    display: none;
    z-index: 1000;
    box-shadow: 0 4px 12px rgba(0,0,0,0.5);
}
.symbol-search-results.open {
    display: block;
}
.symbol-result-item {
    padding: 8px 12px;
    cursor: pointer;
    font-family: var(--font-mono);
    font-size: 13px;
    color: var(--text-main);
    display: flex;
    justify-content: space-between;
}
.symbol-result-item:hover {
    background: var(--bg-hover);
}
.symbol-result-item .vol {
    color: var(--text-muted);
    font-size: 11px;
}

/* Market Selector Button */
.market-selector-btn {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    background: var(--bg-hover);
    border: 1px solid var(--border-color);
    color: var(--text-main);
    padding: 7px 14px;
    border-radius: 6px;
    font-family: var(--font-ui);
    font-size: 14px;
    font-weight: 700;
    cursor: pointer;
    margin: 0;
    flex-shrink: 0;
    transition: var(--transition);
    letter-spacing: 0.3px;
    white-space: nowrap;
}
.market-selector-btn:hover {
    border-color: var(--color-accent);
    color: var(--text-highlight);
    background: var(--bg-panel);
}
.market-selector-btn svg { color: var(--text-muted); flex-shrink: 0; }

/* Market Selector Modal */
.market-modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.6);
    z-index: 9000;
    display: flex;
    align-items: flex-start;
    justify-content: center;
    padding-top: 70px;
    backdrop-filter: blur(2px);
}
.market-modal {
    background: var(--bg-panel);
    border: 1px solid var(--border-color);
    border-radius: 8px;
    width: 560px;
    max-height: 75vh;
    display: flex;
    flex-direction: column;
    box-shadow: 0 20px 60px rgba(0,0,0,0.7);
    animation: modalSlideIn 0.15s ease;
}
@keyframes modalSlideIn {
    from { opacity: 0; transform: translateY(-10px); }
    to   { opacity: 1; transform: translateY(0); }
}
.market-modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 20px;
    border-bottom: 1px solid var(--border-color);
}
.market-modal-title {
    font-size: 15px;
    font-weight: 600;
    color: var(--text-highlight);
}
.market-modal-close {
    background: none;
    border: none;
    color: var(--text-muted);
    font-size: 16px;
    cursor: pointer;
    padding: 4px 8px;
    border-radius: 4px;
    transition: var(--transition);
}
.market-modal-close:hover { background: var(--bg-hover); color: var(--text-main); }
.market-modal-search {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 12px 20px;
    border-bottom: 1px solid var(--border-color);
}
#modalSearchInput {
    flex: 1;
    background: none;
    border: none;
    color: var(--text-main);
    font-family: var(--font-ui);
    font-size: 14px;
    outline: none;
}
#modalSearchInput::placeholder { color: var(--text-muted); }
.market-modal-tabs {
    display: flex;
    gap: 4px;
    padding: 10px 20px;
    border-bottom: 1px solid var(--border-color);
}
.market-tab {
    background: none;
    border: 1px solid transparent;
    color: var(--text-muted);
    padding: 5px 16px;
    border-radius: 4px;
    font-family: var(--font-ui);
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    transition: var(--transition);
}
.market-tab:hover { color: var(--text-main); }
.market-tab.active {
    color: var(--color-accent);
    border-color: var(--border-color);
    background: var(--bg-hover);
}
.market-modal-cols {
    display: grid;
    grid-template-columns: 2fr 1fr 1fr;
    padding: 6px 20px;
    font-size: 11px;
    color: var(--text-muted);
    border-bottom: 1px solid var(--border-color);
}
.market-modal-cols span:not(:first-child) { text-align: right; }
.market-modal-list {
    flex: 1;
    overflow-y: auto;
}
.market-modal-loading {
    text-align: center;
    color: var(--text-muted);
    padding: 40px;
    font-size: 13px;
}
.market-item {
    display: grid;
    grid-template-columns: 2fr 1fr 1fr;
    align-items: center;
    padding: 9px 20px;
    cursor: pointer;
    transition: background 0.1s;
}
.market-item:hover { background: var(--bg-hover); }
.market-item.active-pair { background: var(--bg-hover); }
.market-item .pair-name {
    font-size: 13px;
    font-weight: 600;
    color: var(--text-highlight);
}
.market-item .pair-quote {
    font-size: 11px;
    color: var(--text-muted);
    margin-left: 2px;
    font-weight: 400;
}
.market-item .pair-price {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--text-main);
    text-align: right;
}
.market-item .pair-change {
    font-family: var(--font-mono);
    font-size: 12px;
    text-align: right;
}
.pair-change.up { color: var(--color-buy); }
.pair-change.down { color: var(--color-sell); }
`
}
