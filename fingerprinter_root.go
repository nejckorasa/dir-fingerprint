package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// RootFingerprint represents root fingerprint that comprises all files fingerprints
type RootFingerprint struct {
	Fingerprint  string
	Fingerprints []string
}

func BuildRootFingerprint(files []FileFingerprint) (rootFingerprint RootFingerprint) {

	var fingerprints []string
	for _, f := range files {
		fingerprints = append(fingerprints, f.Fingerprint)
	}

	hasher := sha256.New()
	hasher.Write([]byte( strings.Join(fingerprints, "")))
	fingerprint := hex.EncodeToString(hasher.Sum(nil))

	rootFingerprint = RootFingerprint{fingerprint, fingerprints}
	return
}

func SaveRootFingerprint(fingerprint RootFingerprint, path string) {

	f, err := os.Create(path)
	check(err)
	defer f.Close()

	bytes, err := json.MarshalIndent(fingerprint, "", "  ")
	check(err)
	_, err = fmt.Fprintln(f, string(bytes))
	check(err)

	return
}
