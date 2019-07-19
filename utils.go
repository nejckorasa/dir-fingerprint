package main

import "time"

// timeTrack tracks time of a method
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	Log.Infof("%s took 		%f seconds", name, elapsed.Seconds())
}

// check checks the error
func check(e error) {
	if e != nil {
		panic(e)
	}
}
