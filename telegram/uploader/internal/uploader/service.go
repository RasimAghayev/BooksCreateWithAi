package uploader

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"uploader/internal/config"
	"uploader/internal/scanner"
	"uploader/internal/telegram"
)

type Uploader struct {
	config  *config.Config
	scanner *scanner.Scanner
	store   *Store
}

func New(cfg *config.Config, store *Store) *Uploader {
	return &Uploader{
		config:  cfg,
		scanner: scanner.NewScanner(cfg.UploadDir),
		store:   store,
	}
}

func (u *Uploader) Run(ctx context.Context) error {
	logProcess("Uploader started")
	pairs, err := u.scanner.Scan()
	if err != nil {
		logProcess(fmt.Sprintf("Scan failed: %v", err))
		return fmt.Errorf("scan failed: %w", err)
	}

	if len(pairs) == 0 {
		logProcess("No upload pairs found")
		fmt.Println("No upload pairs found.")
		return nil
	}

	logProcess(fmt.Sprintf("Found %d book(s) to process", len(pairs)))

	for _, pair := range pairs {
		if err := u.processBook(ctx, pair); err != nil {
			logProcess(fmt.Sprintf("Error processing %s: %v", pair.BookFile, err))
			fmt.Printf("Error processing %s: %v\n", pair.BookFile, err)
			continue
		}
	}

	logProcess("Uploader finished")
	return nil
}

func (u *Uploader) processBook(ctx context.Context, pair *scanner.UploadPair) error {
	bookID := pair.BookID
	if bookID == "" {
		bookID = GenerateRandomBookID()
	}

	logProcess(fmt.Sprintf("Processing book: %s [book_id=%s]", pair.Metadata.Title, bookID))
	fmt.Printf("Processing book: %s\n", pair.Metadata.Title)

	if u.store.IsBookProcessed(bookID) {
		logProcess(fmt.Sprintf("Book %s [book_id=%s] already processed. Skipping.", pair.Metadata.Title, bookID))
		fmt.Printf("  Book %s already processed. Skipping.\n", pair.Metadata.Title)
		return nil
	}

	channels := u.matchChannels(pair.Metadata)
	if len(channels) == 0 {
		logProcess(fmt.Sprintf("No channels matched for %s [book_id=%s]", pair.Metadata.Title, bookID))
		fmt.Printf("  No channels matched for %s\n", pair.Metadata.Title)
		return nil
	}

	logProcess(fmt.Sprintf("Matched channels for %s [book_id=%s]: %v", pair.Metadata.Title, bookID, channels))
	fmt.Printf("  Matched channels: %v\n", channels)

	_ = u.store.RegisterBook(BookRecord{
		BookID:          bookID,
		BookFile:        pair.BookFile,
		Title:           pair.Metadata.Title,
		TitleOriginal:   pair.Metadata.TitleOriginal,
		Author:          pair.Metadata.Author,
		Year:            pair.Metadata.Year,
		Version:         pair.Metadata.Version,
		Pages:           pair.Metadata.Pages,
		Level:           pair.Metadata.Level,
		LevelName:       pair.Metadata.LevelName,
		PrimaryLanguage: pair.Metadata.PrimaryLanguage,
		Technologies:    pair.Metadata.Technologies,
		Domains:         pair.Metadata.Domains,
		Tags:            pair.Metadata.Tags,
		ProcessedAt:     time.Now(),
	})

	jsonPath := filepath.Join(pair.Dir, pair.BookFile)
	bookPath := findBookFile(jsonPath)

	const uploadAttempts = 2
	allSuccess := true
	for attempt := 1; attempt <= uploadAttempts; attempt++ {
		logProcess(fmt.Sprintf("Upload attempt %d/%d for %s [book_id=%s]", attempt, uploadAttempts, pair.Metadata.Title, bookID))
		attemptSuccess := true
		for _, ch := range channels {
			channelCfg, ok := u.config.Channels[ch]
			if !ok {
				logProcess(fmt.Sprintf("Warning: channel %s not configured for %s [book_id=%s]", ch, pair.Metadata.Title, bookID))
				fmt.Printf("  Warning: channel %s not configured\n", ch)
				attemptSuccess = false
				continue
			}

			existingStatus, err := u.store.GetUploadStatus(bookID)
			if err != nil {
				logProcess(fmt.Sprintf("Warning: failed to get upload status for %s [book_id=%s]: %v", pair.Metadata.Title, bookID, err))
			} else if existingStatus[ch] == "success" {
				logProcess(fmt.Sprintf("Channel %s already has successful upload for %s [book_id=%s]. Skipping.", ch, pair.Metadata.Title, bookID))
				fmt.Printf("  Channel %s already uploaded. Skipping.\n", ch)
				continue
			}

			if err := u.uploadPdfToChannel(ctx, pair, channelCfg, bookPath, bookID); err != nil {
				logProcess(fmt.Sprintf("Upload to %s failed for %s [book_id=%s]: %v", ch, pair.Metadata.Title, bookID, err))
				fmt.Printf("  Upload to %s failed: %v\n", ch, err)
				u.store.MarkFailed(bookID, pair.BookFile, ch, err.Error())
				attemptSuccess = false
				continue
			}

			logProcess(fmt.Sprintf("Uploaded %s to %s successfully (attempt %d) [book_id=%s]", pair.BookFile, ch, attempt, bookID))
			fmt.Printf("  Uploaded to %s successfully (attempt %d)\n", ch, attempt)
		}
		if !attemptSuccess {
			allSuccess = false
		}
	}

	if allSuccess {
		logProcess(fmt.Sprintf("All uploads successful for %s [book_id=%s]. Cleaning up upload package.", pair.Metadata.Title, bookID))
		fmt.Printf("  All uploads successful. Cleaning up %s\n", pair.BookFile)
		u.cleanup(pair)
		u.archiveCompletedBook(pair)
	} else {
		logProcess(fmt.Sprintf("Some uploads failed for %s [book_id=%s]. Files kept for retry.", pair.Metadata.Title, bookID))
		fmt.Printf("  Some uploads failed. Files kept for retry.\n")
	}

	return nil
}

