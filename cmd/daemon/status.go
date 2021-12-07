package daemon

import (
	"fmt"
	"io"
)

func (d *Daemon) EchoStatus(w io.Writer) {
	fmt.Fprintf(w, "Daemon knows %d stations.\n", len(d.Stations))
	fmt.Fprintf(w, "Currently waiting for %d tasks,\n", len(d.waiting_tasks))
	fmt.Fprintf(w, "recording %d tasks\n", len(d.active_tasks))
	fmt.Fprintf(w, "with %d being finished.\n", len(d.done_tasks))
}
