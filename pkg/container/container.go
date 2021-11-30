package container

import (
	"log"
	"time"
)

func RunBuild() error {
	log.Println("Running build...")
	time.Sleep(time.Minute)
	log.Println("Finished build.")
	return nil
}

func RunRedeploy() error {
	log.Println("Stopping container...")
	time.Sleep(3 * time.Second)
	log.Println("Container stopped.")
	log.Println("Creating/Starting container...")
	time.Sleep(time.Second)
	log.Println("Container started.")
	return nil
}