func (u *Uploader) uploadPdfToChannel(ctx context.Context, pair *scanner.UploadPair, ch config.ChannelConfig, bookPath string, bookID string) error {
	jsonPath := filepath.Join(pair.Dir, pair.BookFile)

	bookData, err := os.ReadFile(bookPath)
	if err != nil {
		return fmt.Errorf("failed to read book file: %w", err)
	}

	uploadName := filepath.Base(bookPath)

	caption := buildCaption(pair.Metadata)

	client := telegram.NewClient(ch.Token, ch.ChatID, time.Duration(u.config.APITimeout)*time.Second)
	result, err := client.SendDocument(ctx, uploadName, bookData, caption)
	if err != nil {
		return err
	}

	logUpload(pair.Metadata.Title, ch.Name, result.FileID, result.MessageID)

	if err := u.updateMetadataWithTelegramIDs(jsonPath, ch.Name, result.FileID, result.MessageID); err != nil {
		logProcess(fmt.Sprintf("Warning: failed to update metadata with telegram IDs for %s [book_id=%s]: %v", pair.BookFile, bookID, err))
		fmt.Printf("  Warning: failed to update metadata with telegram IDs: %v\n", err)
	}

	return u.store.MarkUploaded(bookID, pair.BookFile, ch.Name, result.FileID, result.MessageID)
}

func (u *Uploader) updateMetadataWithTelegramIDs(jsonPath, channel, fileID string, messageID int64) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return err
	}

	var metadata map[string]interface{}
	if err := json.Unmarshal(data, &metadata); err != nil {
		return err
	}

	telegramUploads, ok := metadata["telegram_uploads"].([]interface{})
	if !ok {
		telegramUploads = []interface{}{}
	}

	telegramUploads = append(telegramUploads, map[string]interface{}{
		"channel":      channel,
		"file_id":      fileID,
		"message_id":   messageID,
		"uploaded_at":  time.Now().Format(time.RFC3339),
	})

	metadata["telegram_uploads"] = telegramUploads

	updatedData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(jsonPath, updatedData, 0644)
}

func (u *Uploader) matchChannels(meta *scanner.BookMetadata) []string {
	var matched []string
	matchedMap := make(map[string]bool)

	for _, ch := range u.config.Channels {
		if matchedMap[ch.Name] {
			continue
		}

		if u.channelMatches(meta, ch) {
			matched = append(matched, ch.Name)
			matchedMap[ch.Name] = true
		}
	}

	return matched
}

func (u *Uploader) channelMatches(meta *scanner.BookMetadata, ch config.ChannelConfig) bool {
	rules := ch.Rules

	if match, ok := rules["match"].(string); ok && match == "*" {
		return true
	}

	if languages, ok := rules["languages"].([]interface{}); ok {
		for _, l := range languages {
			if lang, ok := l.(string); ok && lang == meta.PrimaryLanguage {
				return true
			}
		}
	}

	if technologies, ok := rules["technologies"].([]interface{}); ok {
		for _, t := range technologies {
			if tech, ok := t.(string); ok {
				for _, mt := range meta.Technologies {
					if mt == tech {
						return true
					}
				}
			}
		}
	}

	if domains, ok := rules["domains"].([]interface{}); ok {
		for _, d := range domains {
			if domain, ok := d.(string); ok {
				for _, md := range meta.Domains {
					if md == domain {
						return true
					}
				}
			}
		}
	}

	if levels, ok := rules["levels"].([]interface{}); ok {
		for _, l := range levels {
			if level, ok := l.(float64); ok && int(level) == meta.Level {
				return true
			}
		}
	}

	return false
}

