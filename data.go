package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var xdg_data_home string = os.Getenv("XDG_DATA_HOME")

const app_dir = "gemcat"
const data_file = "browser_data"
const cache_dir = "gemcache"

type Data struct {
	Bookmarks []string `json:"bookmarks"`
	History   []string `json:"history"`
}

func getAppDir() string {
	var base_data_dir string
	if xdg_data_home != "" {
		base_data_dir = xdg_data_home
	} else {
		base_data_dir = filepath.Join(os.Getenv("HOME"), ".locali/share")
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

func LoadDataFile() (Data, error) {
	var data Data
	dataFile := getDataFile()

	_, err := os.Stat(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return data, fmt.Errorf("failed to stat data file: %w", err)
	}

	jsonBytes, err := os.ReadFile(dataFile)
	if err != nil {
		return data, fmt.Errorf("failed to read data file: %w", err)
	}

	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return data, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return data, nil
}

func SaveDataFile(data Data) error {
	dataFile := getDataFile()

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	err = os.WriteFile(dataFile, jsonBytes, 0644)
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

func CacheGemFile(rawurl string, content []byte) error {
	full_path, err := pathForURL(rawurl)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	dir := filepath.Dir(full_path)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create cache subdir: %w", err)
	}

	err = os.WriteFile(full_path, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

func LoadFromCache(rawurl string) ([]byte, error) {
	full_path, err := pathForURL(rawurl)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	_, err = os.Stat(full_path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("cache miss: no cached file for %s", rawurl)
		}
		return nil, fmt.Errorf("failed to stat cache file: %w", err)
	}

	content, err := os.ReadFile(full_path)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	return content, nil
}

func IsCacheStale(rawurl string, maxAge time.Duration) (bool, error) {
	full_path, err := pathForURL(rawurl)
	if err != nil {
		return false, fmt.Errorf("invalud URL: %w", err)
	}

	info, err := os.Stat(full_path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil // stale because it doesn't exist
		}
		return false, fmt.Errorf("failed to stat cache file: %w", err)
	}

	age := time.Since(info.ModTime())
	return age > maxAge, nil
}

func pathForURL(rawurl string) (full_path string, err error) {
	host, path, err := parseURL(rawurl)
	if err != nil {
		return "", fmt.Errorf("invalud URL: %w", err)
	}

	return filepath.Join(getCacheDir(), host, path), nil
}

func parseURL(rawurl string) (host, path string, err error) {
	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		return "", "", fmt.Errorf("invalud URL: %w", err)
	}

	host = parsedURL.Host
	path = parsedURL.Path
	if path == "" || path == "/" {
		path = "/index.gmi"
	}

	return host, path, nil
}
