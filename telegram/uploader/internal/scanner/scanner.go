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
	Title             string        `json:"title"`
	TitleOriginal     string        `json:"title_original"`
	Author            string        `json:"author"`
	Year              int           `json:"year"`
	Version           int           `json:"version"`
	Pages             int           `json:"pages"`
	PrimaryLanguage   string        `json:"primary_language"`
	Level             int           `json:"level"`
	LevelName         string        `json:"level_name"`
	LevelIcon         string        `json:"level_icon,omitempty"`
	Domains           []string      `json:"domains"`
	Technologies      []string      `json:"technologies"`
	Tags              []string      `json:"tags"`
	Chapters          []ChapterInfo `json:"chapters"`
	Summary           string        `json:"summary"`
	ISBN              string        `json:"isbn,omitempty"`
	Publisher         string        `json:"publisher,omitempty"`
	Series            string        `json:"series,omitempty"`
}

type UploadPair struct {
	Dir      string
	BookFile string
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

	return &metadata, nil
}
