package recorder

import "os"
import "time"
import "bufio"
import "merkur/storage"
import "net/http"
import "sync"

const (
	bufsize = 1024
)

var (
	DefaultOutputExt = ".mp3"
	OutputFolder = "/home/lenni/radio/"
)

var (
	wait = sync.WaitGroup{}
)

type Task struct {
	Station *storage.Station
	Start, End time.Time
}

func (t *Task) Run() error {
	f, err := os.CreateTemp(OutputFolder, "rec-*" + DefaultOutputExt)
	if err != nil {
		return err
	}

	in, err := http.Get(t.Station.Url)
	if err != nil {
		f.Close()
		return err
	}

	inr := bufio.NewReader(in.Body)
	outr := bufio.NewWriter(f)

	go t.record(inr, outr)
	
	return nil
}

func (t *Task) record(in *bufio.Reader, out *bufio.Writer) {
	wait.Add(1)
	buf := make([]byte, bufsize)

	for true {
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
	}

	out.Flush()
	wait.Done()
}

func Wait() {
	wait.Wait()
}
