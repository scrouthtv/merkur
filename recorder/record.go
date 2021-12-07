package recorder

import (
	"bufio"
	"merkur/storage"
	"net/http"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

const (
	bufsize = 1024
)

var (
	OutputFolder = "/mnt/d/tmp/radio/"
)

var (
	wait = sync.WaitGroup{}
)

type Task struct {
	Station    *storage.Station
	Start, End time.Time
}

func init() {
	if runtime.GOOS == "windows" {
		OutputFolder = "D:\\tmp\\radio\\"
	}
}

func (t *Task) Run() error {
	ext := path.Ext(t.Station.Url)
	f, err := os.CreateTemp(OutputFolder, "rec-*"+ext)
	if err != nil {
		return err
	}

	req, err := http.Get(t.Station.Url)
	if err != nil {
		return err
	}

	inr := bufio.NewReader(req.Body)
	outr := bufio.NewWriter(f)

	wait.Add(1)
	go t.record(inr, outr)

	return nil
}

func (t *Task) record(in *bufio.Reader, out *bufio.Writer) {
	buf := make([]byte, bufsize)
	total := 0

	for {
		if time.Now().After(t.End) {
			println("Ending recording.")
			break
		}

		n, err := in.Read(buf)
		if err != nil {
			panic(err)
		}

		_, err = out.Write(buf[:n])
		if err != nil {
			panic(err)
		}

		total += n
	}

	//println("Written", total, "bytes.")
	out.Flush()
	wait.Done()
}

func Wait() {
	wait.Wait()
}
