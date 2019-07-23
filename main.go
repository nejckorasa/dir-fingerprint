package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"strconv"
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

	Log.Infof("Root 	%s", root)
	Log.Infof("File	%s", RFingFileName)

	// files fingerprints
	ffings, skippedCount, err := BuildFFings(root)
	if err != nil {
		Log.Fatal(err)
	}

	// root fingerprint path
	rfingPath := root + RFingFileName

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

	Log.Infof("Took	%.5f sec", time.Since(start).Seconds())
	Log.Infof("For	%d files", len(ffings))
	Log.Infof("Skip	%d files", skippedCount)
	fmt.Println()
	fmt.Printf("Old		[%s]", oldRfing.Fingerprint)
	fmt.Println()
	fmt.Printf("New		[%s]", newRfing.Fingerprint)
	fmt.Println()
	fmt.Printf("@		%s", rfingPath)
	fmt.Println()
	fmt.Printf("Diff		[%s]", strconv.FormatBool(changed))
	fmt.Println()
}

// parseArgs parses flags/args and returns root
func parseArgs(quitLoadingPrint <-chan struct{}) (root string, err error) {
	flag.StringVar(&RFingFileName, "fing", ".fingerprint", "fingerprint file name")
	flag.BoolVar(&FullFing, "f", false, "include all files fingerprints in fingerprint file")
	debug := flag.Bool("d", false, "debug, turn on debug logging")
	quiet := flag.Bool("q", false, "quiet, turn off logging, only print result")

	flag.Parse()

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

	root = args[0]
	return root, nil
}
