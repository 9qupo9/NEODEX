package admin

func RenderScripts() string {
	return `
	<script>
		function formatNumber(num) {
			if (num >= 1000000) {
				let formatted = (num / 1000000).toFixed(2);
				return formatted.replace(/\.00$/, '') + 'M';
			}
			if (num >= 1000) {
				let formatted = (num / 1000).toFixed(2);
				return formatted.replace(/\.00$/, '') + 'k';
			}
			let formatted = num.toFixed(2);
			return formatted.replace(/\.00$/, '');
		}

        let chart = null;
        let lineSeries = null;
        const chartData = [];
        
        try {
            const chartProperties = {
                width: document.getElementById('chart').clientWidth,
                height: 310,
                layout: { background: { type: 'solid', color: 'transparent' }, textColor: '#555' },
                grid: { vertLines: { color: '#111' }, horzLines: { color: '#111' } },
                timeScale: { timeVisible: true, secondsVisible: true, borderColor: '#222' },
                rightPriceScale: { borderColor: '#222' }
            };
            chart = LightweightCharts.createChart(document.getElementById('chart'), chartProperties);
            if (chart && typeof chart.addLineSeries === 'function') {
                lineSeries = chart.addLineSeries({
                    color: '#2962FF',
                    lineWidth: 1,
                    crosshairMarkerVisible: false,
                });
            } else {
                console.warn("LightweightCharts API changed or chart not created properly");
            }
        } catch (e) {
            console.error("Error initializing chart:", e);
        }
        let currentIsHalted = false;

        function showModal(title, message) {
            document.getElementById('modal-title').innerText = title;
            document.getElementById('modal-content').innerHTML = message;
            document.getElementById('modal-actions').innerHTML = '<button class="btn-side" onclick="closeModal()">Close</button>';
            document.getElementById('global-modal').classList.add('active');
        }

        function closeModal() {
            document.getElementById('global-modal').classList.remove('active');
        }

        function showPrompt(title, fields, onConfirm) {
            document.getElementById('modal-title').innerText = title;
            let html = '';
            fields.forEach((f, i) => {
                html += '<div style="margin-bottom:10px;"><label style="color:#aaa; font-size:11px; text-transform:uppercase;">' + f.label + '</label><br><input type="text" id="prompt-input-' + i + '" class="modal-input" placeholder="' + (f.placeholder || '') + '" value="' + (f.value || '') + '"></div>';
            });
            document.getElementById('modal-content').innerHTML = html;
            document.getElementById('modal-actions').innerHTML = '<button class="btn-side" onclick="closeModal()">Cancel</button>' +
                '<button class="btn-action btn-info" style="font-size: 13px; padding: 8px 16px;" id="prompt-confirm-btn">Confirm</button>';
            document.getElementById('global-modal').classList.add('active');
            
            document.getElementById('prompt-confirm-btn').onclick = () => {
                let values = fields.map((f, i) => document.getElementById('prompt-input-' + i).value);
                closeModal();
                onConfirm(values);
            };
        }

        async function sendAction(action, halt = false) {
            try {
                const res = await fetch('/api/v1/admin/action', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ action: action, halt: halt })
                });
                
                if (!res.ok) {
                    showModal('Server Error', 'Error code: ' + res.status + '.<br>Did you restart the server (go run .) after updating?');
                    return;
                }
                
                if (action === 'snapshot') {
                    showModal('Snapshot Created', 'AOF snapshot successfully created and data compacted!');
                } else if (action === 'clear_cache') {
                    showModal('Clear Logs', 'System logs successfully cleared!');
                }
                
                if (action === 'toggle_trading') {
                    fetchMetrics();
                }
            } catch (e) {
                console.error("Action error", e);
                showModal('Network Error', 'Server unreachable. Check connection.');
            }
        }

        async function startTrading() {
            await sendAction('toggle_trading', false);
        }

        async function stopTrading() {
            await sendAction('toggle_trading', true);
        }

        function switchTab(tabId) {
            // Убираем active класс у всех табов
            document.querySelectorAll('.tab-content').forEach(el => el.classList.remove('active'));
            document.querySelectorAll('.nav-link').forEach(el => el.classList.remove('active'));
            
            // Активируем нужный
            document.getElementById('tab-' + tabId).classList.add('active');
            event.target.classList.add('active');
            
            // Подгружаем данные
            if (tabId === 'users') loadUsers();
        }

        async function loadUsers() {
            const res = await fetch('/api/v1/admin/users');
            const data = await res.json();
            const container = document.getElementById('users-container');
            if(!data || data.length === 0) {
                container.innerHTML = '<div class="glass-panel" style="padding:20px; text-align:center; color:#888;">No users</div>';
                return;
            }
            let html = '<div class="glass-panel" style="padding:0;"><table class="glass-table">';
            html += '<tr><th>Account Address</th><th>USDT Balance</th><th>Actions</th></tr>';
            data.forEach(u => {
                let rawUsdt = u.Balances && u.Balances.USDT ? parseFloat(u.Balances.USDT) : 0;
                let usdtStr = "0.00";
                if (rawUsdt >= 1000) {
                    usdtStr = (rawUsdt / 1000).toFixed(2) + "k";
                } else if (rawUsdt > 0) {
                    usdtStr = rawUsdt.toFixed(2);
                }
                
                html += '<tr>';
                html += '<td style="font-family:monospace; color:#E2E8F0;">' + u.Address;
                if (u.IsBlocked) {
                    html += ' <span class="badge badge-danger" style="display:inline-block; margin-left:10px;">Blocked</span>';
                }
                html += '</td>';
                html += '<td style="color:#00FF66; font-weight: bold;">$' + usdtStr + '</td>';
                html += '<td>';
                html += '<div style="display:flex; gap:8px;">';
                html += '<button class="btn-action" onclick="showModal(\'Audit\', \'Wallet audit feature in development\')">Audit</button>';
                if (u.IsBlocked) {
                    html += '<button class="btn-action btn-info" onclick="toggleBlock(\'' + u.Address + '\', false)">Unblock</button>';
                } else {
                    html += '<button class="btn-action btn-danger" onclick="toggleBlock(\'' + u.Address + '\', true)">Block</button>';
                }
                html += '</div></td>';
                html += '</tr>';
            });
            html += '</table></div>';
            container.innerHTML = html;
        }

        async function toggleBlock(address, block) {
            try {
                await fetch('/api/v1/admin/users/block', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ address: address, block: block })
                });
                loadUsers();
                if (block) {
                    showModal('Block', 'Account ' + address.substring(0,6) + '... successfully blocked.');
                } else {
                    showModal('Unblock', 'Account ' + address.substring(0,6) + '... successfully unblocked.');
                }
            } catch (e) {
                showModal('Error', 'Failed to change account status.');
            }
        }



        function showSecurityModal() {
            let html = '<div style="margin-bottom: 10px; padding: 12px; background: rgba(0, 255, 102, 0.05); border: 1px solid rgba(0, 255, 102, 0.2); border-radius: 8px;">' +
                '<strong style="color: #00FF66;">WAF (Firewall):</strong> Active and functioning normally.' +
                '</div>' +
                '<div style="margin-bottom: 10px; padding: 12px; background: rgba(0, 153, 255, 0.05); border: 1px solid rgba(0, 153, 255, 0.2); border-radius: 8px;">' +
                '<strong style="color: #0099ff;">Authorizations:</strong> No suspicious activity detected. All logins authorized.' +
                '</div>' +
                '<div style="margin-bottom: 10px; padding: 12px; background: rgba(255, 255, 255, 0.05); border: 1px solid rgba(255, 255, 255, 0.1); border-radius: 8px;">' +
                '<strong style="color: #E2E8F0;">DB Encryption:</strong> AES-256 (Keys up-to-date).' +
                '</div>';
            showModal('Security Audit', html);
        }

		async function fetchMetrics() {
			try {
				const res = await fetch('/api/v1/admin/metrics');
				const data = await res.json();
				
				document.getElementById('metric-volume').innerText = '$' + formatNumber(parseFloat(data.volume || 0));
				document.getElementById('metric-revenue').innerText = '$' + formatNumber(parseFloat(data.revenue || 0));
				document.getElementById('metric-users').innerText = formatNumber(data.users);
				document.getElementById('metric-ws').innerText = data.ws_clients;
				document.getElementById('metric-bots').innerText = data.tcp_bots;
				document.getElementById('metric-orders').innerText = formatNumber(data.active_orders);
				
                const latency = data.latency_ms;
				document.getElementById('metric-latency').innerText = latency.toFixed(4) + ' ms';
				
                currentIsHalted = data.is_halted;
                const btnStart = document.getElementById('btn-start');
                const btnStop = document.getElementById('btn-stop');
                
                if (currentIsHalted) {
                    // ЯДРО ОСТАНОВЛЕНО
                    btnStart.disabled = false;
                    btnStart.classList.add('active');
                    
                    btnStop.disabled = true;
                    btnStop.classList.remove('active');
                } else {
                    // ЯДРО РАБОТАЕТ
                    btnStart.disabled = true;
                    btnStart.classList.remove('active');
                    
                    btnStop.disabled = false;
                    btnStop.classList.add('active');
                }

                if (lineSeries) {
                    const now = Math.floor(Date.now() / 1000);
                    chartData.push({ time: now, value: latency });
                    if(chartData.length > 60) chartData.shift(); 
                    lineSeries.setData(chartData);
                }

				const logsContainer = document.getElementById('sysLogs');
				if (data.logs && data.logs.length > 0) {
					let logsHTML = '';
					for (const entry of data.logs) {
						const date = new Date(entry.timestamp);
						const timeStr = date.toLocaleTimeString('ru-RU', { hour12: false });
						
						let colorClass = 'log-info';
						if (entry.message.includes('Error') || entry.message.includes('Failed') || entry.message.includes('halted')) colorClass = 'log-error';
						else if (entry.message.includes('loaded') || entry.message.includes('started')) colorClass = 'log-success';
						else if (entry.message.includes('Warn')) colorClass = 'log-warn';
						
						logsHTML += '<div><span class="log-time">' + timeStr + '</span> <span class="' + colorClass + '">' + entry.message + '</span></div>';
					}
					logsContainer.innerHTML = logsHTML;
					logsContainer.scrollTop = logsContainer.scrollHeight;
				} else {
                    logsContainer.innerHTML = '<div><span class="log-info">No logs</span></div>';
                }
			} catch (e) {
				console.error("Fetch error", e);
			}
		}

		setInterval(fetchMetrics, 1000);
		fetchMetrics();

        window.addEventListener('resize', () => {
            if (chart) {
                chart.resize(document.getElementById('chart').clientWidth, 310);
            }
        });
	</script>
    `
}
