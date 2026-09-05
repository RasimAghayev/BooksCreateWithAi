package uploader

import (
	"context"
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
	pairs, err := u.scanner.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	if len(pairs) == 0 {
		fmt.Println("No upload pairs found.")
		return nil
	}

	for _, pair := range pairs {
		if err := u.processBook(ctx, pair); err != nil {
			fmt.Printf("Error processing %s: %v\n", pair.BookFile, err)
			continue
		}
	}

	return nil
}

func (u *Uploader) processBook(ctx context.Context, pair *scanner.UploadPair) error {
	fmt.Printf("Processing book: %s\n", pair.Metadata.Title)

	if u.store.IsBookProcessed(pair.BookFile) {
		fmt.Printf("  Book %s already processed. Skipping.\n", pair.BookFile)
		return nil
	}

	channels := u.matchChannels(pair.Metadata)
	if len(channels) == 0 {
		fmt.Printf("  No channels matched for %s\n", pair.Metadata.Title)
		return nil
	}

	fmt.Printf("  Matched channels: %v\n", channels)

	_ = u.store.RegisterBook(BookRecord{
		BookFile:       pair.BookFile,
		Title:          pair.Metadata.Title,
		TitleOriginal:  pair.Metadata.TitleOriginal,
		Author:         pair.Metadata.Author,
		Year:           pair.Metadata.Year,
		Version:        pair.Metadata.Version,
		Pages:          pair.Metadata.Pages,
		Level:          pair.Metadata.Level,
		LevelName:      pair.Metadata.LevelName,
		PrimaryLanguage: pair.Metadata.PrimaryLanguage,
		Technologies:   pair.Metadata.Technologies,
		Domains:        pair.Metadata.Domains,
		Tags:           pair.Metadata.Tags,
		ProcessedAt:    time.Now(),
	})

	jsonPath := filepath.Join(pair.Dir, pair.BookFile)

	allSuccess := true
	for _, ch := range channels {
		channelCfg, ok := u.config.Channels[ch]
		if !ok {
			fmt.Printf("  Warning: channel %s not configured\n", ch)
			allSuccess = false
			continue
		}

		if err := u.uploadMetadataToChannel(ctx, pair, channelCfg); err != nil {
			fmt.Printf("  Upload to %s failed: %v\n", ch, err)
			u.store.MarkFailed(pair.BookFile, ch)
			allSuccess = false
			continue
		}

		fmt.Printf("  Uploaded to %s successfully\n", ch)
	}

	if allSuccess {
		fmt.Printf("  All uploads successful. Cleaning up %s\n", pair.BookFile)
		u.cleanup(jsonPath)
	} else {
		fmt.Printf("  Some uploads failed. Files kept for retry.\n")
	}

	return nil
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

func (u *Uploader) uploadMetadataToChannel(ctx context.Context, pair *scanner.UploadPair, ch config.ChannelConfig) error {
	jsonPath := filepath.Join(pair.Dir, pair.BookFile)

	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read metadata file: %w", err)
	}

	caption := buildCaption(pair.Metadata)

	client := telegram.NewClient(ch.Token, ch.ChatID, time.Duration(u.config.APITimeout)*time.Second)
	result, err := client.SendDocument(ctx, pair.BookFile, jsonData, caption)
	if err != nil {
		return err
	}

	return u.store.MarkUploaded(pair.BookFile, ch.Name, result.FileID, result.MessageID)
}

func (u *Uploader) cleanup(jsonPath string) {
	_ = os.Remove(jsonPath)
	pdfPath := strings.Replace(jsonPath, ".json", ".pdf", 1)
	_ = os.Remove(pdfPath)
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
			chapters += " | "
		}
		chapters += fmt.Sprintf("%d. %s (%s)", ch.Chapter, ch.Title, ch.Pages)
	}

	return fmt.Sprintf("📚 %s\n👤 %s   📅 %d   📖 v%d\n💻 %s   🎯 Level %d — %s\n🏷 %s\n📑 %s",
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
