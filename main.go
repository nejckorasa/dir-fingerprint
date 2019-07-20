package main

import (
	"errors"
	"flag"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"time"
)

var Log = logrus.New()

// fingerprint file name
var RFingFileName string

func init() {
	Log.Formatter = new(prefixed.TextFormatter)
}

func main() {
	start := time.Now()

	root, err := parseArgs()
	if err != nil {
		Log.Fatal(err)
		return
	}

	Log.Infof("Root 	%s", root)
	Log.Infof("File	%s", RFingFileName)

	// files fingerprints
	ffings := BuildFFings(root)
	// root fingerprint
	rfing := BuildRFing(ffings)

	rfingPath := root + RFingFileName
	SaveRFing(rfing, rfingPath)

	Log.Infof("All Done	[%s](%.4f) 	@ %s", rfing.root[:6], time.Since(start).Seconds(), rfingPath)
	Log.Infof("Processed	%d files", len(ffings))
	Log.Infof("RFprint	[%s]", rfing.root)
	Log.Debugf("FFprints	%s", rfing.files)
}

// parseArgs parses flags/args and returns root
func parseArgs() (root string, err error) {
	flag.StringVar(&RFingFileName, "fing", ".fingerprint", "fingerprint file name")
	debug := flag.Bool("d", false, "debug")

	flag.Parse()

	if *debug {
		Log.Level = logrus.DebugLevel
	} else {
		Log.Level = logrus.InfoLevel
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return root, errors.New("wrong usage, missing root argument")
	}

	root = args[0]
	return root, nil
}
