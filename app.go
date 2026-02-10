package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// TMDBProxy represents a TMDB proxy server
type TMDBProxy struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	ImageURL string `json:"imageUrl"`
	Type     string `json:"type"` // "direct" or "cors"
	IsOnline bool   `json:"isOnline"`
}

// TMDB API Key
const tmdbAPIKey = "4ef0d7355d9ffb5151e987764708ce96"

// Predefined TMDB proxy servers
var tmdbProxies = []TMDBProxy{
	{ID: "apn_render", Name: "APN Render", URL: "https://apn-latest.onrender.com/https://api.themoviedb.org", ImageURL: "https://imagetmdb.com", Type: "cors"},
	{ID: "tmdb_direct", Name: "TMDB Direct", URL: "https://api.themoviedb.org", ImageURL: "https://image.tmdb.org", Type: "direct"},
}

// HTTP clients with optimized settings
var (
	httpClient     = &http.Client{Timeout: 15 * time.Second}
	httpClientFast = &http.Client{Timeout: 5 * time.Second}
)

// App struct
type App struct{}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// OnStartup is called when the app starts (Wails v3 lifecycle)
func (a *App) OnStartup() {
	loadSettingsFromFile()
}

// PublicServer represents a public Jackett server
type PublicServer struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Key      string `json:"key"`
	Type     string `json:"type"`
	IsOnline bool   `json:"isOnline"`
}

// TorrentResult represents a torrent search result
type TorrentResult struct {
	Title       string `json:"title"`
	Size        int64  `json:"size"`
	SizeStr     string `json:"sizeStr"`
	Seeders     int    `json:"seeders"`
	Leechers    int    `json:"leechers"`
	MagnetURI   string `json:"magnetUri"`
	Link        string `json:"link"`
	Tracker     string `json:"tracker"`
	PublishDate string `json:"publishDate"`
	Category    string `json:"category"`
}

// TMDBMovie represents a movie from TMDB
type TMDBMovie struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	OriginalTitle string  `json:"original_title"`
	Overview      string  `json:"overview"`
	PosterPath    string  `json:"poster_path"`
	BackdropPath  string  `json:"backdrop_path"`
	ReleaseDate   string  `json:"release_date"`
	VoteAverage   float64 `json:"vote_average"`
	VoteCount     int     `json:"vote_count"`
	Popularity    float64 `json:"popularity"`
	GenreIDs      []int   `json:"genre_ids"`
}

// TMDBTVShow represents a TV show from TMDB
type TMDBTVShow struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	OriginalName  string   `json:"original_name"`
	Overview      string   `json:"overview"`
	PosterPath    string   `json:"poster_path"`
	BackdropPath  string   `json:"backdrop_path"`
	FirstAirDate  string   `json:"first_air_date"`
	VoteAverage   float64  `json:"vote_average"`
	VoteCount     int      `json:"vote_count"`
	Popularity    float64  `json:"popularity"`
	GenreIDs      []int    `json:"genre_ids"`
	OriginCountry []string `json:"origin_country"`
}

// TMDBTVSearchResponse represents TMDB TV search response
type TMDBTVSearchResponse struct {
	Page         int          `json:"page"`
	Results      []TMDBTVShow `json:"results"`
	TotalPages   int          `json:"total_pages"`
	TotalResults int          `json:"total_results"`
}

// TMDBGenre represents a genre from TMDB
type TMDBGenre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TMDBGenresResponse represents TMDB genres response
type TMDBGenresResponse struct {
	Genres []TMDBGenre `json:"genres"`
}

// TMDBSearchResponse represents TMDB search response
type TMDBSearchResponse struct {
	Page         int         `json:"page"`
	Results      []TMDBMovie `json:"results"`
	TotalPages   int         `json:"total_pages"`
	TotalResults int         `json:"total_results"`
}

// JackettResponse represents Jackett API response
type JackettResponse struct {
	Results []struct {
		Title                string  `json:"Title"`
		Size                 int64   `json:"Size"`
		Seeders              int     `json:"Seeders"`
		Peers                int     `json:"Peers"`
		MagnetUri            string  `json:"MagnetUri"`
		Link                 string  `json:"Link"`
		Tracker              string  `json:"Tracker"`
		PublishDate          string  `json:"PublishDate"`
		CategoryDesc         string  `json:"CategoryDesc"`
		TrackerId            string  `json:"TrackerId"`
		TrackerType          string  `json:"TrackerType"`
		Guid                 string  `json:"Guid"`
		Details              string  `json:"Details"`
		DownloadVolumeFactor float64 `json:"DownloadVolumeFactor"`
		UploadVolumeFactor   float64 `json:"UploadVolumeFactor"`
	} `json:"Results"`
	Indexers []struct {
		ID      string `json:"ID"`
		Name    string `json:"Name"`
		Status  int    `json:"Status"`
		Results int    `json:"Results"`
		Error   string `json:"Error"`
	} `json:"Indexers"`
}

