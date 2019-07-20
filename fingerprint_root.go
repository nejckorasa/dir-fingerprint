package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// RFing represents root fingerprint that comprises all files fingerprints
type RFing struct {
	root  string
	files []string
}

// BuildRFing builds fingerprint for all provided files
func BuildRFing(ffings []FFing) (rfing RFing) {
	var fings []string
	for _, f := range ffings {
		fings = append(fings, f.Fing)
	}

	// need to sort fingerprints to ensure correct root fingerprint
	sort.Strings(fings)

	hasher := sha256.New()
	if _, err := hasher.Write([]byte(strings.Join(fings, ""))); err != nil {
		panic(err)
	}

	fing := hex.EncodeToString(hasher.Sum(nil))
	rfing = RFing{fing, fings}
	return
}

// SaveRFing saves fingerprint in path
func SaveRFing(fingerprint RFing, path string) {

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	bytes, err := json.MarshalIndent(fingerprint, "", "  ")
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintln(f, string(bytes))
	if err != nil {
		panic(err)
	}

	return
}
