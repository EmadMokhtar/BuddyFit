package internal

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Video struct {
	Author     string
	Title      string
	ID         string
	Transcript string
}

func (v *Video) URL() string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", v.ID)
}

func (v *Video) UnescapedURL() string {
	unescape, err := url.QueryUnescape(v.URL())
	if err != nil {
		return ""
	}
	return unescape
}

func NewVideoFromFileName(fileName, filePath string) (*Video, error) {
	filNameNoExt := strings.TrimSuffix(fileName, ".en.srt.txt")

	parts := strings.Split(filNameNoExt, Sep)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid file name format: %s", fileName)
	}
	title := parts[1]
	id := parts[0]
	author := filepath.Dir(filePath)

	content, err := os.ReadFile(filepath.Join(filePath, fileName))
	if err != nil {
		return nil, err
	}

	return &Video{
		Title:      title,
		ID:         id,
		Author:     author,
		Transcript: string(content),
	}, nil
}
