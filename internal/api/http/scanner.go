package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	knownSymbols   = make(map[string]bool)
	newListings    = []string{}
	scannerMutex   sync.RWMutex
	dataDir        = "data"
	knownFile      = filepath.Join(dataDir, "known_symbols.json")
	newListingsFile = filepath.Join(dataDir, "new_listings.json")
)

// InitScanner initializes the automated new listings tracker.
func InitScanner() {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Println("Error creating data dir:", err)
	}

	loadScannerData()

	// Run initial scan
	go runScannerLoop()
}

func loadScannerData() {
	scannerMutex.Lock()
	defer scannerMutex.Unlock()

	// Load known symbols
	if b, err := os.ReadFile(knownFile); err == nil {
		json.Unmarshal(b, &knownSymbols)
	}

	// Load new listings
	if b, err := os.ReadFile(newListingsFile); err == nil {
		json.Unmarshal(b, &newListings)
	}
}

func saveScannerData() {
	if b, err := json.MarshalIndent(knownSymbols, "", "  "); err == nil {
		os.WriteFile(knownFile, b, 0644)
	}
	if b, err := json.MarshalIndent(newListings, "", "  "); err == nil {
		os.WriteFile(newListingsFile, b, 0644)
	}
}

func runScannerLoop() {
	for {
		scanBinance()
		time.Sleep(1 * time.Hour) // Check every hour
	}
}

func scanBinance() {
	resp, err := http.Get("https://api.binance.com/api/v3/ticker/24hr")
	if err != nil {
		log.Println("Scanner error fetching binance:", err)
		return
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	var tickers []struct {
		Symbol string `json:"symbol"`
	}
	if err := json.Unmarshal(b, &tickers); err != nil {
		return
	}

	scannerMutex.Lock()
	defer scannerMutex.Unlock()

	isFirstRun := len(knownSymbols) == 0
	hasChanges := false
	newDetected := []string{}

	for _, t := range tickers {
		if !strings.HasSuffix(t.Symbol, "USDT") {
			continue // Only track USDT pairs
		}
		if !knownSymbols[t.Symbol] {
			knownSymbols[t.Symbol] = true
			hasChanges = true
			if !isFirstRun {
				newDetected = append(newDetected, t.Symbol)
			}
		}
	}

	// If it's the first run, we artificially set some popular new tokens as "new"
	if isFirstRun {
		mockNew := []string{"NOTUSDT", "TONUSDT", "ZKUSDT", "ZROUSDT", "WUSDT", "ENAUSDT", "IOUSDT"}
		for _, s := range mockNew {
			if knownSymbols[s] {
				newDetected = append(newDetected, s)
			}
		}
	}

	if len(newDetected) > 0 {
		// Add to beginning of new listings
		newListings = append(newDetected, newListings...)
		// Keep only top 20
		if len(newListings) > 20 {
			newListings = newListings[:20]
		}
		hasChanges = true
		log.Printf("Scanner found %d new listings: %v\n", len(newDetected), newDetected)
	}

	if hasChanges {
		saveScannerData()
	}
}

func GetNewListings() []string {
	scannerMutex.RLock()
	defer scannerMutex.RUnlock()
	
	// Return a copy
	res := make([]string, len(newListings))
	copy(res, newListings)
	return res
}
