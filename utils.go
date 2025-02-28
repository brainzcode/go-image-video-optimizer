package main

import (
	"path/filepath"
	"strings"
)

func getFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".webp", ".tiff", ".bmp", ".dng":
		return "image"
	case ".mp4", ".avi", ".mov", ".mkv", ".flv", ".wmv":
		return "video"
	default:
		return "unknown"
	}
}
