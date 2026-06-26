package styles

func Vars() string {
	return `
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=Roboto+Mono:wght@400;500;700&display=swap');

:root {
    /* Premium Dark Chrome / Titanium Theme */
    --bg-dark: #0d0e12;
    --bg-panel: #14151a;
    --bg-hover: #1e1f26;
    --border-color: #262831;
    
    --text-main: #d1d5db;
    --text-muted: #6b7280;
    --text-highlight: #ffffff;
    
    /* Muted, non-acidic Buy/Sell colors */
    --color-buy: #2ebd85;
    --color-buy-bg: rgba(46, 189, 133, 0.1);
    --color-sell: #e0294a;
    --color-sell-bg: rgba(224, 41, 74, 0.1);
    
    /* Chrome / Silver Accent */
    --color-accent: #9ca3af;
    --color-accent-hover: #e5e7eb;
    
    /* Primary / Brand color */
    --primary-color: #3468d0;
    
    --font-ui: 'Inter', sans-serif;
    --font-mono: 'Roboto Mono', monospace;
    
    --transition: all 0.2s ease;
}

* { box-sizing: border-box; margin: 0; padding: 0; }
::-webkit-scrollbar { width: 4px; height: 4px; }
::-webkit-scrollbar-track { background: var(--bg-dark); }
::-webkit-scrollbar-thumb { background: var(--bg-hover); border-radius: 4px; }
::-webkit-scrollbar-thumb:hover { background: var(--color-accent); }

body {
    font-family: var(--font-ui);
    background-color: var(--bg-dark);
    color: var(--text-main);
    overflow: hidden;
    height: 100vh;
    font-size: 13px;
    -webkit-font-smoothing: antialiased;
}
`
}