func (u *Uploader) cleanup(pair *scanner.UploadPair) {
	jsonPath := filepath.Join(pair.Dir, pair.BookFile)
	bookPath := findBookFile(jsonPath)
	_ = os.Remove(bookPath)
	_ = os.Remove(jsonPath)
}

// findBookFile json sidecar yoluna uyğun kitab faylını tapır (hər hansı genişlənmə).
func findBookFile(jsonPath string) string {
	base := strings.TrimSuffix(jsonPath, ".json")
	for _, ext := range []string{".pdf", ".epub", ".djvu"} {
		p := base + ext
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return base + ".pdf"
}

func (u *Uploader) archiveCompletedBook(pair *scanner.UploadPair) {
	exeDir, err := os.Executable()
	if err != nil {
		logProcess(fmt.Sprintf("Failed to get executable path for archiving %s: %v", pair.Metadata.Title, err))
		return
	}

	projectRoot := filepath.Join(filepath.Dir(exeDir), "../../")
	booksDir := filepath.Join(projectRoot, "Books")
	booksReadDir := filepath.Join(projectRoot, "books_read")

	bookDirName := strings.TrimSuffix(pair.BookFile, ".json")
	srcDir := filepath.Join(booksDir, bookDirName)
	// Yeni iyerarxiya (SYSTEM_PROMPT.md Bölmə 3): books_read/{primary_language}/L{level}/{filename_safe}
	langDir := "General"
	if pair.Metadata.PrimaryLanguage != "" {
		langDir = pair.Metadata.PrimaryLanguage
	}
	dstDir := filepath.Join(booksReadDir, langDir, fmt.Sprintf("L%d", pair.Metadata.Level), bookDirName)

	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		logProcess(fmt.Sprintf("Book directory not found for archiving: %s", srcDir))
		return
	}

	if err := os.MkdirAll(filepath.Dir(dstDir), 0755); err != nil {
		logProcess(fmt.Sprintf("Failed to create books_read dir: %v", err))
		return
	}

	if err := os.Rename(srcDir, dstDir); err != nil {
		logProcess(fmt.Sprintf("Failed to move book to books_read: %v", err))
		return
	}

	sourceDir := filepath.Join(dstDir, "source")
	if err := os.RemoveAll(sourceDir); err != nil {
		logProcess(fmt.Sprintf("Failed to delete source dir for %s: %v", bookDirName, err))
	} else {
		logProcess(fmt.Sprintf("Deleted source/ for %s", bookDirName))
	}

	logProcess(fmt.Sprintf("Archived %s to %s", pair.Metadata.Title, dstDir))
	fmt.Printf("  Archived to %s\n", dstDir)
}

func logProcess(message string) {
	logPath := getLogPath("process.log")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(f, "[%s] %s\n", timestamp, message)
}

func logUpload(title, channel, fileID string, messageID int64) {
	logPath := getLogPath("upload.log")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(f, "[%s] BOOK: %s | CHANNEL: %s | FILE_ID: %s | MESSAGE_ID: %d\n", timestamp, title, channel, fileID, messageID)
}

func getLogPath(filename string) string {
	exeDir, err := os.Executable()
	if err != nil {
		return filename
	}
	logDir := filepath.Join(filepath.Dir(exeDir), "../../logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return filename
	}
	return filepath.Join(logDir, filename)
}

func buildCaption(meta *scanner.BookMetadata) string {
	tags := ""
	for i, tag := range meta.Tags {
		if i > 0 {
			tags += " "
		}
		tags += tag
	}

	chapters := ""
	for i, ch := range meta.Chapters {
		if i > 0 {
			chapters += "\n"
		}
		chapters += fmt.Sprintf("📑 %d. %s (%s)", ch.Chapter, ch.Title, ch.Pages)
	}

	return fmt.Sprintf("📚 %s\n👤 %s  \n📅 %d  \n📖 v%d\n💻 %s  \n🎯 Level %d — %s\n🏷 %s\n%s",
		meta.TitleOriginal,
		meta.Author,
		meta.Year,
		meta.Version,
		meta.PrimaryLanguage,
		meta.Level,
		meta.LevelName,
		tags,
		chapters,
	)
}
