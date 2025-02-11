package main

import (
	"log"
	"task-manager/cli"
	"task-manager/database"
)

func main() {
	// Initialize MongoDB database instance
	err := database.NewDBInstance()
	if err != nil {
		log.Fatal(err)
	}

	// Start the task CLI
	cli.TaskCLI()
	if err := cli.Run(); err != nil {
		log.Fatal(err)
	}
}
