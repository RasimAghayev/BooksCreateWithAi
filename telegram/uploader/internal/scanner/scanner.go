package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type ChapterInfo struct {
	Chapter int    `json:"chapter"`
	Title   string `json:"title"`
	Pages   string `json:"pages"`
}

type BookMetadata struct {
	BookID           string        `json:"book_id"`
	Title            string        `json:"title"`
	TitleOriginal    string        `json:"title_original"`
	Author           string        `json:"author"`
	Year             int           `json:"year"`
	Version          int           `json:"version"`
	Pages            int           `json:"pages"`
	PrimaryLanguage  string        `json:"primary_language"`
	Level            int           `json:"level"`
	LevelName        string        `json:"level_name"`
	LevelIcon        string        `json:"level_icon,omitempty"`
	Domains          []string      `json:"domains"`
	Technologies     []string      `json:"technologies"`
	Tags             []string      `json:"tags"`
	Chapters         []ChapterInfo `json:"chapters"`
	Summary          string        `json:"summary"`
	ISBN             string        `json:"isbn,omitempty"`
	Publisher        string        `json:"publisher,omitempty"`
	Series           string        `json:"series,omitempty"`
	TelegramUploads  []TelegramUpload `json:"telegram_uploads,omitempty"`
}

type TelegramUpload struct {
	Channel     string    `json:"channel"`
	FileID      string    `json:"file_id"`
	MessageID   int64     `json:"message_id"`
	UploadedAt  string    `json:"uploaded_at"`
}

type UploadPair struct {
	Dir      string
	BookFile string
	BookID   string
	Metadata *BookMetadata
}

type Scanner struct {
	UploadDir string
}

func NewScanner(uploadDir string) *Scanner {
	return &Scanner{UploadDir: uploadDir}
}

func (s *Scanner) Scan() ([]*UploadPair, error) {
	entries, err := os.ReadDir(s.UploadDir)
	if err != nil {
		return nil, err
	}

	var pairs []*UploadPair

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if ext != ".json" {
			continue
		}

		metadata, err := s.loadMetadata(filepath.Join(s.UploadDir, entry.Name()))
		if err != nil {
			fmt.Printf("Warning: failed to load metadata for %s: %v\n", entry.Name(), err)
			continue
		}

		pairs = append(pairs, &UploadPair{
			Dir:      s.UploadDir,
			BookFile: entry.Name(),
			BookID:   metadata.BookID,
			Metadata: metadata,
		})
	}

	return pairs, nil
}

func (s *Scanner) loadMetadata(path string) (*BookMetadata, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var metadata BookMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}

	if metadata.Title == "" {
		var nested struct {
			BookID    string `json:"book_id"`
			Book struct {
				Title   string `json:"title"`
				Author  string `json:"author"`
				Year    int    `json:"year"`
				Version int    `json:"version"`
				Pages   int    `json:"page_count"`
			} `json:"book"`
			TitleOriginal string   `json:"title_original"`
			Classification struct {
				PrimaryLanguage string   `json:"primary_language"`
				Level           int      `json:"level"`
				LevelName       string   `json:"level_name"`
				LevelIcon       string   `json:"level_icon"`
				Domains         []string `json:"domains"`
				Technologies    []string `json:"technologies"`
				Tags            []string `json:"tags"`
			} `json:"classification"`
			Summary  string        `json:"summary"`
			Chapters []ChapterInfo `json:"chapters"`
			ISBN     string        `json:"isbn"`
			Publisher string       `json:"publisher"`
			Series   string        `json:"series"`
		}
		if err := json.Unmarshal(data, &nested); err == nil {
			if metadata.BookID == "" {
				metadata.BookID = nested.BookID
			}
			if metadata.Title == "" {
				metadata.Title = nested.Book.Title
			}
			if metadata.TitleOriginal == "" {
				metadata.TitleOriginal = nested.Book.Title
			}
			if metadata.Author == "" {
				metadata.Author = nested.Book.Author
			}
			if metadata.Year == 0 {
				metadata.Year = nested.Book.Year
			}
			if metadata.Version == 0 {
				metadata.Version = nested.Book.Version
			}
			if metadata.Pages == 0 {
				metadata.Pages = nested.Book.Pages
			}
			if metadata.PrimaryLanguage == "" {
				metadata.PrimaryLanguage = nested.Classification.PrimaryLanguage
			}
			if metadata.Level == 0 {
				metadata.Level = nested.Classification.Level
			}
			if metadata.LevelName == "" {
				metadata.LevelName = nested.Classification.LevelName
			}
			if len(metadata.Domains) == 0 {
				metadata.Domains = nested.Classification.Domains
			}
			if len(metadata.Technologies) == 0 {
				metadata.Technologies = nested.Classification.Technologies
			}
			if len(metadata.Tags) == 0 {
				metadata.Tags = nested.Classification.Tags
			}
			if metadata.Summary == "" {
				metadata.Summary = nested.Summary
			}
			if len(metadata.Chapters) == 0 {
				metadata.Chapters = nested.Chapters
			}
			if metadata.ISBN == "" {
				metadata.ISBN = nested.ISBN
			}
			if metadata.Publisher == "" {
				metadata.Publisher = nested.Publisher
			}
			if metadata.Series == "" {
				metadata.Series = nested.Series
			}
		}
	}

	return &metadata, nil
}
