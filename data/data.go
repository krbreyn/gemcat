package data

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var xdg_data_home string = os.Getenv("XDG_DATA_HOME")

const app_dir = "gemcat"
const data_file = "browser_state"
const cache_dir = "gemcache"

func getAppDir() string {
	var base_data_dir string
	if xdg_data_home != "" {
		base_data_dir = xdg_data_home
	} else {
		base_data_dir = filepath.Join(os.Getenv("HOME"), ".local/share")
	}

	return filepath.Join(base_data_dir, app_dir)
}

func getDataFile() string {
	app_data_dir := getAppDir()

	err := os.MkdirAll(app_data_dir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create data dir: %v\n", err)
		os.Exit(1)
	}

	return filepath.Join(
		app_data_dir,
		fmt.Sprintf("%s.json", data_file),
	)
}

func LoadDataFile() ([]byte, error) {
	dataFile := getDataFile()

	data, err := os.ReadFile(dataFile)
	if err != nil {
		return data, fmt.Errorf("failed to read data file: %w", err)
	}

	return data, nil
}

func SaveDataFile(data []byte) error {
	dataFile := getDataFile()

	err := os.WriteFile(dataFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write data file: %w", err)
	}

	return nil
}

func getCacheDir() string {
	cache_path := filepath.Join(getAppDir(), cache_dir)
	err := os.MkdirAll(cache_path, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create data dir: %v\n", err)
		os.Exit(1)
	}

	return cache_path
}

func CacheGemFile(u *url.URL, content []byte) error {
	relativePath := NormalizeGemPath(u)
	fullPath := filepath.Join(getCacheDir(), relativePath)

	dir := filepath.Dir(fullPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("cache error: failed to create cache subdir: %w", err)
	}

	err = os.WriteFile(fullPath, content, 0644)
	if err != nil {
		return fmt.Errorf("cache error: failed to write cache file: %w", err)
	}

	return nil
}

func IsCacheMiss(err error) bool {
	return strings.Contains(err.Error(), "cache miss")
}

func LoadFromCache(u *url.URL) ([]byte, error) {
	fullPath := filepath.Join(getCacheDir(), NormalizeGemPath(u))

	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("cache miss: no cached file for %s", u.String())
		}
		return nil, fmt.Errorf("failed to stat cache file: %w", err)
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	return content, nil
}

func IsCacheStale(u *url.URL, maxAge time.Duration) (bool, error) {
	fullPath := filepath.Join(getCacheDir(), NormalizeGemPath(u))

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil // stale because it doesn't exist
		}
		return false, fmt.Errorf("failed to stat cache file: %w", err)
	}

	age := time.Since(info.ModTime())
	return age > maxAge, nil
}

func NormalizeGemPath(u *url.URL) string {
	host := u.Host
	path := u.Path

	if path == "" || path == "/" {
		path = "/index.gmi"
	} else if strings.HasSuffix(path, "/") {
		path = filepath.Join(path, "index.gmi")
	} else {
		// If it doesn't look like a file, append .gmi
		if filepath.Ext(path) == "" {
			path += ".gmi"
		}
	}

	// Avoid any weird escaping issues
	return filepath.Join(host, path)
}
