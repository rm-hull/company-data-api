package internal

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
)

func isUrl(uri string) bool {
	for _, prefix := range []string{"https://", "http://"} {
		if strings.HasPrefix(uri, prefix) {
			return true
		}
	}
	return false
}

func TransientDownload(uri string, handler func(tmpfile string) error) error {
	if !isUrl(uri) {
		return handler(uri)
	}
	log.Printf("Retrieving: %s", uri)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch from %s: %w", uri, err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close body: %v", err)
		}
	}()

	if resp.StatusCode > 299 {
		return fmt.Errorf("http status response from %s: %s", uri, resp.Status)
	}

	tmp, err := os.CreateTemp("", "download-*")
	if err != nil {
		return err
	}
	tmpfile := tmp.Name()
	log.Printf("Downloading content (%s) to %s", humanize.Bytes(uint64(resp.ContentLength)), tmpfile)
	defer func() {
		log.Printf("Removing temporary file: %s", tmpfile)
		if err := os.Remove(tmpfile); err != nil {
			log.Printf("failed to remove file %s: %v", tmpfile, err)
		}
	}()

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		return fmt.Errorf("failed to copy response body: %w", err)
	}
	if err := tmp.Close(); err != nil {
		log.Printf("failed to close temporary file: %v", err)
	}

	return handler(tmpfile)
}
