package container

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func IsComponentDirectory(path string) bool {
	buildFile := filepath.Join(path, "Dockerfile")
	_, err := os.Stat(buildFile)
	return !os.IsNotExist(err)
}

func GetComponentDirectories(projectPath string) []string {
	contents, err := ioutil.ReadDir(projectPath)
	components := []string{}

	if err != nil {
		log.Fatal(err)
	}

	for _, content := range contents {
		path := filepath.Join(projectPath, content.Name())
		if content.IsDir() && IsComponentDirectory(path) {
			components = append(components, path)
		}
	}

	return components
}

func IsFileInComponent(componentPath string, filename string) bool {
	return strings.Contains(filename, componentPath)
}

func RunBuild(buildPath string) error {
	log.Printf("Running build on '%s'...\n", buildPath)
	time.Sleep(time.Minute)
	log.Printf("Finished build on '%s'.", buildPath)
	return nil
}

func RunRedeploy(componentName string) error {
	log.Printf("Stopping container '%s'...\n", componentName)
	time.Sleep(3 * time.Second)
	log.Println("Container stopped.")
	log.Printf("Creating/Starting container '%s'...", componentName)
	time.Sleep(time.Second)
	log.Println("Container started.")
	return nil
}
