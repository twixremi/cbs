package network

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type DownloadResult struct {
	Path string
	Hash string
	Err  error
}

type Client struct {
	httpClient *http.Client
	semaphore  chan struct{}
}

func NewClient(maxConcurrent int) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		semaphore:  make(chan struct{}, maxConcurrent),
	}
}

func (c *Client) Download(urls []string, dest string) (string, error) {
	c.semaphore <- struct{}{}
	defer func() { <-c.semaphore }()

	var lastErr error
	for _, url := range urls {
		fmt.Printf("[NETWORK] Trying mirror: %s\n", url)
		for i := 0; i < 2; i++ {
			hash, err := c.doDownload(url, dest)
			if err == nil {
				return hash, nil
			}
			lastErr = err
			fmt.Printf("[NETWORK] Attempt %d failed for %s: %v\n", i+1, url, err)
			time.Sleep(1 * time.Second)
		}
	}

	return "", fmt.Errorf("all mirrors failed. Last error: %v", lastErr)
}

func (c *Client) doDownload(url, dest string) (string, error) {
	out, err := os.OpenFile(dest, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer out.Close()

	stat, _ := out.Stat()
	start := stat.Size()

	req, _ := http.NewRequest("GET", url, nil)
	if start > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", start))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusRequestedRangeNotSatisfiable {
		// Possibly file already complete or server doesn't support Range
		// For now, assume it might be done or error
	} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	// Verify SHA-256
	return ComputeHash(dest)
}

func ComputeHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
