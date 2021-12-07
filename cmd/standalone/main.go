package main

import "time"
import "merkur/storage"
import "merkur/recorder"

func main() {
	s := storage.Station{Name: "SWR3", Url: "https://liveradio.swr.de/sw331ch/swr3/play.aac"}
	t := recorder.Task{Station: &s, Start: time.Now(), End: time.Now().Add(10 * time.Second)}

	t.Run()
	recorder.Wait()
}
