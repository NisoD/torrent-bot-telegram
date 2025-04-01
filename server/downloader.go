package server

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/cenkalti/rain/torrent"
)

// DownloadProgress - state of a download
type DownloadProgress struct {
	Status          string
	PercentComplete float64
	BytesCompleted  int64
	BytesTotal      int64
	Peers           int
}

// File with it's path
type TorrentFile struct {
	ID       int
	Name     string
	Path     string
	Selected bool
}

// Handler
type Downloader struct {
	session      *torrent.Session
	torrent      *torrent.Torrent
	downloadPath string
	files        []TorrentFile
	logger       *Logger
	mu           sync.Mutex
}

func NewDownloader(downloadPath string, logger *Logger) *Downloader {
	if logger == nil {
		// Create a default logger that outputs to stdout if none provided
		log, _ := NewLogger(filepath.Join(downloadPath, "logs"), true)
		logger = log
	}

	return &Downloader{
		downloadPath: downloadPath,
		logger:       logger,
	}
}

// GetTorrentInfo retrieves information about the torrent without starting the download
func (d *Downloader) GetTorrentInfo(magnetLink string) ([]TorrentFile, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// session config
	cfg := torrent.DefaultConfig
	cfg.DataDir = d.downloadPath

	// new session
	ses, err := torrent.NewSession(cfg)
	if err != nil {
		d.logger.LogError("Failed to create torrent session: %v", err)
		return nil, err
	}
	d.session = ses

	// Add magnet link
	tor, err := ses.AddURI(magnetLink, nil)
	if err != nil {
		d.logger.LogError("Failed to add magnet URI: %v", err)
		return nil, err
	}
	d.torrent = tor

	// Wait for metadata and inform user
	d.logger.LogInfo("Fetching torrent metadata...")
	metadataComplete := make(chan struct{})
	go func() {
		for range time.Tick(time.Second) {
			s := tor.Stats()
			if s.Status == torrent.Downloading ||
				s.Status == torrent.Stopped ||
				s.Status == torrent.Seeding {
				close(metadataComplete)
				return
			}
		}
	}()

	
	select {
	case <-metadataComplete:
		// Metadata fetched successfully
	case <-time.After(60 * time.Second): // This can be prolonged for rare low seed usage
		return nil, errors.New("timeout while fetching torrent metadata")
	}

	// Populate files list
	torrentFiles, err := tor.Files()
	if err != nil {
		return nil, err
	}

	d.files = make([]TorrentFile, len(torrentFiles))

	for i, file := range torrentFiles {
		d.files[i] = TorrentFile{
			ID:       i,
			Name:     filepath.Base(file.Path()),
			Path:     file.Path(),
			Selected: true,
		}
	}
	
	tor.Stop()
	return d.files, nil
}

// Important Feature - Lets user select files to download (Not availabe for IOS users in ISH + rtorrent usage) 
func (d *Downloader) SelectFiles(fileIDs []int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.torrent == nil {
		return errors.New("no active torrent")
	}

	// map for quick lookup 
	selectedMap := make(map[int]bool)
	for _, id := range fileIDs {
		selectedMap[id] = true
	}

	// Update selection status
	for i := range d.files {
		d.files[i].Selected = selectedMap[i]
	}

	/**File selection is done by the torrent client, and not by setting priorities 
	The torrent client will download the files selected by default
	if not selected they are skiped 
	The d.files array contain the selected files
	**/
	return nil
}

// Download starts downloading selected files from a torrent
// Returns a channel that will receive progress updates
func (d *Downloader) Download() (chan DownloadProgress, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// torrent is set
	if d.torrent == nil {
		return nil, errors.New("no active torrent")
	}

	// To update user of the prog we make a channel
	progressChan := make(chan DownloadProgress)

	d.torrent.Start()
	d.logger.LogInfo("Starting download of selected files")

	// goroutine for monitoring progress
	go func() {
		defer close(progressChan)

		for range time.Tick(time.Second) {
			// Check if torrent is still valid
			if d.torrent == nil {
				d.logger.LogError("Torrent is no longer valid")
				return
			}

			s := d.torrent.Stats()
			// Precentage for readability
			var percentComplete float64
			if s.Bytes.Total > 0 {
				percentComplete = float64(s.Bytes.Completed) / float64(s.Bytes.Total) * 100
			}

			// user is updated
			progressInfo := DownloadProgress{
				Status:          s.Status.String(),
				PercentComplete: percentComplete,
				BytesCompleted:  s.Bytes.Completed,
				BytesTotal:      s.Bytes.Total,
				Peers:           s.Peers.Total,
			}

			progressChan <- progressInfo

			// Proggress update can be more/less 
			if int(percentComplete)%10 == 0 {
				d.logger.LogInfo("Download progress: %.2f%% (Status: %s, Peers: %d)",
					percentComplete, s.Status.String(), s.Peers.Total)
			}

			// if seeding download is done
			if s.Status == torrent.Seeding {
				d.logger.LogInfo("Download complete")
				return
			}
		}
	}()

	return progressChan, nil
}

// GetDownloadedFiles returns the list of downloaded files
func (d *Downloader) GetDownloadedFiles() ([]TorrentFile, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.torrent == nil {
		return nil, errors.New("no active torrent")
	}

	// Update file status
	torrentFiles, err := d.torrent.Files()
	if err != nil {
		return nil, errors.New("problem loading torrent files")
	}

	for i, file := range torrentFiles {
		if i < len(d.files) && d.files[i].Selected {
			// Check if file exists and update its path
			fullPath := filepath.Join(d.downloadPath, file.Path())
			if _, err := os.Stat(fullPath); err == nil {
				d.files[i].Path = fullPath
			} else {
				foundPath, err := findFileInSubdirectories(d.downloadPath, file.Path())
				if err == nil {
					d.files[i].Path = foundPath
				}
			}
		}
	}

	// Return only selected files
	var downloadedFiles []TorrentFile
	for _, file := range d.files {
		if file.Selected {
			downloadedFiles = append(downloadedFiles, file)
		}
	}

	return downloadedFiles, nil
}

func findFileInSubdirectories(root, relativePath string) (string, error) {
	var foundPath string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Base(path) == filepath.Base(relativePath) {
			if foundPath == "" {
				foundPath = path
				return errors.New("found")
			}
		}
		return nil
	})
	if err != nil && err.Error() == "found" {
		return foundPath, nil
	}
	if foundPath != "" {
		return foundPath, nil
	}
	return "", errors.New("file not found")
}

// Close cleans up resources
func (d *Downloader) Close() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.session != nil {
		d.session.Close()
		d.session = nil
		d.torrent = nil
	}

	d.logger.LogInfo("Downloader closed")
}

// FormatBytes returns a human-readable byte string
func FormatBytes(bytes int64) string {
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