// Predefined public servers
var publicServers = []PublicServer{
	{ID: "jacred_xyz", Name: "Jacred.xyz", URL: "jacred.xyz", Key: "", Type: "jackett"},
	{ID: "jac_red_ru", Name: "Jac-red.ru", URL: "jac-red.ru", Key: "", Type: "jackett"},
	{ID: "jac_red", Name: "Jac.red", URL: "jac.red", Key: "", Type: "jackett"},
}

// GetPublicServers returns list of available public servers with health check
func (a *App) GetPublicServers() []PublicServer {
	var wg sync.WaitGroup
	servers := make([]PublicServer, len(publicServers))
	copy(servers, publicServers)

	for i := range servers {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			servers[idx].IsOnline = a.checkServerHealth(servers[idx])
		}(i)
	}
	wg.Wait()
	return servers
}

// GetTMDBProxies returns list of TMDB proxy servers with health check
func (a *App) GetTMDBProxies() []TMDBProxy {
	var wg sync.WaitGroup
	proxies := make([]TMDBProxy, len(tmdbProxies))
	copy(proxies, tmdbProxies)

	for i := range proxies {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			proxies[idx].IsOnline = a.checkTMDBProxyHealth(proxies[idx])
		}(i)
	}
	wg.Wait()
	return proxies
}

// checkTMDBProxyHealth checks if a TMDB proxy is online
func (a *App) checkTMDBProxyHealth(proxy TMDBProxy) bool {
	urlStr := fmt.Sprintf("%s/3/movie/popular?api_key=%s&language=ru-RU&page=1", proxy.URL, tmdbAPIKey)
	resp, err := httpClientFast.Get(urlStr)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

// getSelectedTMDBProxy returns the currently selected TMDB proxy
func (a *App) getSelectedTMDBProxy() *TMDBProxy {
	for _, p := range tmdbProxies {
		if p.ID == appSettings.TMDBProxy {
			return &p
		}
	}
	// Default to first proxy
	if len(tmdbProxies) > 0 {
		return &tmdbProxies[0]
	}
	return nil
}

// GetImageURL returns the proxied image URL for TMDB images
func (a *App) GetImageURL(path string, size string) string {
	proxy := a.getSelectedTMDBProxy()
	if proxy == nil || path == "" {
		return ""
	}

	// Handle different proxy types
	imageBase := proxy.ImageURL
	if strings.Contains(imageBase, "weserv.nl") {
		// weserv.nl format: https://images.weserv.nl/?url=image.tmdb.org/t/p/w500/path
		return fmt.Sprintf("%s/t/p/%s%s", imageBase, size, path)
	}

	// Standard format
	return fmt.Sprintf("%s/t/p/%s%s", imageBase, size, path)
}

// GetCurrentImageBase returns the current image base URL
func (a *App) GetCurrentImageBase() string {
	proxy := a.getSelectedTMDBProxy()
	if proxy == nil {
		return "https://image.tmdb.org"
	}
	return proxy.ImageURL
}

// CheckServerByID checks if a specific server is online by ID
func (a *App) CheckServerByID(serverID string) bool {
	for _, s := range publicServers {
		if s.ID == serverID {
			return a.checkServerHealth(s)
		}
	}
	return false
}

// CheckTMDBProxyByID checks if a specific TMDB proxy is online by ID
func (a *App) CheckTMDBProxyByID(proxyID string) bool {
	for _, p := range tmdbProxies {
		if p.ID == proxyID {
			return a.checkTMDBProxyHealth(p)
		}
	}
	return false
}

// checkServerHealth checks if a server is online
func (a *App) checkServerHealth(server PublicServer) bool {
	client := &http.Client{Timeout: 5 * time.Second}

	protocol := "https://"
	apiPath := "/api/v2.0/indexers/status:healthy/results"
	urlStr := fmt.Sprintf("%s%s%s?apikey=%s", protocol, server.URL, apiPath, server.Key)

	resp, err := client.Get(urlStr)
	if err != nil {
		// Try http if https fails
		protocol = "http://"
		urlStr = fmt.Sprintf("%s%s%s?apikey=%s", protocol, server.URL, apiPath, server.Key)
		resp, err = client.Get(urlStr)
		if err != nil {
			return false
		}
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200 || resp.StatusCode == 401
}

// SearchTorrents searches for torrents on the specified server
func (a *App) SearchTorrents(serverID string, query string) ([]TorrentResult, error) {
	var server *PublicServer
	for _, s := range publicServers {
		if s.ID == serverID {
			server = &s
			break
		}
	}

	if server == nil {
		return nil, fmt.Errorf("server not found: %s", serverID)
	}

	client := &http.Client{Timeout: 30 * time.Second}

	encodedQuery := url.QueryEscape(query)
	apiPath := "/api/v2.0/indexers/status:healthy/results"

	var resp *http.Response
	var err error
	var urlStr string

	// Retry logic for rate limiting (429)
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt*2) * time.Second) // Wait 2s, 4s between retries
		}

		// Try https first
		urlStr = fmt.Sprintf("https://%s%s?apikey=%s&Query=%s", server.URL, apiPath, server.Key, encodedQuery)

		req, _ := http.NewRequest("GET", urlStr, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

		resp, err = client.Do(req)
		if err != nil {
			// Try http
			urlStr = fmt.Sprintf("http://%s%s?apikey=%s&Query=%s", server.URL, apiPath, server.Key, encodedQuery)
			req, _ = http.NewRequest("GET", urlStr, nil)
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
			resp, err = client.Do(req)
			if err != nil {
				continue // Retry
			}
		}

		// If rate limited, retry
		if resp.StatusCode == 429 {
			resp.Body.Close()
			continue
		}

		break
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("no response from server after %d retries", maxRetries)
	}
	defer resp.Body.Close()

	// Check for rate limiting after all retries
	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("сервер перегружен (429). Попробуйте другой сервер или подождите минуту")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Check if response is HTML (error page)
	bodyStr := string(body)
	if len(bodyStr) > 0 && bodyStr[0] == '<' {
		return nil, fmt.Errorf("сервер вернул ошибку (status: %d). Попробуйте другой сервер", resp.StatusCode)
	}

	var jackettResp JackettResponse
	if err := json.Unmarshal(body, &jackettResp); err != nil {
		// Try to provide more context
		preview := bodyStr
		if len(preview) > 100 {
			preview = preview[:100]
		}
		return nil, fmt.Errorf("failed to parse response: %v (preview: %s)", err, preview)
	}

	results := make([]TorrentResult, 0, len(jackettResp.Results))
	for _, r := range jackettResp.Results {
		results = append(results, TorrentResult{
			Title:       r.Title,
			Size:        r.Size,
			SizeStr:     formatBytes(r.Size),
			Seeders:     r.Seeders,
			Leechers:    r.Peers - r.Seeders,
			MagnetURI:   r.MagnetUri,
			Link:        r.Link,
			Tracker:     r.Tracker,
			PublishDate: r.PublishDate,
			Category:    r.CategoryDesc,
		})
	}

	// Sort by seeders
	sort.Slice(results, func(i, j int) bool {
		return results[i].Seeders > results[j].Seeders
	})

	return results, nil
}

