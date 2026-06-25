package styles

func Orderform() string {
	return `
/* ================= ORDER FORM ================= */
.order-form-panel { background: var(--bg-panel); }
.order-tabs {
    display: flex;
    border-bottom: 1px solid var(--border-color);
}
.tab {
    flex: 1;
    text-align: center;
    padding: 10px 0;
    cursor: pointer;
    color: var(--text-muted);
    font-weight: 500;
    border-bottom: 2px solid transparent;
    transition: var(--transition);
}
.tab.active { color: var(--color-accent); border-bottom-color: var(--color-accent); }

/* Quote Currency Tabs */
.quote-tabs {
    display: flex;
    gap: 2px;
    padding: 8px 12px 4px 12px;
    border-bottom: 1px solid var(--border-color);
}
.quote-tab {
    background: none;
    border: none;
    color: var(--text-muted);
    font-family: var(--font-ui);
    font-size: 12px;
    font-weight: 600;
    padding: 4px 10px;
    border-radius: 4px;
    cursor: pointer;
    transition: var(--transition);
    white-space: nowrap;
}
.quote-tab:hover { color: var(--text-main); background: var(--bg-hover); }
.quote-tab.active {
    color: var(--color-accent);
    background: var(--bg-hover);
}

.side-toggle {
    display: flex;
    background: var(--bg-dark);
    border-radius: 4px;
    padding: 2px;
}
.side-btn {
    flex: 1;
    text-align: center;
    padding: 8px 0;
    cursor: pointer;
    border-radius: 2px;
    color: var(--text-muted);
    font-weight: 600;
    transition: var(--transition);
}
.side-btn.active.buy { background: var(--color-buy); color: #fff; }
.side-btn.active.sell { background: var(--color-sell); color: #fff; }

.order-forms-wrapper { display: flex; flex-direction: column; padding: 12px; gap: 12px; flex: 1; overflow-y: auto;}
.balance-row { display: flex; justify-content: space-between; color: var(--text-muted); font-size: 12px; }

.input-group {
    background: var(--bg-dark);
    border: 1px solid var(--border-color);
    border-radius: 2px;
    display: flex;
    align-items: center;
    padding: 0 10px;
    height: 36px;
    transition: var(--transition);
}
.input-group:focus-within { border-color: var(--color-accent); }
.input-label { color: var(--text-muted); padding-right: 10px; white-space: nowrap; }

/* Hide browser default number arrows */
input[type=number]::-webkit-inner-spin-button, 
input[type=number]::-webkit-outer-spin-button { 
  -webkit-appearance: none; 
  margin: 0; 
}
input[type=number] { -moz-appearance: textfield; }

.input-group input {
    background: transparent;
    border: none;
    color: var(--text-main);
    flex: 1;
    text-align: right;
    font-family: var(--font-mono);
    outline: none;
    width: 100%;
}
.input-suffix { padding-left: 10px; color: var(--text-main); }

.slider-container { margin: 15px 0 25px 0; padding: 0 7px; }
.slider-track-wrap { position: relative; height: 16px; display: flex; align-items: center; }
.slider-bg { position: absolute; width: 100%; height: 4px; background: var(--border-color); border-radius: 2px; }
.slider-fill { position: absolute; width: 0%; height: 4px; background: var(--color-buy); border-radius: 2px; transition: width 0.1s, background 0.2s; z-index: 1; }
.slider-marks { position: absolute; width: 100%; display: flex; justify-content: space-between; z-index: 2; pointer-events: none; }
.slider-mark { width: 10px; height: 10px; background: var(--bg-panel); border: 2px solid var(--border-color); border-radius: 50%; pointer-events: auto; cursor: pointer; transition: var(--transition); transform: rotate(45deg); }
.slider-mark:hover { border-color: var(--color-accent); }
.slider-range { 
    -webkit-appearance: none; 
    width: 100%; 
    background: transparent; 
    position: absolute; 
    z-index: 3; 
    cursor: pointer;
    margin: 0;
}
.slider-range::-webkit-slider-thumb { 
    -webkit-appearance: none; 
    width: 16px; 
    height: 16px; 
    background: #fff; 
    border: 3px solid var(--color-buy); 
    border-radius: 50%; 
    cursor: grab;
    margin-top: -6px;
    box-shadow: 0 0 5px rgba(0,0,0,0.5);
    transition: border-color 0.2s;
}
.slider-range::-webkit-slider-runnable-track {
    width: 100%;
    height: 4px;
    background: transparent;
}
.slider-range:active::-webkit-slider-thumb { cursor: grabbing; transform: scale(1.1); }

/* Dynamic Slider Colors based on mode */
.slider-container.sell-mode .slider-fill { background: var(--color-sell); }
.slider-container.sell-mode .slider-range::-webkit-slider-thumb { border-color: var(--color-sell); }
.slider-container.buy-mode .slider-fill { background: var(--color-buy); }
.slider-container.buy-mode .slider-range::-webkit-slider-thumb { border-color: var(--color-buy); }

.btn {
    width: 100%; height: 36px; border: none; border-radius: 2px;
    font-weight: bold; cursor: pointer; font-size: 13px; color: #fff;
    transition: var(--transition);
}
.btn-buy { background: var(--color-buy); }
.btn-buy:hover { background: #0cba76; }
.btn-sell { background: var(--color-sell); }
.btn-sell:hover { background: #e03f54; }

.action-row { display: flex; gap: 8px; }

/* Slider % labels */
.slider-labels {
    display: flex;
    justify-content: space-between;
    margin-top: 6px;
    padding: 0 2px;
}
.slider-labels span {
    font-size: 10px;
    color: var(--text-muted);
    cursor: pointer;
    transition: color 0.15s;
}
.slider-labels span:hover { color: var(--color-accent); }

/* Account Balance Block */
.account-balance-block {
    border-top: 1px solid var(--border-color);
    padding: 14px 12px 12px 12px;
    display: flex;
    flex-direction: column;
    gap: 10px;
}
.acct-bal-header {
    display: flex;
    align-items: center;
    gap: 6px;
    color: var(--text-muted);
    font-size: 11px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}
.acct-bal-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
}
.acct-bal-coin {
    font-size: 12px;
    font-weight: 600;
    color: var(--text-muted);
}
.acct-bal-amount {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--text-main);
}
.acct-bal-total-row {
    border-top: 1px solid var(--border-color);
    padding-top: 8px;
    margin-top: 2px;
}
.acct-bal-total-row .acct-bal-amount {
    color: var(--color-accent);
    font-weight: 600;
}
`
}
