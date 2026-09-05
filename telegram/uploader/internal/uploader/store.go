package uploader

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type BookRecord struct {
	BookFile      string    `json:"book_file"`
	Title         string    `json:"title"`
	TitleOriginal string    `json:"title_original"`
	Author        string    `json:"author"`
	Year          int       `json:"year"`
	Version       int       `json:"version"`
	Pages         int       `json:"pages"`
	Level         int       `json:"level"`
	LevelName     string    `json:"level_name"`
	PrimaryLanguage string `json:"primary_language"`
	Technologies  []string  `json:"technologies"`
	Domains       []string  `json:"domains"`
	Tags          []string  `json:"tags"`
	ProcessedAt   time.Time `json:"processed_at"`
	UploadedAt    time.Time `json:"uploaded_at,omitempty"`
	TelegramFileID string `json:"telegram_file_id,omitempty"`
	TelegramMessageID int64 `json:"telegram_message_id,omitempty"`
}

type ChapterInfo struct {
	Chapter int    `json:"chapter"`
	Title   string `json:"title"`
	Pages   string `json:"pages"`
}

type UploadRecord struct {
	BookFile     string    `json:"book_file"`
	Channel      string    `json:"channel"`
	Status       string    `json:"status"`
	ProcessedAt  time.Time `json:"processed_at"`
	UploadedAt   time.Time `json:"uploaded_at,omitempty"`
	TelegramFileID string `json:"telegram_file_id,omitempty"`
	TelegramMessageID int64 `json:"telegram_message_id,omitempty"`
}

type Store struct {
	path string
	mu   sync.RWMutex
	data map[string]map[string]UploadRecord
	books map[string]BookRecord
}

func NewStore(dbPath string) (*Store, error) {
	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, err
	}

	s := &Store{
		path:  absPath,
		data:  make(map[string]map[string]UploadRecord),
		books: make(map[string]BookRecord),
	}

	if err := s.load(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if len(data) == 0 {
		return nil
	}

	var records []UploadRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return err
	}

	s.data = make(map[string]map[string]UploadRecord)
	for _, r := range records {
		if s.data[r.BookFile] == nil {
			s.data[r.BookFile] = make(map[string]UploadRecord)
		}
		s.data[r.BookFile][r.Channel] = r
	}

	return nil
}

func (s *Store) save() error {
	s.mu.RLock()
	var records []UploadRecord
	for _, channels := range s.data {
		for _, r := range channels {
			records = append(records, r)
		}
	}
	s.mu.RUnlock()

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, data, 0644)
}

func (s *Store) IsBookProcessed(bookFile string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.books[bookFile]
	return ok
}

func (s *Store) RegisterBook(record BookRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.books[record.BookFile] = record
	return s.saveBooks()
}

func (s *Store) saveBooks() error {
	data, err := json.MarshalIndent(s.books, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func (s *Store) GetUploadStatus(bookFile string) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	channels, ok := s.data[bookFile]
	if !ok {
		return make(map[string]string), nil
	}

	status := make(map[string]string)
	for ch, r := range channels {
		status[ch] = r.Status
	}

	return status, nil
}

func (s *Store) MarkUploaded(bookFile, channel, fileID string, messageID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data[bookFile] == nil {
		s.data[bookFile] = make(map[string]UploadRecord)
	}

	s.data[bookFile][channel] = UploadRecord{
		BookFile:          bookFile,
		Channel:           channel,
		Status:            "success",
		ProcessedAt:       time.Now(),
		UploadedAt:        time.Now(),
		TelegramFileID:    fileID,
		TelegramMessageID: messageID,
	}

	return s.save()
}

func (s *Store) MarkFailed(bookFile, channel string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data[bookFile] == nil {
		s.data[bookFile] = make(map[string]UploadRecord)
	}

	s.data[bookFile][channel] = UploadRecord{
		BookFile:    bookFile,
		Channel:     channel,
		Status:      "failed",
		ProcessedAt: time.Now(),
		UploadedAt:  time.Now(),
	}

	return s.save()
}

func (s *Store) DeleteRecord(bookFile string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, bookFile)
	delete(s.books, bookFile)
	return s.save()
}

func (s *Store) Close() error {
	return nil
}

func GenerateBookID(title, author string, year, version, pages int) string {
	h := sha256.New()
	h.Write([]byte(title))
	h.Write([]byte(author))
	h.Write([]byte(string(rune(year))))
	h.Write([]byte(string(rune(version))))
	h.Write([]byte(string(rune(pages))))
	return hex.EncodeToString(h.Sum(nil))[:16]
}

func GenerateRandomBookID() string {
	buf := make([]byte, 16)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}
