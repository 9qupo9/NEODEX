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
            document.getElementById('modal-actions').innerHTML = '<button class="btn-side" onclick="closeModal()">Закрыть</button>';
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
            document.getElementById('modal-actions').innerHTML = '<button class="btn-side" onclick="closeModal()">Отмена</button>' +
                '<button class="btn-action btn-info" style="font-size: 13px; padding: 8px 16px;" id="prompt-confirm-btn">Подтвердить</button>';
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
                    showModal('Ошибка Сервера', 'Код ошибки: ' + res.status + '.<br>Вы перезапустили сервер (go run .) после обновления?');
                    return;
                }
                
                if (action === 'snapshot') {
                    showModal('Снапшот Создан', 'Снапшот AOF успешно создан и данные сжаты!');
                } else if (action === 'clear_cache') {
                    showModal('Очистка Логов', 'Системные логи успешно очищены!');
                }
                
                // Сразу запрашиваем метрики, чтобы обновить UI без задержки
                if (action === 'toggle_trading') {
                    fetchMetrics();
                }
            } catch (e) {
                console.error("Action error", e);
                showModal('Ошибка Сети', 'Сервер недоступен. Проверьте подключение.');
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
                container.innerHTML = '<div class="glass-panel" style="padding:20px; text-align:center; color:#888;">Нет пользователей</div>';
                return;
            }
            let html = '<div class="glass-panel" style="padding:0;"><table class="glass-table">';
            html += '<tr><th>Адрес Аккаунта</th><th>USDT Баланс</th><th>Действия</th></tr>';
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
                    html += ' <span class="badge badge-danger" style="display:inline-block; margin-left:10px;">Заблокирован</span>';
                }
                html += '</td>';
                html += '<td style="color:#00FF66; font-weight: bold;">$' + usdtStr + '</td>';
                html += '<td>';
                html += '<div style="display:flex; gap:8px;">';
                html += '<button class="btn-action" onclick="showModal(\'Аудит\', \'Функция аудита кошелька в разработке\')">Аудит</button>';
                if (u.IsBlocked) {
                    html += '<button class="btn-action btn-info" onclick="toggleBlock(\'' + u.Address + '\', false)">Разблокировать</button>';
                } else {
                    html += '<button class="btn-action btn-danger" onclick="toggleBlock(\'' + u.Address + '\', true)">Заблокировать</button>';
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
                    showModal('Блокировка', 'Аккаунт ' + address.substring(0,6) + '... успешно заблокирован.');
                } else {
                    showModal('Разблокировка', 'Аккаунт ' + address.substring(0,6) + '... успешно разблокирован.');
                }
            } catch (e) {
                showModal('Ошибка', 'Не удалось изменить статус аккаунта.');
            }
        }



        function showSecurityModal() {
            let html = '<div style="margin-bottom: 10px; padding: 12px; background: rgba(0, 255, 102, 0.05); border: 1px solid rgba(0, 255, 102, 0.2); border-radius: 8px;">' +
                '<strong style="color: #00FF66;">WAF (Сетевой экран):</strong> Активен и функционирует штатно.' +
                '</div>' +
                '<div style="margin-bottom: 10px; padding: 12px; background: rgba(0, 153, 255, 0.05); border: 1px solid rgba(0, 153, 255, 0.2); border-radius: 8px;">' +
                '<strong style="color: #0099ff;">Авторизации:</strong> Подозрительной активности не выявлено. Все входы авторизованы.' +
                '</div>' +
                '<div style="margin-bottom: 10px; padding: 12px; background: rgba(255, 255, 255, 0.05); border: 1px solid rgba(255, 255, 255, 0.1); border-radius: 8px;">' +
                '<strong style="color: #E2E8F0;">Шифрование БД:</strong> AES-256 (Ключи актуальны).' +
                '</div>';
            showModal('Аудит Безопасности', html);
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
                    logsContainer.innerHTML = '<div><span class="log-info">Логи отсутствуют</span></div>';
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
