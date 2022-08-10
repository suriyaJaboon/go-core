package main

import (
	"log"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

const format = "2006010215:04"

type E struct {
	K int
	V string
}

type fsn struct {
	wtc  *fsnotify.Watcher
	path string
	e    chan E
	err  chan error
	done chan bool
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = watcher.Close() }()

	f := fsn{
		wtc:  watcher,
		path: "/tmp/logs",
		e:    make(chan E),
		err:  make(chan error),
		done: make(chan bool),
	}

	defer close(f.e)
	defer close(f.err)
	defer close(f.done)

	//go f.wk(1)
	for _, c := range []int{1, 2, 3, 4} {
		go f.wk(c)
	}

	go f.errs()

	if err = f.wtc.Add(f.path); err != nil {
		panic(err)
	}
	go f.events()

	<-f.done
}

func (fs *fsn) events() {
	for {
		select {
		case event, ok := <-fs.wtc.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				//log.Println("modified file:", event.Name)
				fs.rd(time.Now())
			}
		case err, ok := <-fs.wtc.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func (fs *fsn) rd(t time.Time) {
	if rds, err := os.ReadDir(fs.path); err == nil {
		for k, rd := range rds {
			var f os.FileInfo
			if f, err = rd.Info(); err == nil {
				if f.ModTime().Format(format) != t.Format(format) {
					fs.e <- E{K: k, V: rd.Name()}
					//log.Println(fi.ModTime().Format(format), rd.Name())
				}
			}
		}
	}
}

func (fs *fsn) wk(core int) {
	for e := range fs.e {
		//time.Sleep(time.Second * time.Duration(file.k))
		if e.K == 2 {
			fs.err <- &os.PathError{
				Op:   "worker-" + strconv.Itoa(core) + "-core",
				Path: e.V,
				Err:  syscall.ENODEV,
			}
		} else {
			log.Println(core, e.V)
		}
	}
}

func (fs *fsn) errs() {
	for e := range fs.err {
		log.Println("ERROR->", e)
	}
}
