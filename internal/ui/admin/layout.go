package admin

import (
	"dex/internal/ui/styles"
)

// RenderLayout собирает и возвращает полную HTML страницу админ-панели
func RenderLayout() string {
	return `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>NEODEX | Admin</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
    <script src="https://unpkg.com/lightweight-charts/dist/lightweight-charts.standalone.production.js"></script>
    <style>
        ` + styles.Vars() + `
        
        :root {
            --bg-base: #000000;
            --bg-panel: #0A0A0A;
            --border-color: #222222;
            --text-main: #FFFFFF;
            --text-muted: #888888;
            --accent-color: #2962FF;
            --success-color: #00C853;
            --danger-color: #D50000;
            --warning-color: #FFAB00;
        }

        body {
            margin: 0;
            background-color: var(--bg-base);
            color: var(--text-main);
            font-family: 'Inter', -apple-system, sans-serif;
            -webkit-font-smoothing: antialiased;
            display: flex;
            height: 100vh;
            overflow: hidden;
        }

        /* Главная область */
        .main-content {
            flex: 1;
            display: flex;
            flex-direction: column;
            overflow-y: auto;
            padding: 40px;
        }

        .header {
            display: flex;
            justify-content: space-between;
            align-items: flex-end;
            margin-bottom: 32px;
        }

        h1 {
            margin: 0;
            font-size: 28px;
            font-weight: 500;
            letter-spacing: -0.5px;
        }

        .sys-status {
            font-family: 'JetBrains Mono', monospace;
            font-size: 12px;
            color: var(--success-color);
            display: flex;
            align-items: center;
            gap: 8px;
        }

        .sys-status::before {
            content: '';
            display: block;
            width: 6px;
            height: 6px;
            background-color: var(--success-color);
            border-radius: 50%;
        }

        /* Сетка дашборда */
        .dashboard-grid {
            display: grid;
            grid-template-columns: repeat(12, 1fr);
            gap: 24px;
        }

        .panel {
            background-color: var(--bg-panel);
            border: 1px solid var(--border-color);
            border-radius: 4px;
            padding: 24px;
            display: flex;
            flex-direction: column;
        }

        .panel-title {
            font-size: 12px;
            text-transform: uppercase;
            letter-spacing: 1px;
            color: var(--text-muted);
            margin-bottom: 16px;
            font-weight: 600;
        }

        /* Скроллбары */
        ::-webkit-scrollbar { width: 4px; }
        ::-webkit-scrollbar-track { background: transparent; }
        ::-webkit-scrollbar-thumb { background: #333; }
        ::-webkit-scrollbar-thumb:hover { background: #555; }
    </style>
    
    <!-- Компонентные стили -->
    <style>
        ` + GetSidebarCSS() + `
        ` + GetMetricsCSS() + `
        ` + GetChartCSS() + `
        ` + GetLogsCSS() + `
        
        /* Стили для табов */
        .tab-content {
            display: none;
            animation: fadeIn 0.3s ease;
        }
        .tab-content.active {
            display: block;
        }
        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(10px); }
            to { opacity: 1; transform: translateY(0); }
        }

        /* Премиум Таблицы */
        .glass-table {
            width: 100%;
            border-collapse: separate;
            border-spacing: 0;
            text-align: left;
            margin-top: 0;
        }
        .glass-table th {
            padding: 16px 20px;
            color: #8892b0;
            font-size: 11px;
            text-transform: uppercase;
            letter-spacing: 1px;
            font-weight: 600;
            border-bottom: 1px solid rgba(255,255,255,0.05);
            background: rgba(0,0,0,0.2);
        }
        .glass-table td {
            padding: 16px 20px;
            color: #e2e8f0;
            font-size: 13px;
            border-bottom: 1px solid rgba(255,255,255,0.03);
            transition: background 0.2s;
        }
        .glass-table tr:hover td {
            background: rgba(255,255,255,0.02);
        }

        /* Премиум Бейджи */
        .badge {
            padding: 6px 12px;
            border-radius: 20px;
            font-size: 10px;
            font-weight: bold;
            text-transform: uppercase;
            letter-spacing: 1px;
            display: flex;
            align-items: center;
            gap: 6px;
            box-shadow: 0 0 10px rgba(0,0,0,0.2);
            background: rgba(255,255,255,0.05);
            border: 1px solid rgba(255,255,255,0.1);
            color: #ccc;
        }
        .badge::before {
            content: '';
            display: inline-block;
            width: 6px;
            height: 6px;
            border-radius: 50%;
            background: #888;
        }
        .badge-info { background: rgba(0, 153, 255, 0.1); color: #0099ff; border-color: rgba(0, 153, 255, 0.2); }
        .badge-info::before { background: #0099ff; box-shadow: 0 0 5px #0099ff; }

        .badge-warning { background: rgba(255, 153, 0, 0.1); color: #ff9900; border-color: rgba(255, 153, 0, 0.2); }
        .badge-warning::before { background: #ff9900; box-shadow: 0 0 5px #ff9900; }

        .badge-purple { background: rgba(200, 0, 255, 0.1); color: #c800ff; border-color: rgba(200, 0, 255, 0.2); }
        .badge-purple::before { background: #c800ff; box-shadow: 0 0 5px #c800ff; }

        .badge-success { background: rgba(0, 255, 102, 0.1); color: #00FF66; border-color: rgba(0, 255, 102, 0.2); }
        .badge-success::before { background: #00FF66; box-shadow: 0 0 5px #00FF66; }

        .badge-danger { background: rgba(255, 0, 51, 0.1); color: #FF0033; border-color: rgba(255, 0, 51, 0.2); }
        .badge-danger::before { background: #FF0033; box-shadow: 0 0 5px #FF0033; }


        /* Карточки Безопасности */
        .security-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 20px;
        }
        .sec-card {
            background: rgba(255, 255, 255, 0.02);
            border: 1px solid rgba(255, 255, 255, 0.05);
            border-radius: 12px;
            padding: 24px;
            display: flex;
            flex-direction: column;
            gap: 16px;
            transition: all 0.3s ease;
        }
        .sec-card:hover {
            background: rgba(255, 255, 255, 0.04);
            border-color: rgba(255, 255, 255, 0.1);
            transform: translateY(-2px);
        }
        .sec-title {
            font-size: 12px;
            color: #8892b0;
            text-transform: uppercase;
            letter-spacing: 1px;
            font-weight: 600;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .sec-val {
            font-size: 24px;
            font-weight: 700;
            color: #E2E8F0;
            font-family: monospace;
        }

        /* Кнопки действий (Таблицы) */
        .btn-action {
            background: rgba(255, 255, 255, 0.03);
            border: 1px solid rgba(255, 255, 255, 0.1);
            color: #aaa;
            padding: 6px 14px;
            border-radius: 6px;
            font-size: 11px;
            font-weight: 600;
            text-transform: uppercase;
            cursor: pointer;
            transition: all 0.3s ease;
        }
        .btn-action:hover {
            background: rgba(255, 255, 255, 0.08);
            transform: translateY(-1px);
        }
        .btn-action.btn-info {
            color: #0099ff;
            border-color: rgba(0, 153, 255, 0.3);
            background: rgba(0, 153, 255, 0.05);
        }
        .btn-action.btn-info:hover {
            background: rgba(0, 153, 255, 0.15);
            box-shadow: 0 0 10px rgba(0, 153, 255, 0.2);
        }
        .btn-action.btn-danger {
            color: #FF0033;
            border-color: rgba(255, 0, 51, 0.3);
            background: rgba(255, 0, 51, 0.05);
        }
        .btn-action.btn-danger:hover {
            background: rgba(255, 0, 51, 0.15);
            box-shadow: 0 0 10px rgba(255, 0, 51, 0.2);
        }


        /* Модальные окна */
        .modal-overlay {
            position: fixed;
            top: 0; left: 0; right: 0; bottom: 0;
            background: rgba(0, 0, 0, 0.7);
            backdrop-filter: blur(5px);
            display: none;
            align-items: center;
            justify-content: center;
            z-index: 1000;
            opacity: 0;
            transition: opacity 0.3s ease;
        }
        .modal-overlay.active {
            display: flex;
            opacity: 1;
        }
        .modal-container {
            background: #111;
            border: 1px solid rgba(255,255,255,0.1);
            border-radius: 12px;
            width: 400px;
            max-width: 90%;
            box-shadow: 0 10px 30px rgba(0,0,0,0.5);
            transform: translateY(20px);
            transition: transform 0.3s ease;
            overflow: hidden;
        }
        .modal-overlay.active .modal-container {
            transform: translateY(0);
        }
        .modal-header {
            padding: 16px 20px;
            border-bottom: 1px solid rgba(255,255,255,0.05);
            font-weight: 600;
            font-size: 14px;
            color: #E2E8F0;
            background: rgba(255,255,255,0.02);
        }
        .modal-body {
            padding: 20px;
            color: #aaa;
            font-size: 13px;
            line-height: 1.5;
        }
        .modal-footer {
            padding: 16px 20px;
            border-top: 1px solid rgba(255,255,255,0.05);
            display: flex;
            justify-content: flex-end;
            gap: 10px;
            background: rgba(255,255,255,0.01);
        }
        .modal-input {
            width: 100%;
            background: rgba(0,0,0,0.5);
            border: 1px solid rgba(255,255,255,0.1);
            color: #fff;
            padding: 10px;
            border-radius: 6px;
            margin-top: 8px;
            font-family: monospace;
            outline: none;
        }
        .modal-input:focus {
            border-color: rgba(255, 153, 0, 0.5);
            box-shadow: 0 0 5px rgba(255, 153, 0, 0.2);
        }
    </style>
</head>
<body>
    ` + RenderSidebar() + `

    <!-- Глобальное Модальное Окно -->
    <div id="global-modal" class="modal-overlay">
        <div class="modal-container">
            <div id="modal-title" class="modal-header">Уведомление</div>
            <div id="modal-content" class="modal-body"></div>
            <div id="modal-actions" class="modal-footer">
                <button class="btn-side" onclick="closeModal()">Закрыть</button>
            </div>
        </div>
    </div>

    <main class="main-content">
        <div id="tab-overview" class="tab-content active">
            <header class="header">
                <h1>Системный Обзор</h1>
                <div style="display:flex; gap:16px; align-items:center;">
                    <div class="badge badge-success">ALL SYSTEMS OPERATIONAL</div>
                    <button class="btn-action btn-info" onclick="showSecurityModal()">АУДИТ БЕЗОПАСНОСТИ</button>
                </div>
            </header>

            <div class="dashboard-grid">
                ` + RenderMetrics() + `
                ` + RenderChart() + `
                ` + RenderLogs() + `
            </div>
        </div>

        <div id="tab-users" class="tab-content">
            <header class="header">
                <h1>Пользователи</h1>
                <div class="badge badge-info">УПРАВЛЕНИЕ АККАУНТАМИ</div>
            </header>
            <div id="users-container"></div>
        </div>

    </main>

    ` + RenderScripts() + `
</body>
</html>`
}
