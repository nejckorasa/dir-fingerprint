package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// FFing represents file fingerprint
type FFing struct {
	Fing, Name string
}

// FWalkResult holds the result of file walk, including error if it occurred
type FWalkResult struct {
	ffing *FFing
	err   error
}

// BuildFFings builds files fingerprints for files in tree rooted at root
func BuildFFings(root string) []FFing {
	cs := make(chan FWalkResult, 10)
	walkFiles(root, cs)
	return receiveFFings(cs)
}

// walkFiles walks the file tree rooted at root and creates fingerprint for each file
// sending FFing to cs channel
func walkFiles(root string, cs chan FWalkResult) {
	wg := &sync.WaitGroup{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		Log.Debugf("Walking	%s", path)
		wg.Add(1)
		go buildFFing(wg, cs, path, info, root, err)
		return nil
	})
	if err != nil {
		panic(err)
	}

	// wait for wait group and close channels
	go func(wg *sync.WaitGroup, cs chan FWalkResult) {
		wg.Wait()
		close(cs)
	}(wg, cs)
}

// buildFFing creates file fingerprint for file on specified path and sends FFing to ch channel
func buildFFing(wg *sync.WaitGroup, cs chan<- FWalkResult, path string, info os.FileInfo, root string, err error) {
	defer wg.Done()

	if err != nil {
		cs <- FWalkResult{nil, err}
		return
	}

	if !info.IsDir() && info.Name() != RFingFileName {
		start := time.Now()
		hasher := sha256.New()

		f, err := os.Open(path)
		if err != nil {
			cs <- FWalkResult{nil, err}
			return
		}

		if _, err := io.Copy(hasher, f); err != nil {
			cs <- FWalkResult{nil, err}
			return
		}
		f.Close()

		ffing := FFing{hex.EncodeToString(hasher.Sum(nil)), info.Name()}
		Log.Infof("Done	[%s](%.4f) 	@ %s", ffing.Fing[:6], time.Since(start).Seconds(), strings.TrimLeft(path, root))
		cs <- FWalkResult{&ffing, nil}
	}
}

// receiveFFings receives FFings from channel
func receiveFFings(cs chan FWalkResult) (files []FFing) {
	done := make(chan bool, 1)

	go func(cs <-chan FWalkResult, done chan<- bool) {
		for walkResult := range cs {

			// if there is an error, panic
			if walkResult.err != nil {
				panic(walkResult.err.Error())
			}
			files = append(files, *walkResult.ffing)
		}
		done <- true
	}(cs, done)

	<-done
	return
}
