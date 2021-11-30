package watcher

import (
	"log"
	"os"
	"syscall"

	"github.com/michael-valdron/docker-auto-rebuild/pkg/container"
	"github.com/sevlyar/go-daemon"
)

var (
	stop = make(chan bool)
	done = make(chan bool)
)

func termHandler(sig os.Signal) error {
	log.Println("terminating...")
	stop <- true
	if sig == syscall.SIGQUIT {
		<-done
	}
	return daemon.ErrStop
}

func redeployHandler(sig os.Signal) error {
	return container.RunRedeploy()
}

func addCommands(signal *string) {
	daemon.AddCommand(daemon.StringFlag(signal, "quit"), syscall.SIGQUIT, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "redeploy"), syscall.SIGHUP, redeployHandler)
}

func sendSignal(cxt *daemon.Context) error {
	d, err := cxt.Search()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	daemon.SendCommands(d)
	return err
}

func listenForSignals() error {
	err := daemon.ServeSignals()
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
	return err
}

func RunDaemon(cxt *daemon.Context, worker func(*daemon.Context, <-chan bool, chan<- bool), signal *string) {
	addCommands(signal)

	if len(daemon.ActiveFlags()) > 0 {
		sendSignal(cxt)
		return
	}

	d, err := cxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cxt.Release()

	go worker(cxt, stop, done)

	listenForSignals()

	log.Println("daemon terminated")
}
