package main

import (
	"os"
	"runtime"
	"strconv"
	"syscall"
)

var raws = []string{
	"raw-0",
	"raw-1",
	"raw-2",
	"raw-3",
	"raw-4",
	"raw-5",
	"raw-6",
	"raw-7",
	"raw-8",
	"raw-9",
	"raw-10",
}

func worker(core int, jobs <-chan string, e chan<- error) {
	for job := range jobs {
		e <- &os.PathError{
			Op:   "worker-" + strconv.Itoa(core) + "-core",
			Path: job,
			Err:  syscall.ENODEV,
		}
	}
}

func main() {
	var cores = runtime.NumCPU() / 2
	runtime.GOMAXPROCS(cores)

	jobs := make(chan string, len(raws))
	err := make(chan error, len(raws))
	for core := 0; core <= cores; core++ {
		go worker(core, jobs, err)
	}
	defer close(jobs)
	defer close(err)

	for _, raw := range raws {
		jobs <- raw
	}

	for core := 0; core < len(raws); core++ {
		e := <-err
		println(e.Error())
	}
}
