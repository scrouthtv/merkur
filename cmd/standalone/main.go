package main

import (
	"fmt"
	"merkur/recorder"
	"merkur/storage"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 4 {
		printHelp()
		os.Exit(1)
	}

	url := os.Args[1]
	folder := os.Args[2]
	duration, err := time.ParseDuration(os.Args[3])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		printHelp()
		os.Exit(1)
	}

	if duration < 3*time.Second {
		fmt.Fprintln(os.Stderr, "Refusing to record less than 3 seconds")
		printHelp()
		os.Exit(1)
	}

	recorder.OutputFolder = folder
	s := storage.Station{Url: url}
	t := recorder.Task{Station: &s, Start: time.Now(), End: time.Now().Add(duration)}

	err = t.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		printHelp()
		os.Exit(1)
	}

	recorder.Wait()
}

func printHelp() {
	fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "<url> <folder> <duration>")
}
