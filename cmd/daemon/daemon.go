package daemon

import (
	"merkur/recorder"
	"merkur/storage"
	"time"
)

var (
	recordLoopPause = 1 * time.Second
)

type Daemon struct {
	Stations      []*storage.Station
	waiting_tasks []*recorder.Task
	active_tasks  []*recorder.Task
	done_tasks    []*recorder.Task
}

func (d *Daemon) AddTask(task *recorder.Task) {
	d.waiting_tasks = append(d.waiting_tasks, task)
}

func (d *Daemon) RecordLoop() {
	for {
		next, task := d.NextStart()
		if next != nil && time.Until(*next) <= recordLoopPause {
			d.start(task)
		}

		time.Sleep(1 * time.Second)
	}
}

func (d *Daemon) NextStart() (t *time.Time, task *recorder.Task) {
	t = nil
	task = nil

	for _, tk := range d.waiting_tasks {
		if t == nil || tk.Start.Before(*t) {
			t = &tk.Start
			task = tk
		}
	}

	return
}

func (d *Daemon) start(task *recorder.Task) {
	task.Done = func() {
		d.done(task)
	}

	for i, t := range d.waiting_tasks {
		if t == task {
			d.waiting_tasks = append(d.waiting_tasks[:i], d.waiting_tasks[i+1:]...)
		}
	}

	d.active_tasks = append(d.active_tasks, task)
	task.Run()
}

func (d *Daemon) done(task *recorder.Task) {
	for i, t := range d.active_tasks {
		if t == task {
			d.active_tasks = append(d.active_tasks[:i], d.active_tasks[i+1:]...)
		}
	}

	d.done_tasks = append(d.done_tasks, task)
}
