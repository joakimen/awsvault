package fs

import (
	"os"
	"path/filepath"
)

func XDGCacheDir() string {
	return filepath.Join(os.Getenv("HOME"), ".cache")
}

func XDGDataDir() string {
	return filepath.Join(os.Getenv("HOME"), ".local", "share")
}
