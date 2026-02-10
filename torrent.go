package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// DownloadClient represents a download manager client
type DownloadClient struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Command string `json:"command"`
	Args    string `json:"args"` // %s will be replaced with URL/magnet
}

// Supported download clients (priority order)
var downloadClients = []DownloadClient{
	{ID: "fdm", Name: "Free Download Manager", Command: "fdm", Args: "%s"},
	{ID: "qbittorrent", Name: "qBittorrent", Command: "qbittorrent", Args: "%s"},
	{ID: "transmission", Name: "Transmission", Command: "transmission-gtk", Args: "%s"},
	{ID: "deluge", Name: "Deluge", Command: "deluge", Args: "%s"},
	{ID: "ktorrent", Name: "KTorrent", Command: "ktorrent", Args: "%s"},
}

// Download represents a download entry
type Download struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	MagnetURI  string  `json:"magnetUri"`
	TorrentURL string  `json:"torrentUrl,omitempty"`
	Progress   float64 `json:"progress"`
	Status     string  `json:"status"` // "added", "sent_to_fdm", "error"
	Error      string  `json:"error,omitempty"`
	SavePath   string  `json:"savePath"`
	AddedAt    int64   `json:"addedAt"`
	Client     string  `json:"client"` // Which client was used
}

// DownloadManager manages downloads history
type DownloadManager struct {
	downloads       map[string]*Download
	mu              sync.RWMutex
	savePath        string
	dataFile        string
	preferredClient string
}

var downloadManager *DownloadManager

// InitDownloadManager initializes the download manager
func (a *App) InitDownloadManager() error {
	if downloadManager != nil {
		return nil
	}

	// Get default paths
	homeDir, _ := os.UserHomeDir()

	// Data file for persistence (only config, no download folder needed - FDM handles downloads)
	configDir := filepath.Join(homeDir, ".config", "Knit")
	os.MkdirAll(configDir, 0755)
	dataFile := filepath.Join(configDir, "downloads.json")

	downloadManager = &DownloadManager{
		downloads:       make(map[string]*Download),
		savePath:        "", // Not used - FDM handles download paths
		dataFile:        dataFile,
		preferredClient: "fdm", // Default to FDM
	}

	// Load existing downloads
	downloadManager.loadFromFile()

	return nil
}

// loadFromFile loads downloads history from file
func (dm *DownloadManager) loadFromFile() {
	data, err := os.ReadFile(dm.dataFile)
	if err != nil {
		return
	}

	var downloads []*Download
	if err := json.Unmarshal(data, &downloads); err != nil {
		return
	}

	for _, d := range downloads {
		dm.downloads[d.ID] = d
	}
}

// saveToFile saves downloads history to file
func (dm *DownloadManager) saveToFile() {
	dm.mu.RLock()
	downloads := make([]*Download, 0, len(dm.downloads))
	for _, d := range dm.downloads {
		downloads = append(downloads, d)
	}
	dm.mu.RUnlock()

	data, err := json.MarshalIndent(downloads, "", "  ")
	if err != nil {
		return
	}

	os.WriteFile(dm.dataFile, data, 0644)
}

// GetAvailableClients returns list of available download clients
func (a *App) GetAvailableClients() []DownloadClient {
	available := []DownloadClient{}
	for _, client := range downloadClients {
		if _, err := exec.LookPath(client.Command); err == nil {
			available = append(available, client)
		}
	}
	return available
}

// findAvailableClient finds the first available download client
func findAvailableClient() (*DownloadClient, error) {
	// First try preferred client (FDM)
	for _, client := range downloadClients {
		if _, err := exec.LookPath(client.Command); err == nil {
			return &client, nil
		}
	}
	return nil, fmt.Errorf("no download manager found. Please install Free Download Manager, qBittorrent, or another torrent client")
}

// sendToFDM sends a magnet/torrent to Free Download Manager
func sendToFDM(uri string) error {
	// FDM accepts URLs directly as argument
	cmd := exec.Command("fdm", uri)
	return cmd.Start()
}

