package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

func TestParallelDownload(t *testing.T) {
	filename := "5MB.zip"
	expectedHash := "b3215c06647bc550406a9c8ccc378756"

	oldArgs := os.Args
	os.Args = []string{"pget", "-chunks", "10", "http://ipv4.download.thinkbroadband.com/5MB.zip"}
	defer func() { os.Args = oldArgs }()

	main()
	defer func() {
		err := os.Remove(filename)
		if err != nil {
			log.Printf("error deleting file: %s\n", err)
		}
	}()

	hash, err := md5Sum(filename)
	if err != nil {
		t.Fatalf("error hashing file: %s", err)
	}

	if hash != expectedHash {
		t.Fatalf("hash = %s, want %s", hash, expectedHash)
	}
}

func md5Sum(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("error closing file: %s\n", err)
		}
	}()

	hash := md5.New()

	_, err = io.Copy(hash, file)
	if err != nil {
		return "", fmt.Errorf("error hashing file: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)[:16]), nil
}
