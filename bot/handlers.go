package bot

import (
	"BotTelegram/server"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// UserState represents the current state of a user in the conversation
type UserState int

const (
	StateNone UserState = iota
	StateAwaitingMagnet
	StateSelectingFiles
	StateDownloading
)

// UserSession represents a user's session
type UserSession struct {
	State       UserState
	Downloader  *server.Downloader
	MagnetLink  string
	Files       []server.TorrentFile
	ProgressMsg *tgbotapi.Message
}

// Global sessions map
var (
	sessions = make(map[int64]*UserSession)
	mu       sync.Mutex
)

// getSession gets or creates a user session
func (b *Bot) getSession(chatID int64) *UserSession {
	mu.Lock()
	defer mu.Unlock()

	session, exists := sessions[chatID]
	if !exists {
		session = &UserSession{
			State: StateNone,
			Downloader: server.NewDownloader(
				b.Config.AppConfig.DownloadPath,
				b.Logger,
			),
		}
		sessions[chatID] = session
	}
	return session
}

// handleMessage processes incoming messages and routes them to the appropriate handler
func (b *Bot) handleMessage(message *tgbotapi.Message) {
	// Get user session
	session := b.getSession(message.Chat.ID)

	// Check if it's a command
	if message.IsCommand() {
		b.handleCommand(message, session)
		return
	}

	// Handle message based on state
	switch session.State {
	case StateAwaitingMagnet:
		if strings.HasPrefix(message.Text, "magnet:") {
			b.handleMagnetLink(message, session)
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID,
				"That doesn't look like a magnet link. Please send a valid magnet link starting with 'magnet:'.")
			b.Config.API.Send(msg)
		}

	case StateSelectingFiles:
		b.handleFileSelection(message, session)

	default:
		// Default behavior - check if it's a magnet link
		if strings.HasPrefix(message.Text, "magnet:") {
			b.handleMagnetLink(message, session)
		} else {
			// Default response for unrecognized messages
			msg := tgbotapi.NewMessage(message.Chat.ID,
				"Send me a magnet link to download a torrent or use /help to see available commands.")
			b.Config.API.Send(msg)
		}
	}
}

// handleCommand processes bot commands
func (b *Bot) handleCommand(message *tgbotapi.Message, session *UserSession) {
	var reply string

	switch message.Command() {
	case "start":
		reply = "Welcome to Torrent Downloader Bot!\nSend me a magnet link, and I'll download it for you."
		session.State = StateAwaitingMagnet

	case "help":
		reply = "Available commands:\n" +
			"/start - Start the bot\n" +
			"/help - Show this help message\n" +
			"/cancel - Cancel current operation\n" +
			"\nOr simply send a magnet link to download a torrent."

	case "cancel":
		if session.Downloader != nil {
			session.Downloader.Close()
		}
		session.State = StateNone
		session.MagnetLink = ""
		session.Files = nil
		reply = "Current operation cancelled. You can send a new magnet link to start again."

	default:
		reply = "Unknown command. Use /help to see available commands."
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, reply)
	b.Config.API.Send(msg)
}

// handleMagnetLink processes magnet links and starts fetching metadata
func (b *Bot) handleMagnetLink(message *tgbotapi.Message, session *UserSession) {
	magnetLink := message.Text
	chatID := message.Chat.ID

	// Reset session
	if session.Downloader != nil {
		session.Downloader.Close()
	}
	session.MagnetLink = magnetLink
	session.Files = nil
	session.State = StateAwaitingMagnet

	// Send initial response
	msg := tgbotapi.NewMessage(chatID, "Fetching torrent metadata... This might take a moment.")
	sentMsg, err := b.Config.API.Send(msg)
	if err != nil {
		b.Logger.LogError("Error sending message: %v", err)
		return
	}

	// Fetch torrent info
	go func() {
		files, err := session.Downloader.GetTorrentInfo(magnetLink)
		if err != nil {
			updateMsg := tgbotapi.NewEditMessageText(chatID, sentMsg.MessageID,
				fmt.Sprintf("Error fetching torrent information: %v", err))
			b.Config.API.Send(updateMsg)
			return
		}

		// Store files in session
		session.Files = files
		session.State = StateSelectingFiles

		// Build file list message
		var sb strings.Builder
		sb.WriteString("Select files to download by sending their numbers separated by commas (e.g., '1,3,5'):\n\n")

		for i, file := range files {
			//remove server.formatbytes since server.TorrentFile no longer has the Size field.
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, file.Name))
		}

		sb.WriteString("\nOr send 'all' to download all files.")

		// Send file selection message
		updateMsg := tgbotapi.NewEditMessageText(chatID, sentMsg.MessageID, sb.String())
		b.Config.API.Send(updateMsg)
	}()
}

// handleFileSelection processes file selection from user
func (b *Bot) handleFileSelection(message *tgbotapi.Message, session *UserSession) {
	chatID := message.Chat.ID
	selection := strings.TrimSpace(message.Text)

	// Check if user wants all files
	if strings.ToLower(selection) == "all" {
		// Select all files
		var fileIDs []int
		for i := range session.Files {
			fileIDs = append(fileIDs, i)
		}

		err := session.Downloader.SelectFiles(fileIDs)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error selecting files: %v", err))
			b.Config.API.Send(msg)
			return
		}

		b.startDownload(message, session)
		return
	}

	// Parse file numbers
	parts := strings.Split(selection, ",")
	var fileIDs []int

	for _, part := range parts {
		num, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Invalid selection format. Please use numbers separated by commas (e.g., '1,3,5').")
			b.Config.API.Send(msg)
			return
		}

		// Adjust for 0-based indexing
		fileID := num - 1

		// Validate file ID
		if fileID < 0 || fileID >= len(session.Files) {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Invalid file number: %d. Please select numbers between 1 and %d.",
				num, len(session.Files)))
			b.Config.API.Send(msg)
			return
		}

		fileIDs = append(fileIDs, fileID)
	}

	// Apply selection
	err := session.Downloader.SelectFiles(fileIDs)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error selecting files: %v", err))
		b.Config.API.Send(msg)
		return
	}

	b.startDownload(message, session)
}

