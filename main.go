package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"
	"time"

	"github.com/k0kubun/pp"
)

const (
	daemonStart = 1 + iota
	daemonSuccess
	daemonFailed
)

var (
	filename string
	dir      string
	wait     string
	waitdur  time.Duration
)

func init() {
	var err error
	waitdur, err = time.ParseDuration(wait)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	// ldflags
	const layout = `1405`
	pp.Println(waitdur)
	Daemon(func() {
		for {
			copyloop(dir, waitdur, func(imgpath string) string {
				pp.Println(imgpath)
				ex, err := os.Executable()
				if err != nil {
					panic(err)
				}
				exPath := path.Dir(ex)
				ext := filepath.Ext(imgpath)
				return path.Join(exPath, filename+time.Now().Format(layout)+ext)
			})
		}
	})
}

func parent() (err error) {
	args := []string{"--ld"}
	args = append(args, os.Args[1:]...)

	r, w, err := os.Pipe()
	if err != nil {
		return err
	}

	cmd := exec.Command(os.Args[0], args...)
	cmd.ExtraFiles = []*os.File{w}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err = cmd.Start(); err != nil {
		return err
	}

	ch := make(chan int, 1)
	go func() {
		buf := make([]byte, 1)
		r.Read(buf)

		if int(buf[0]) == daemonSuccess {
			ch <- daemonSuccess
		}
	}()

	select {
	case <-time.After(30 * time.Second):
		os.Exit(1)
		return errors.New("time up: it cannot start child process")
	case <-ch:
		return nil
	}
}

// Daemon daemonize served function process easily.
func Daemon(fn func()) {
	// is daemon or parent?
	child := flag.Bool("ld", false, "ld process")
	flag.Parse()

	if !*child {
		if err := parent(); err != nil {
			log.Fatalf("Error occurred [%v]", err)
			os.Exit(1)
		}
		return
	}

	var err error
	// notify the status of a child process to the parent process.
	pipe := os.NewFile(uintptr(3), "pipe")
	if pipe != nil {
		defer pipe.Close()
		if err != nil {
			pipe.Write([]byte{daemonFailed})
		} else {
			pipe.Write([]byte{daemonSuccess})
		}
	}

	signal.Ignore(syscall.SIGCHLD)

	syscall.Close(0)
	syscall.Close(1)
	syscall.Close(2)

	// be a process group leader
	syscall.Setsid()

	syscall.Umask(022)

	syscall.Chdir("/")

	// main
	fn()
}
