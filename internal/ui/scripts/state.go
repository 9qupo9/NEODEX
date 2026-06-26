package scripts

func State() string {
	return `
let ws;

const asksContainer = document.getElementById('asksContainer');
const bidsContainer = document.getElementById('bidsContainer');
const marketTradesContainer = document.getElementById('marketTradesContainer');
const marketTradesTable = document.getElementById('marketTradesTable');
const marketTradesEmptyState = document.getElementById('marketTradesEmptyState');


const historyTabBtns = document.querySelectorAll('.history-tab-btn');
const historyTable = document.getElementById('historyTable');
const historyTableBody = document.getElementById('historyTableBody');
const historyEmptyState = document.getElementById('historyEmptyState');

// State
let availUSDT = 0;
let availBase = 0; // Using base instead of hardcoded BTC
let currentSide = 'BUY'; // BUY or SELL
let currentType = 'MARKET'; // LIMIT, MARKET, STOP_LIMIT
let currentMarketPrice = 0;
let currentTickSize = 0.01;
let currentSymbol = 'btcusdt';
let currentBase = 'BTC';
let currentQuote = 'USDT';
let allBinanceSymbols = []; // Store fetched symbols
let dynamicNewListings = [];

// UI Elements
const sideBtns = document.querySelectorAll('.side-btn');
const tabBtns = document.querySelectorAll('.tab');
const availBalanceEl = document.getElementById('availBalance');
const submitBtn = document.getElementById('submitOrderBtn');
const stopInputGroup = document.getElementById('stopInputGroup');
const priceInputGroup = document.getElementById('priceInputGroup');
const priceInput = document.getElementById('priceInput');
const stopInput = document.getElementById('stopInput');
const qtyInput = document.getElementById('qtyInput');
const totalEl = document.getElementById('totalUSDT');

const qtySlider = document.getElementById('qtySlider');
const sliderFill = document.getElementById('sliderFill');
const sliderContainer = document.querySelector('.slider-container');

`
}
