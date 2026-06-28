package admin

func GetSidebarCSS() string {
	return `
        /* Боковая панель */
        .sidebar {
            width: 260px;
            background-color: var(--bg-panel);
            border-right: 1px solid var(--border-color);
            display: flex;
            flex-direction: column;
            padding: 24px;
            box-sizing: border-box;
        }

        .logo {
            font-family: 'JetBrains Mono', monospace;
            font-size: 24px;
            font-weight: 600;
            color: var(--text-main);
            text-decoration: none;
            letter-spacing: -0.5px;
            margin-bottom: 48px;
            display: flex;
            align-items: center;
        }

        .nav-link {
            display: block;
            padding: 12px 16px;
            color: var(--text-muted);
            text-decoration: none;
            font-size: 14px;
            font-weight: 500;
            border-radius: 4px;
            margin-bottom: 8px;
            transition: all 0.2s;
            border: 1px solid transparent;
        }

        .nav-link.active {
            background-color: rgba(41, 98, 255, 0.1);
            color: var(--accent-color);
            border-color: rgba(41, 98, 255, 0.2);
        }

        .nav-link:hover:not(.active) {
            background-color: rgba(255, 255, 255, 0.05);
            color: var(--text-main);
        }

        .sidebar-bottom {
            margin-top: auto;
        }

        .logout-btn {
            display: block;
            width: 100%;
            padding: 12px;
            background: transparent;
            border: 1px solid var(--border-color);
            color: var(--text-muted);
            text-align: center;
            font-size: 13px;
            border-radius: 4px;
            cursor: pointer;
            transition: all 0.2s;
            text-decoration: none;
            box-sizing: border-box;
        }

        .logout-btn:hover {
            border-color: var(--text-main);
            color: var(--text-main);
        }

        /* Премиальный дизайн кнопок управления в сайдбаре */
        .sidebar-controls {
            margin-top: 32px;
            display: flex;
            flex-direction: column;
            gap: 12px;
            padding-top: 24px;
            border-top: 1px solid rgba(255,255,255,0.05);
        }

        .side-heading {
            font-size: 11px;
            color: #999;
            text-transform: uppercase;
            margin-bottom: 4px;
            letter-spacing: 1.5px;
            font-weight: 600;
        }

        .btn-side {
            background: rgba(255, 255, 255, 0.06);
            border: 1px solid rgba(255, 255, 255, 0.1);
            color: #E2E8F0;
            padding: 12px 14px;
            font-size: 13px;
            font-family: 'Inter', sans-serif;
            font-weight: 500;
            border-radius: 6px;
            cursor: pointer;
            transition: all 0.2s ease;
            display: flex;
            align-items: center;
            justify-content: flex-start;
            gap: 12px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.2);
        }

        .btn-side:hover {
            background: rgba(255, 255, 255, 0.12);
            border-color: rgba(255, 255, 255, 0.2);
            color: #ffffff;
            transform: translateY(-1px);
            box-shadow: 0 4px 8px rgba(0,0,0,0.3);
        }

        .led-indicator {
            width: 10px;
            height: 10px;
            border-radius: 50%;
            background-color: #333;
            box-shadow: inset 0 1px 2px rgba(0,0,0,0.5), 0 0 2px #111;
            transition: all 0.3s ease;
        }

        /* Состояния LED-индикаторов */
        .btn-start.active {
            color: #00FF66;
            border-color: rgba(0, 255, 102, 0.3);
            background: rgba(0, 255, 102, 0.08);
            text-shadow: 0 0 8px rgba(0, 255, 102, 0.4);
        }
        .btn-start.active .led-indicator {
            background-color: #00FF66;
            box-shadow: 0 0 10px #00FF66, 0 0 20px #00FF66;
        }
        .btn-start:disabled {
            opacity: 0.3;
            cursor: not-allowed;
            background: transparent;
            border-color: rgba(255,255,255,0.05);
            color: #666;
            transform: none;
            box-shadow: none;
        }

        .btn-stop.active {
            color: #FF0033;
            border-color: rgba(255, 0, 51, 0.3);
            background: rgba(255, 0, 51, 0.08);
            text-shadow: 0 0 8px rgba(255, 0, 51, 0.4);
        }
        .btn-stop.active .led-indicator {
            background-color: #FF0033;
            box-shadow: 0 0 10px #FF0033, 0 0 20px #FF0033;
        }
        .btn-stop:disabled {
            opacity: 0.3;
            cursor: not-allowed;
            background: transparent;
            border-color: rgba(255,255,255,0.05);
            color: #666;
            transform: none;
            box-shadow: none;
        }
    `
}

func RenderSidebar() string {
	return `
    <aside class="sidebar">
        <a href="#" class="logo">
            NEODEX
        </a>
        <nav>
            <a href="#" class="nav-link active" onclick="switchTab('overview')">System Overview</a>
            <a href="#" class="nav-link" onclick="switchTab('users')">Users</a>
        </nav>
        
        <div class="sidebar-controls">
            <div class="side-heading">System Commands</div>
            <button class="btn-side" onclick="sendAction('snapshot')">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path><polyline points="17 21 17 13 7 13 7 21"></polyline><polyline points="7 3 7 8 15 8"></polyline></svg>
                Create Snapshot
            </button>
            <button class="btn-side" onclick="sendAction('clear_cache')">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
                Clear Logs
            </button>
            
            <div class="side-heading" style="margin-top: 16px;">Exchange Engine</div>
            <button id="btn-start" class="btn-side btn-start active" onclick="startTrading()">
                <div class="led-indicator"></div>
                Start Trading
            </button>
            <button id="btn-stop" class="btn-side btn-stop" onclick="stopTrading()" disabled>
                <div class="led-indicator"></div>
                Stop Trading
            </button>
        </div>

        <div class="sidebar-bottom">
            <a href="/" class="logout-btn" onclick="document.cookie='admin_session=; Max-Age=0; path=/;'">Logout</a>
        </div>
    </aside>`
}