// tmdbRequest makes a request to TMDB API
func (a *App) tmdbRequest(path string) ([]byte, error) {
	proxy := a.getSelectedTMDBProxy()
	if proxy == nil {
		return nil, fmt.Errorf("no TMDB proxy available")
	}

	urlStr := fmt.Sprintf("%s/3%s", proxy.URL, path)
	if !strings.Contains(path, "api_key") {
		if strings.Contains(path, "?") {
			urlStr += "&api_key=" + tmdbAPIKey
		} else {
			urlStr += "?api_key=" + tmdbAPIKey
		}
	}

	resp, err := httpClient.Get(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to TMDB: %v", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// SearchTMDB searches for movies on TMDB via proxy
func (a *App) SearchTMDB(query string) (*TMDBSearchResponse, error) {
	path := fmt.Sprintf("/search/movie?api_key=%s&query=%s&language=ru-RU", tmdbAPIKey, url.QueryEscape(query))
	body, err := a.tmdbRequest(path)
	if err != nil {
		return nil, err
	}

	var resp TMDBSearchResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	return &resp, nil
}

// GetTMDBMovieDetails gets detailed movie info from TMDB via proxy
func (a *App) GetTMDBMovieDetails(movieID int) (*TMDBMovie, error) {
	path := fmt.Sprintf("/movie/%d?api_key=%s&language=ru-RU", movieID, tmdbAPIKey)
	body, err := a.tmdbRequest(path)
	if err != nil {
		return nil, err
	}

	var movie TMDBMovie
	if err := json.Unmarshal(body, &movie); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	return &movie, nil
}

// GetGenres gets list of movie genres from TMDB
func (a *App) GetGenres() (*TMDBGenresResponse, error) {
	path := fmt.Sprintf("/genre/movie/list?api_key=%s&language=ru-RU", tmdbAPIKey)
	body, err := a.tmdbRequest(path)
	if err != nil {
		return nil, err
	}

	var resp TMDBGenresResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	return &resp, nil
}

// GetPopularMovies gets popular movies from TMDB via proxy
func (a *App) GetPopularMovies(page int) (*TMDBSearchResponse, error) {
	if page < 1 {
		page = 1
	}
	path := fmt.Sprintf("/movie/popular?api_key=%s&language=ru-RU&page=%d", tmdbAPIKey, page)
	body, err := a.tmdbRequest(path)
	if err != nil {
		return nil, err
	}

	var resp TMDBSearchResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	return &resp, nil
}

// SearchTMDBTV searches for TV shows on TMDB via proxy
func (a *App) SearchTMDBTV(query string) (*TMDBTVSearchResponse, error) {
	path := fmt.Sprintf("/search/tv?api_key=%s&query=%s&language=ru-RU", tmdbAPIKey, url.QueryEscape(query))
	body, err := a.tmdbRequest(path)
	if err != nil {
		return nil, err
	}

	var resp TMDBTVSearchResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	return &resp, nil
}

// GetPopularTVShows gets popular TV shows from TMDB via proxy
func (a *App) GetPopularTVShows(page int) (*TMDBTVSearchResponse, error) {
	if page < 1 {
		page = 1
	}
	path := fmt.Sprintf("/tv/popular?api_key=%s&language=ru-RU&page=%d", tmdbAPIKey, page)
	body, err := a.tmdbRequest(path)
	if err != nil {
		return nil, err
	}

	var resp TMDBTVSearchResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	return &resp, nil
}

// formatBytes converts bytes to human readable string
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Settings struct for app configuration
type Settings struct {
	SelectedServer string `json:"selectedServer"`
	TMDBProxy      string `json:"tmdbProxy"`
	Theme          string `json:"theme"`
}

var appSettings = Settings{
	SelectedServer: "jacred_xyz",
	TMDBProxy:      "apn_render",
	Theme:          "dark",
}

// getSettingsPath returns the path to settings file
func getSettingsPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	return filepath.Join(configDir, "Knit", "settings.json")
}

// loadSettingsFromFile loads settings from disk
func loadSettingsFromFile() {
	path := getSettingsPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return // Use defaults
	}
	json.Unmarshal(data, &appSettings)
}