// startDownload begins the download process
func (b *Bot) startDownload(message *tgbotapi.Message, session *UserSession) {
	chatID := message.Chat.ID

	// Send initial download message
	msg := tgbotapi.NewMessage(chatID, "Starting download...")
	sentMsg, err := b.Config.API.Send(msg)
	if err != nil {
		b.Logger.LogError("Error sending message: %v", err)
		return
	}

	// Store progress message
	session.ProgressMsg = &sentMsg
	session.State = StateDownloading

	// Start download
	progressChan, err := session.Downloader.Download()
	if err != nil {
		updateMsg := tgbotapi.NewEditMessageText(chatID, sentMsg.MessageID,
			fmt.Sprintf("Error starting download: %v", err))
		b.Config.API.Send(updateMsg)
		return
	}

	// Monitor progress
	go func() {
		lastUpdate := time.Now()

		for progress := range progressChan {
			// Update UI every 3 seconds to avoid Telegram API rate limits
			if time.Since(lastUpdate) >= 3*time.Second {
				statusMsg := fmt.Sprintf("Status: %s\nProgress: %.2f%%\nDownloaded: %s / %s\nPeers: %d",
					progress.Status,
					progress.PercentComplete,
					server.FormatBytes(progress.BytesCompleted),
					server.FormatBytes(progress.BytesTotal),
					progress.Peers)

				updateMsg := tgbotapi.NewEditMessageText(chatID, sentMsg.MessageID, statusMsg)
				b.Config.API.Send(updateMsg)
				lastUpdate = time.Now()
			}
		}

		// Download complete, get files
		files, err := session.Downloader.GetDownloadedFiles()
		if err != nil {
			updateMsg := tgbotapi.NewEditMessageText(chatID, sentMsg.MessageID,
				fmt.Sprintf("Download failed: %v", err))
			b.Config.API.Send(updateMsg)
			return
		}

		// Send completion message
		updateMsg := tgbotapi.NewEditMessageText(chatID, sentMsg.MessageID,
			fmt.Sprintf("Download complete! Uploading %d files...", len(files)))
		b.Config.API.Send(updateMsg)

		// Upload files to Telegram
		b.uploadFiles(chatID, files)

		// Reset state
		session.State = StateNone
	}()
}

// uploadFiles uploads downloaded files to Telegram
func (b *Bot) uploadFiles(chatID int64, files []server.TorrentFile) {
	for _, file := range files {
		// Open the file
		f, err := os.Open(file.Path)
		if err != nil {
			errorMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error opening file: %v", err))
			b.Config.API.Send(errorMsg)
			continue
		}
		defer f.Close()

		// Create file upload
		fileUpload := tgbotapi.FileReader{
			Name:   file.Name,
			Reader: f,
		}

		// Check if it's a readable text file
		fileExt := strings.ToLower(filepath.Ext(file.Name))
		if isTextFile(fileExt) {
			// Reset file pointer
			f.Seek(0, 0)

			// Read file content and send as text if small enough
			content, err := io.ReadAll(f)
			if err == nil && len(content) < 4096 { // Telegram message limit
				textMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("File: %s\n\n```\n%s\n```", file.Name, string(content)))
				textMsg.ParseMode = "Markdown"
				b.Config.API.Send(textMsg)
				continue // Skip document upload
			}
			f.Seek(0, 0) // Reset pointer again for document upload if needed
		}

		// Send as document
		doc := tgbotapi.NewDocument(chatID, fileUpload)
		doc.Caption = fmt.Sprintf("File: %s", file.Name)

		// Send the file
		_, err = b.Config.API.Send(doc)
		if err != nil {
			b.Logger.LogError("Failed to upload file %s: %v", file.Name, err)
			errorMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error uploading file %s: %v", file.Name, err))
			b.Config.API.Send(errorMsg)
		}

		// Small delay between uploads
		time.Sleep(1 * time.Second)
	}

	// Send completion message
	completeMsg := tgbotapi.NewMessage(chatID, "All files have been uploaded! Send another magnet link to download more files.")
	b.Config.API.Send(completeMsg)
}

// isTextFile checks if the file extension indicates a text file
func isTextFile(ext string) bool {
	textExtensions := map[string]bool{
		".txt":  true,
		".log":  true,
		".md":   true,
		".json": true,
		".csv":  true,
		".xml":  true,
		".html": true,
		".htm":  true,
		".css":  true,
		".js":   true,
		".py":   true,
		".go":   true,
		".c":    true,
		".cpp":  true,
		".h":    true,
		".java": true,
		".php":  true,
		".rb":   true,
		".sh":   true,
		".bat":  true,
		".ps1":  true,
		".yaml": true,
		".yml":  true,
		".toml": true,
		".ini":  true,
		".cfg":  true,
		".conf": true,
	}

	return textExtensions[ext]
}
