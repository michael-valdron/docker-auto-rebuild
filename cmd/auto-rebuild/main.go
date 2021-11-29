package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/michael-valdron/docker-auto-rebuild/pkg/watcher"
	"github.com/sevlyar/go-daemon"
)

const NAME = "auto-rebuild"

type Flags struct {
	signal *string
	target *string
}

func createFlags() *Flags {
	return &Flags{
		signal: flag.String("s", "", `Send signal to the daemon:
			quit — graceful shutdown
			stop — fast shutdown
			reload — reloading the configuration file`),
		target: flag.String("t", ".", "Target directory to watch, e.g. /path/to/target"),
	}
}

func createArgsString(flags *Flags) []string {
	args := []string{NAME}

	if *flags.target != "." {
		args[0] = args[0] + fmt.Sprintf(" -t %s", *flags.target)
	}

	return args
}

func runWatch(stop <-chan bool, done chan<- bool, workingDir string) {
	log.Println("- - - - - - - - - - - - - - -")
	log.Printf("Watching '%s'...\n", workingDir)
	watcher.Watch(stop, done)
}

func main() {
	flags := createFlags()

	flag.Parse()

	watcher.RunDaemon(&daemon.Context{
		PidFileName: fmt.Sprintf("%s.pid", NAME),
		PidFilePerm: 0644,
		LogFileName: fmt.Sprintf("%s.log", NAME),
		LogFilePerm: 0640,
		WorkDir:     *flags.target,
		Umask:       027,
		Args:        createArgsString(flags),
	}, func(cxt *daemon.Context, stop <-chan bool, done chan<- bool) {
		runWatch(stop, done, cxt.WorkDir)
	}, flags.signal)
}
