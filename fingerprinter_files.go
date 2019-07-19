package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// FileFingerprint represents file fingerprint
type FileFingerprint struct {
	Fingerprint, Name string
}

// BuildFileFingerprints builds files fingerprints for files in tree rooted at root
func BuildFileFingerprints(root string) []FileFingerprint {
	cs := make(chan FileFingerprint, 10)
	walkFiles(root, cs)
	return receiveFileFingerprints(cs)
}

// walkFiles walks the file tree rooted at root and creates fingerprint for each file
// sending FileFingerprint to cs channel
func walkFiles(root string, cs chan FileFingerprint) {

	wg := &sync.WaitGroup{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		Log.Debug("Walking " + path)
		wg.Add(1)
		go buildFileFingerprint(wg, cs, path, info, root)
		return nil
	})

	check(err)

	// wait for wait group and close channel
	go func(wg *sync.WaitGroup, cs chan FileFingerprint) {
		wg.Wait()
		close(cs)
	}(wg, cs)
}

// buildFileFingerprint creates file fingerprint for file on specified path and sends FileFingerprint to ch channel
func buildFileFingerprint(wg *sync.WaitGroup, cs chan<- FileFingerprint, path string, info os.FileInfo, root string) {

	defer wg.Done()

	if !info.IsDir() && info.Name() != FingerprintFileName {
		hasher := sha256.New()

		f, err := os.Open(path)
		check(err)

		defer f.Close()

		_, err = io.Copy(hasher, f)
		check(err)

		metadata := FileFingerprint{hex.EncodeToString(hasher.Sum(nil)), info.Name()}
		cs <- metadata

		Log.Infof("Done		[%s] -> %s", metadata.Fingerprint[:6], strings.TrimLeft(path, root))
	}
}

// receiveFileFingerprints receives FileFingerprints from channel
func receiveFileFingerprints(cs chan FileFingerprint) (files []FileFingerprint) {

	done := make(chan bool, 1)

	go func(cs <-chan FileFingerprint, done chan<- bool) {

		for file := range cs {
			files = append(files, file)
		}

		done <- true
	}(cs, done)

	<-done
	return
}

