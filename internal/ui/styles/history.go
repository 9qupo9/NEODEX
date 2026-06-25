package styles

func History() string {
	return `
/* ================= HISTORY ================= */
/* Делаем панель истории более компактной */
.history-panel { height: 200px; background: var(--bg-panel); border-top: 1px solid var(--border-color); display: flex; flex-direction: column;}
.history-tabs { display: flex; padding: 0 20px; border-bottom: 1px solid var(--border-color); }
.history-tabs .tab { flex: none; padding: 10px 20px; }

.table-container { flex: 1; overflow-y: auto; }
table { width: 100%; border-collapse: collapse; text-align: left; font-size: 12px; }
th { color: var(--text-muted); font-weight: normal; padding: 8px 20px; position: sticky; top: 0; background: var(--bg-panel); z-index: 10;}
td { padding: 6px 20px; border-bottom: 1px solid var(--border-color); font-family: var(--font-mono); }
tr:hover td { background: var(--bg-hover); }
.empty-state { padding: 30px; text-align: center; color: var(--text-muted); }
`
}
