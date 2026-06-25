package styles

func Layout() string {
	return `
.app-container {
    display: flex;
    flex-direction: column;
    height: 100%;
}

.main-layout {
    display: grid;
    /* 3 columns: left sidebar, center chart, right sidebar */
    grid-template-columns: 280px 1fr 280px;
    grid-template-rows: 1fr;
    flex: 1;
    overflow: hidden;
    gap: 1px;
    background: var(--border-color);
}
.panel { background: var(--bg-dark); display: flex; flex-direction: column; overflow: hidden; }
.chart-panel { display: flex; flex-direction: column; }

`
}
