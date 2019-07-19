package main

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"time"
)

var Log = logrus.New()
var FingerprintFileName = ".fingerprint"

var root = "/Users/nejckorasa/Downloads/"

func init() {
	Log.Formatter = new(prefixed.TextFormatter)
	Log.Level = logrus.DebugLevel
}

func main() {

	defer timeTrack(time.Now(), "All")

	Log.Infof("Root = 	%s", root)

	fileFingerprints := BuildFileFingerprints(root)
	rootFingerprint := BuildRootFingerprint(fileFingerprints)

	rootFingerprintLocation := root + FingerprintFileName
	SaveRootFingerprint(rootFingerprint, rootFingerprintLocation)

	Log.Infof("Processed		%d fileFingerprints", len(fileFingerprints))
	Log.Infof("Fingerprint	%s", rootFingerprint)
	Log.Infof("Location		%s", rootFingerprintLocation)
}