// sendToClient sends a magnet/torrent to a download client
func sendToClient(client *DownloadClient, uri string) error {
	args := strings.Replace(client.Args, "%s", uri, 1)
	cmd := exec.Command(client.Command, args)
	return cmd.Start()
}

// AddMagnet adds a magnet link to the download manager (FDM or other)
func (a *App) AddMagnet(magnetURI string, name string) (*Download, error) {
	if downloadManager == nil {
		if err := a.InitDownloadManager(); err != nil {
			return nil, err
		}
	}

	// Find available client
	client, err := findAvailableClient()
	if err != nil {
		return nil, err
	}

	// Parse magnet to get hash
	hash := extractHashFromMagnet(magnetURI)
	if hash == "" {
		hash = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	downloadManager.mu.Lock()

	// Check if already exists
	if existing, ok := downloadManager.downloads[hash]; ok {
		downloadManager.mu.Unlock()
		// Re-send to client
		go sendToClient(client, magnetURI)
		return existing, nil
	}

	download := &Download{
		ID:        hash,
		Name:      name,
		MagnetURI: magnetURI,
		Status:    "sending",
		SavePath:  downloadManager.savePath,
		AddedAt:   time.Now().Unix(),
		Client:    client.Name,
	}

	downloadManager.downloads[hash] = download
	downloadManager.mu.Unlock()

	// Send to download client
	go func() {
		err := sendToClient(client, magnetURI)
		downloadManager.mu.Lock()
		if err != nil {
			download.Status = "error"
			download.Error = fmt.Sprintf("Failed to send to %s: %v", client.Name, err)
		} else {
			download.Status = "sent_to_" + client.ID
		}
		downloadManager.mu.Unlock()
		downloadManager.saveToFile()
	}()

	downloadManager.saveToFile()
	return download, nil
}

// AddTorrentURL downloads a .torrent file and sends it to FDM
func (a *App) AddTorrentURL(torrentURL string, name string) (*Download, error) {
	if downloadManager == nil {
		if err := a.InitDownloadManager(); err != nil {
			return nil, err
		}
	}

	// Find available client
	client, err := findAvailableClient()
	if err != nil {
		return nil, err
	}

	hash := fmt.Sprintf("%d", time.Now().UnixNano())

	downloadManager.mu.Lock()

	download := &Download{
		ID:         hash,
		Name:       name,
		TorrentURL: torrentURL,
		Status:     "downloading_torrent",
		SavePath:   downloadManager.savePath,
		AddedAt:    time.Now().Unix(),
		Client:     client.Name,
	}

	downloadManager.downloads[hash] = download
	downloadManager.mu.Unlock()

	// Download torrent file and send to client
	go func() {
		torrentPath, err := downloadTorrentFile(torrentURL, name, downloadManager.savePath)
		if err != nil {
			downloadManager.mu.Lock()
			download.Status = "error"
			download.Error = err.Error()
			downloadManager.mu.Unlock()
			downloadManager.saveToFile()
			return
		}

		// Send to download client
		err = sendToClient(client, torrentPath)
		downloadManager.mu.Lock()
		if err != nil {
			download.Status = "error"
			download.Error = fmt.Sprintf("Failed to send to %s: %v", client.Name, err)
		} else {
			download.Status = "sent_to_" + client.ID
		}
		downloadManager.mu.Unlock()
		downloadManager.saveToFile()
	}()

	downloadManager.saveToFile()
	return download, nil
}

// downloadTorrentFile downloads a .torrent file and returns its path
func downloadTorrentFile(torrentURL string, name string, savePath string) (string, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(torrentURL)
	if err != nil {
		return "", fmt.Errorf("failed to download torrent: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download torrent: HTTP %d", resp.StatusCode)
	}

	// Create safe filename
	safeName := strings.Map(func(r rune) rune {
		if r == '/' || r == '\\' || r == ':' || r == '*' || r == '?' || r == '"' || r == '<' || r == '>' || r == '|' {
			return '_'
		}
		return r
	}, name)
	if !strings.HasSuffix(strings.ToLower(safeName), ".torrent") {
		safeName += ".torrent"
	}

	torrentPath := filepath.Join(savePath, safeName)

	// Save torrent file
	file, err := os.Create(torrentPath)
	if err != nil {
		return "", fmt.Errorf("failed to create torrent file: %v", err)
	}

	_, err = io.Copy(file, resp.Body)
	file.Close()
	if err != nil {
		os.Remove(torrentPath)
		return "", fmt.Errorf("failed to save torrent file: %v", err)
	}

	return torrentPath, nil
}

// GetDownloads returns all downloads
func (a *App) GetDownloads() []*Download {
	if downloadManager == nil {
		return []*Download{}
	}

	downloadManager.mu.RLock()
	defer downloadManager.mu.RUnlock()

	downloads := make([]*Download, 0, len(downloadManager.downloads))
	for _, d := range downloadManager.downloads {
		downloads = append(downloads, d)
	}

	return downloads
}

// RemoveDownload removes a download from history
func (a *App) RemoveDownload(id string) error {
	if downloadManager == nil {
		return fmt.Errorf("download manager not initialized")
	}

	downloadManager.mu.Lock()
	delete(downloadManager.downloads, id)
	downloadManager.mu.Unlock()

	downloadManager.saveToFile()
	return nil
}

// ClearDownloads clears all downloads history
func (a *App) ClearDownloads() error {
	if downloadManager == nil {
		return nil
	}

	downloadManager.mu.Lock()
	downloadManager.downloads = make(map[string]*Download)
	downloadManager.mu.Unlock()

	downloadManager.saveToFile()
	return nil
}

// OpenDownloadFolder opens the download folder
func (a *App) OpenDownloadFolder() error {
	if downloadManager == nil {
		homeDir, _ := os.UserHomeDir()
		return openFolder(filepath.Join(homeDir, "Downloads", "Knit"))
	}
	return openFolder(downloadManager.savePath)
}

// OpenFDM opens Free Download Manager application
func (a *App) OpenFDM() error {
	client, err := findAvailableClient()
	if err != nil {
		return err
	}
	cmd := exec.Command(client.Command)
	return cmd.Start()
}

// Helper functions

func extractHashFromMagnet(magnetURI string) string {
	// Extract hash from magnet:?xt=urn:btih:HASH...
	lower := strings.ToLower(magnetURI)
	if idx := strings.Index(lower, "btih:"); idx != -1 {
		start := idx + 5
		end := start
		for end < len(magnetURI) && magnetURI[end] != '&' {
			end++
		}
		return magnetURI[start:end]
	}
	return ""
}

func openFolder(path string) error {
	cmd := exec.Command("xdg-open", path)
	return cmd.Start()
}

// CopyMagnetToClipboard copies magnet link to clipboard
func (a *App) CopyMagnetToClipboard(magnetURI string) error {
	// Try xclip
	cmd := exec.Command("xclip", "-selection", "clipboard")
	cmd.Stdin = strings.NewReader(magnetURI)
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Try xsel
	cmd = exec.Command("xsel", "--clipboard", "--input")
	cmd.Stdin = strings.NewReader(magnetURI)
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Try wl-copy (Wayland)
	cmd = exec.Command("wl-copy")
	cmd.Stdin = strings.NewReader(magnetURI)
	if err := cmd.Run(); err == nil {
		return nil
	}

	return fmt.Errorf("no clipboard tool found (install xclip, xsel, or wl-clipboard)")
}

// GetMagnetFromTorrentURL extracts magnet URI from .torrent file URL (if possible)
func (a *App) GetMagnetFromTorrentURL(torrentURL string) string {
	// Some trackers provide magnet links directly - check if URL is already a magnet
	if strings.HasPrefix(torrentURL, "magnet:") {
		return torrentURL
	}

	// Try to convert torrent URL to magnet (basic approach)
	parsed, err := url.Parse(torrentURL)
	if err != nil {
		return ""
	}

	// If URL has hash parameter, try to construct magnet
	hash := parsed.Query().Get("hash")
	if hash == "" {
		hash = parsed.Query().Get("btih")
	}
	if hash != "" {
		return fmt.Sprintf("magnet:?xt=urn:btih:%s", hash)
	}

	return ""
}
