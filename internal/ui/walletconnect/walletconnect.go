package walletconnect

// IWalletConnectService определяет методы для интеграции логики подключения кошелька.
// Принцип SOLID: Dependency Inversion - UI будет зависеть от абстракции, а не от конкретной реализации.
type IWalletConnectService interface {
	GetScript() string
	GetButtonHTML() string
}

// WalletConnectService реализует IWalletConnectService.
// Принцип SOLID: Single Responsibility - класс отвечает только за подключение кошелька.
type WalletConnectService struct{}

// NewWalletConnectService создает новый инстанс сервиса подключения кошелька.
func NewWalletConnectService() IWalletConnectService {
	return &WalletConnectService{}
}

// GetScript возвращает JS логику для подключения Web3 кошелька.
func (s *WalletConnectService) GetScript() string {
	return `
// Глобальное состояние кошелька
let isWalletConnected = false;
let walletAddress = "";

// Функция для обновления интерфейса кошелька
function updateWalletUI() {
    const btn = document.getElementById('headerConnectWalletBtn');
    const textSpan = document.getElementById('walletBtnText');
    if (btn && textSpan) {
        if (isWalletConnected) {
            textSpan.innerText = walletAddress.substring(0, 6) + '...' + walletAddress.substring(walletAddress.length - 4);
            btn.classList.add('connected');
        } else {
            textSpan.innerText = 'Connect Wallet';
            btn.classList.remove('connected');
        }
    }
    
    if (typeof updateUI === 'function') updateUI();
    if (typeof fetchBalance === 'function') fetchBalance();
    if (typeof fetchOrders === 'function') fetchOrders();
    if (typeof fetchPositions === 'function') fetchPositions();
}

// При загрузке страницы проверяем, был ли кошелек подключен ранее
window.addEventListener('DOMContentLoaded', async () => {
    const savedWallet = localStorage.getItem('neodex_wallet_connected');
    if (savedWallet === 'true' && typeof window.ethereum !== 'undefined') {
        try {
            // Тихо проверяем доступные аккаунты без вызова модалки
            const accounts = await window.ethereum.request({ method: 'eth_accounts' });
            if (accounts && accounts.length > 0) {
                walletAddress = accounts[0];
                isWalletConnected = true;
                updateWalletUI();
            } else {
                localStorage.removeItem('neodex_wallet_connected');
            }
        } catch (e) {
            console.error("Ошибка авто-подключения кошелька", e);
        }
    }
});

// Логика подключения кошелька (Wallet Connect)
async function connectWallet() {
    if (typeof window.ethereum === 'undefined') {
        alert('Пожалуйста, установите MetaMask или другой Web3 кошелек.');
        return;
    }
    
    if (isWalletConnected) {
        isWalletConnected = false;
        walletAddress = "";
        localStorage.removeItem('neodex_wallet_connected');
        updateWalletUI();
        return;
    }
    
    try {
        // Заставляем MetaMask всегда спрашивать подтверждение и выбор аккаунта
        await window.ethereum.request({
            method: 'wallet_requestPermissions',
            params: [{ eth_accounts: {} }]
        });
        
        const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
        if (accounts && accounts.length > 0) {
            walletAddress = accounts[0];
            isWalletConnected = true;
            localStorage.setItem('neodex_wallet_connected', 'true');
            updateWalletUI();
        }
    } catch (e) {
        console.error('Ошибка подключения кошелька:', e);
    }
}
`
}

// GetButtonHTML возвращает HTML код кнопки для шапки сайта.
func (s *WalletConnectService) GetButtonHTML() string {
	return `
<style>
.wallet-connect-btn {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    background: linear-gradient(135deg, rgba(156, 163, 175, 0.1) 0%, rgba(156, 163, 175, 0.05) 100%);
    border: 1px solid var(--border-color);
    color: var(--text-highlight);
    padding: 7px 16px;
    border-radius: 8px;
    font-family: var(--font-ui);
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    transition: var(--transition);
    letter-spacing: 0.3px;
    position: relative;
    overflow: hidden;
}

.wallet-connect-btn::before {
    content: '';
    position: absolute;
    top: 0; left: 0; right: 0; bottom: 0;
    background: linear-gradient(135deg, rgba(255,255,255,0.1), transparent);
    opacity: 0;
    transition: opacity 0.2s ease;
}

.wallet-connect-btn:hover {
    border-color: var(--color-accent);
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(0,0,0,0.3);
}

.wallet-connect-btn:hover::before {
    opacity: 1;
}

.wallet-connect-btn svg {
    color: var(--color-accent);
    transition: var(--transition);
}

.wallet-connect-btn.connected {
    background: var(--bg-hover);
    border-color: var(--color-buy);
    color: var(--color-buy);
}

.wallet-connect-btn.connected svg {
    color: var(--color-buy);
}

/* Пульсирующая точка для подключенного состояния */
.wallet-status-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background-color: var(--color-buy);
    display: none;
    box-shadow: 0 0 8px var(--color-buy);
}

.wallet-connect-btn.connected .wallet-status-dot {
    display: block;
    animation: pulseDot 2s infinite;
}

@keyframes pulseDot {
    0% { opacity: 1; box-shadow: 0 0 8px rgba(46, 189, 133, 0.8); }
    50% { opacity: 0.4; box-shadow: 0 0 2px rgba(46, 189, 133, 0.2); }
    100% { opacity: 1; box-shadow: 0 0 8px rgba(46, 189, 133, 0.8); }
}
</style>
<button id="headerConnectWalletBtn" class="wallet-connect-btn" style="margin-right:15px;" onclick="connectWallet()">
    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 12V8H6a2 2 0 0 1-2-2c0-1.1.9-2 2-2h12v4"/><path d="M4 6v12c0 1.1.9 2 2 2h14v-4"/><path d="M18 12a2 2 0 0 0-2 2c0 1.1.9 2 2 2h4v-4h-4z"/></svg>
    <span id="walletBtnText">Connect Wallet</span>
    <div class="wallet-status-dot"></div>
</button>
`
}
