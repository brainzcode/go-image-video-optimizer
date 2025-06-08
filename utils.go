package main

import (
	"path/filepath"
	"strings"
)

func GetFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".webp", ".tiff", ".bmp", ".dng":
		return "image"
	case ".mp4", ".avi", ".mov", ".mkv", ".flv", ".wmv", ".webm", ".m4v", ".3gp", ".gif":
		return "video"
	default:
		return "unknown"
	}
}