// saveSettingsToFile saves settings to disk
func saveSettingsToFile() error {
	path := getSettingsPath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(appSettings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// GetSettings returns current app settings
func (a *App) GetSettings() Settings {
	return appSettings
}

// SaveSettings saves app settings
func (a *App) SaveSettings(settings Settings) error {
	appSettings = settings
	return saveSettingsToFile()
}

// SearchWithMovie searches for torrents with movie title
func (a *App) SearchWithMovie(movie TMDBMovie, serverID string) ([]TorrentResult, error) {
	var results []TorrentResult
	var err error

	// Try original (English) title first - usually better results
	if movie.OriginalTitle != "" {
		results, err = a.SearchTorrents(serverID, movie.OriginalTitle)
		if err == nil && len(results) > 0 {
			return results, nil
		}
	}

	// Try Russian title
	if movie.Title != "" && movie.Title != movie.OriginalTitle {
		results, err = a.SearchTorrents(serverID, movie.Title)
		if err == nil && len(results) > 0 {
			return results, nil
		}
	}

	// Try original title with year
	if movie.OriginalTitle != "" && movie.ReleaseDate != "" && len(movie.ReleaseDate) >= 4 {
		year := movie.ReleaseDate[:4]
		query := fmt.Sprintf("%s %s", movie.OriginalTitle, year)
		results, err = a.SearchTorrents(serverID, query)
		if err == nil && len(results) > 0 {
			return results, nil
		}
	}

	// Return last error or empty results
	if err != nil {
		return nil, err
	}
	return results, nil
}

// SearchWithTVShow searches for torrents with TV show title
func (a *App) SearchWithTVShow(show TMDBTVShow, serverID string) ([]TorrentResult, error) {
	var results []TorrentResult
	var err error

	// Try original (English) name first - usually better results
	if show.OriginalName != "" {
		results, err = a.SearchTorrents(serverID, show.OriginalName)
		if err == nil && len(results) > 0 {
			return results, nil
		}
	}

	// Try Russian name
	if show.Name != "" && show.Name != show.OriginalName {
		results, err = a.SearchTorrents(serverID, show.Name)
		if err == nil && len(results) > 0 {
			return results, nil
		}
	}

	// Try original name with year
	if show.OriginalName != "" && show.FirstAirDate != "" && len(show.FirstAirDate) >= 4 {
		year := show.FirstAirDate[:4]
		query := fmt.Sprintf("%s %s", show.OriginalName, year)
		results, err = a.SearchTorrents(serverID, query)
		if err == nil && len(results) > 0 {
			return results, nil
		}
	}

	// Return last error or empty results
	if err != nil {
		return nil, err
	}
	return results, nil
}
