package main

import (
	"fmt"

	"github.com/aws/amazon-ecs-event-stream-handler/cmd"
	"github.com/aws/amazon-ecs-event-stream-handler/logger"

	log "github.com/cihub/seelog"

	"os"
)

const errorCode = 1

func main() {
	defer log.Flush()
	err := logger.InitLogger()
	if err != nil {
		fmt.Printf("Could not initialize logger: %+v", err)
	}
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Criticalf("Error executing: %v", err)
		os.Exit(1)
	}
}
