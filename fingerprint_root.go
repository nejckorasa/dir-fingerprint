package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

// RFing represents root fingerprint
type RFing struct {
	Fingerprint       string
	FilesFingerprints []string
	Changed           bool // marks if the root fingerprint changed
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
	rfing = RFing{fing, fings, false}
	return
}

// SaveRFing saves fingerprint in path
func SaveRFing(rfing RFing, path string) {

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if !FullFing {
		rfing.FilesFingerprints = nil
	}

	bytes, err := json.MarshalIndent(rfing, "", "  ")
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintln(f, string(bytes))
	if err != nil {
		panic(err)
	}

	return
}

// ReadRFing reads fingerprint in path
func ReadRFing(path string) (rfing *RFing, err error) {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)
	if err := json.Unmarshal(byteValue, &rfing); err != nil {
		return nil, err
	}

	return rfing, nil
}

// Compare compares root fingerprints
func Compare(oldRfing *RFing, newRfing *RFing) bool {
	if oldRfing == nil {
		return true
	}
	return oldRfing.Fingerprint != newRfing.Fingerprint
}
