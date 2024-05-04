package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadAndSave(url string) (*os.File, error) {
	confPath := getConfigPath()
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	out, err := os.Create(confPath)

	if err != nil {
		return nil, err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return out, err
}

func downloadImage(url string) (string, error) {
	res, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", errors.New("non 200 status code")
	}

	cacheDir := os.TempDir()

	file, err := os.Create(filepath.Join(cacheDir, "wallpaper"))

	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, res.Body)

	if err != nil {
		return "", err
	}

	err = file.Close()

	if err != nil {
		return "!", err
	}

	return file.Name(), nil
}

func SetFromUrl(url string) error {
	file, err := downloadImage(url)

	if err != nil {
		return err
	}

	return SetFromFile(file)
}

func GetResolutionName(width, height int) string {
	switch {
	case width == 3840 && height == 2160:
		return "4k"
	case width == 2560 && height == 1440:
		return "1440p"
	case width == 1920 && height == 1080:
		return "1080p"
	case width == 1280 && height == 720:
		return "720p"
	default:
		return fmt.Sprintf("%dx%d", width, height)
	}
}
