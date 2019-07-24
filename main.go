package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
	"strconv"
	"strings"
	"time"
)

var Log = logrus.New()

// fingerprint file name
var RFingFileName string

// create full fingerprint file = include all files fingerprints
var FullFing bool

func init() {
	Log.Formatter = new(prefixed.TextFormatter)
}

func main() {
	start := time.Now()

	quitLoadingPrint := make(chan struct{})
	defer close(quitLoadingPrint)
	root, err := parseArgs(quitLoadingPrint)
	if err != nil {
		Log.Fatal(err)
		return
	}
	if root == nil {
		return
	}

	Log.Infof("Root 	%s", *root)
	Log.Infof("File	%s", RFingFileName)

	// files fingerprints
	ffings, skippedCount, err := BuildFFings(*root)
	if err != nil {
		Log.Fatal(err)
	}

	// root fingerprint path
	rfingPath := *root + RFingFileName

	// load old root fingerprint
	oldRfing, err := ReadRFing(rfingPath)
	if err != nil {
		Log.Error(err)
	}

	// creat new root fingerprint
	newRfing := BuildRFing(ffings)

	// indicates if the fingerprint has changed
	changed := Compare(oldRfing, &newRfing)
	newRfing.Changed = changed

	// save new root fingerprint
	SaveRFing(newRfing, rfingPath)

	outputResults(time.Since(start).Seconds(), ffings, skippedCount, oldRfing, newRfing, rfingPath, changed)
}

// parseArgs parses flags/args and returns root
func parseArgs(quitLoadingPrint <-chan struct{}) (root *string, err error) {
	flag.StringVar(&RFingFileName, "f", ".fingerprint", "fingerprint file name")
	flag.BoolVar(&FullFing, "files", false, "files, include all files fingerprints in fingerprint file, mind that there might me a lot of them")
	debug := flag.Bool("d", false, "debug, turn on debug logging")
	quiet := flag.Bool("q", false, "quiet, turn off logging, only print result")
	help := flag.Bool("h", false, "help, display usage")

	flag.Parse()

	if *help {
		flag.Usage()
		return nil, nil
	}

	switch {
	case *debug:
		Log.Level = logrus.DebugLevel
	case *quiet:
		// extra logging to show we are still alive
		go func() {
			for {
				select {
				case <-quitLoadingPrint:
					return
				default:
					for range time.Tick(2 * time.Second) {
						fmt.Printf(".")
					}
				}
			}
		}()
		Log.Level = logrus.FatalLevel
	default:
		Log.Level = logrus.InfoLevel
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return root, errors.New("wrong usage, missing root argument")
	}

	root = &args[0]
	if !strings.HasSuffix(*root, string(os.PathSeparator)) {
		tmp := *root + string(os.PathSeparator)
		root = &tmp
	}
	return root, nil
}

func outputResults(seconds float64, ffings []FFing, skippedCount int, oldRfing *RFing, newRfing RFing, rfingPath string, rfingChanged bool) {
	Log.Infof("Took	%.5f sec", seconds)
	Log.Infof("For	%d files", len(ffings))
	Log.Infof("Skip	%d files\n", skippedCount)
	oldFingerprint := ""
	if oldRfing != nil {
		oldFingerprint = oldRfing.Fingerprint
	}
	fmt.Printf("Old		[%s]\n", oldFingerprint)
	fmt.Printf("New		[%s]\n", newRfing.Fingerprint)
	fmt.Printf("@		%s\n", rfingPath)
	fmt.Printf("Diff		%s\n", strconv.FormatBool(rfingChanged))
}
